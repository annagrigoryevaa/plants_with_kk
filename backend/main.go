package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Port              string
	DataDir           string
	AppBaseURL        string
	KeycloakIssuerURL string
	KeycloakClientID  string
}

type App struct {
	cfg          Config
	store        *Store
	tokenChecker *TokenValidator
	mux          *http.ServeMux
}

type Store struct {
	mu      sync.RWMutex
	data    PersistedData
	path    string
	baseDir string
}

type PersistedData struct {
	Plants []Plant         `json:"plants"`
	Events []CalendarEvent `json:"events"`
}

type Plant struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Species   string       `json:"species"`
	Room      string       `json:"room"`
	Notes     string       `json:"notes"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Photos    []PlantPhoto `json:"photos"`
}

type PlantPhoto struct {
	ID         string    `json:"id"`
	URL        string    `json:"url"`
	FileName   string    `json:"fileName"`
	Note       string    `json:"note"`
	CapturedAt time.Time `json:"capturedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CalendarEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Date      string    `json:"date"`
	Actions   []string  `json:"actions"`
	PlantIDs  []string  `json:"plantIds"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PlantPayload struct {
	Name    string `json:"name"`
	Species string `json:"species"`
	Room    string `json:"room"`
	Notes   string `json:"notes"`
}

type EventPayload struct {
	Title    string   `json:"title"`
	Date     string   `json:"date"`
	Actions  []string `json:"actions"`
	PlantIDs []string `json:"plantIds"`
	Notes    string   `json:"notes"`
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type TokenValidator struct {
	issuerURL string
	clientID  string
	http      *http.Client

	mu         sync.RWMutex
	jwksURI    string
	keys       map[string]*rsa.PublicKey
	lastSynced time.Time
}

type openIDConfiguration struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

type jwksPayload struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	KID string   `json:"kid"`
	KTY string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5C []string `json:"x5c"`
}

type tokenClaims struct {
	Issuer            string      `json:"iss"`
	Audience          interface{} `json:"aud"`
	AuthorizedParty   string      `json:"azp"`
	Expiration        int64       `json:"exp"`
	NotBefore         int64       `json:"nbf"`
	IssuedAt          int64       `json:"iat"`
	PreferredUsername string      `json:"preferred_username"`
	Email             string      `json:"email"`
	Name              string      `json:"name"`
}

func main() {
	cfg := mustLoadConfig()
	store, err := NewStore(cfg.DataDir)
	if err != nil {
		log.Fatalf("store init failed: %v", err)
	}

	validator := NewTokenValidator(cfg.KeycloakIssuerURL, cfg.KeycloakClientID)
	app := &App{
		cfg:          cfg,
		store:        store,
		tokenChecker: validator,
		mux:          http.NewServeMux(),
	}
	app.routes()

	log.Printf("plant tracker listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, app.withLogging(app.mux)); err != nil {
		log.Fatal(err)
	}
}

func mustLoadConfig() Config {
	cfg := Config{
		Port:              valueOrDefault("PORT", "8080"),
		DataDir:           valueOrDefault("DATA_DIR", "./data"),
		AppBaseURL:        valueOrDefault("APP_BASE_URL", "http://localhost:8080"),
		KeycloakIssuerURL: os.Getenv("KEYCLOAK_ISSUER_URL"),
		KeycloakClientID:  os.Getenv("KEYCLOAK_CLIENT_ID"),
	}

	if cfg.KeycloakIssuerURL == "" || cfg.KeycloakClientID == "" {
		log.Fatal("KEYCLOAK_ISSUER_URL and KEYCLOAK_CLIENT_ID are required")
	}

	return cfg
}

func NewStore(baseDir string) (*Store, error) {
	if err := os.MkdirAll(filepath.Join(baseDir, "uploads"), 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join(baseDir, "db.json")
	s := &Store{
		path:    path,
		baseDir: baseDir,
		data: PersistedData{
			Plants: []Plant{},
			Events: []CalendarEvent{},
		},
	}

	body, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, s.saveLocked()
		}
		return nil, err
	}

	if len(body) > 0 {
		if err := json.Unmarshal(body, &s.data); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) saveLocked() error {
	body, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, body, 0o644)
}

func (s *Store) Snapshot() PersistedData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return PersistedData{
		Plants: append([]Plant(nil), s.data.Plants...),
		Events: append([]CalendarEvent(nil), s.data.Events...),
	}
}

func (s *Store) CreatePlant(payload PlantPayload) (Plant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	plant := Plant{
		ID:        newID("plant"),
		Name:      strings.TrimSpace(payload.Name),
		Species:   strings.TrimSpace(payload.Species),
		Room:      strings.TrimSpace(payload.Room),
		Notes:     strings.TrimSpace(payload.Notes),
		CreatedAt: now,
		UpdatedAt: now,
		Photos:    []PlantPhoto{},
	}
	s.data.Plants = append(s.data.Plants, plant)
	return plant, s.saveLocked()
}

func (s *Store) UpdatePlant(id string, payload PlantPayload) (Plant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Plants {
		if s.data.Plants[i].ID == id {
			s.data.Plants[i].Name = strings.TrimSpace(payload.Name)
			s.data.Plants[i].Species = strings.TrimSpace(payload.Species)
			s.data.Plants[i].Room = strings.TrimSpace(payload.Room)
			s.data.Plants[i].Notes = strings.TrimSpace(payload.Notes)
			s.data.Plants[i].UpdatedAt = time.Now().UTC()
			return s.data.Plants[i], s.saveLocked()
		}
	}

	return Plant{}, os.ErrNotExist
}

func (s *Store) DeletePlant(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := slices.IndexFunc(s.data.Plants, func(plant Plant) bool {
		return plant.ID == id
	})
	if index < 0 {
		return os.ErrNotExist
	}

	s.data.Plants = append(s.data.Plants[:index], s.data.Plants[index+1:]...)
	for i := range s.data.Events {
		s.data.Events[i].PlantIDs = deleteString(s.data.Events[i].PlantIDs, id)
	}
	return s.saveLocked()
}

func (s *Store) AddPhoto(plantID string, file io.Reader, originalName, note string, capturedAt time.Time) (PlantPhoto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Plants {
		if s.data.Plants[i].ID != plantID {
			continue
		}

		photoID := newID("photo")
		ext := strings.ToLower(filepath.Ext(originalName))
		if ext == "" {
			ext = ".jpg"
		}
		fileName := photoID + ext
		relativePath := filepath.Join("uploads", fileName)
		absolutePath := filepath.Join(s.baseDir, relativePath)

		dst, err := os.Create(absolutePath)
		if err != nil {
			return PlantPhoto{}, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			return PlantPhoto{}, err
		}

		photo := PlantPhoto{
			ID:         photoID,
			URL:        "/uploads/" + fileName,
			FileName:   originalName,
			Note:       strings.TrimSpace(note),
			CapturedAt: capturedAt.UTC(),
			CreatedAt:  time.Now().UTC(),
		}
		s.data.Plants[i].Photos = append([]PlantPhoto{photo}, s.data.Plants[i].Photos...)
		s.data.Plants[i].UpdatedAt = time.Now().UTC()
		return photo, s.saveLocked()
	}

	return PlantPhoto{}, os.ErrNotExist
}

func (s *Store) CreateEvent(payload EventPayload) (CalendarEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	event := CalendarEvent{
		ID:        newID("event"),
		Title:     strings.TrimSpace(payload.Title),
		Date:      payload.Date,
		Actions:   uniqueStrings(payload.Actions),
		PlantIDs:  uniqueStrings(payload.PlantIDs),
		Notes:     strings.TrimSpace(payload.Notes),
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.data.Events = append(s.data.Events, event)
	return event, s.saveLocked()
}

func (s *Store) UpdateEvent(id string, payload EventPayload) (CalendarEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Events {
		if s.data.Events[i].ID == id {
			s.data.Events[i].Title = strings.TrimSpace(payload.Title)
			s.data.Events[i].Date = payload.Date
			s.data.Events[i].Actions = uniqueStrings(payload.Actions)
			s.data.Events[i].PlantIDs = uniqueStrings(payload.PlantIDs)
			s.data.Events[i].Notes = strings.TrimSpace(payload.Notes)
			s.data.Events[i].UpdatedAt = time.Now().UTC()
			return s.data.Events[i], s.saveLocked()
		}
	}

	return CalendarEvent{}, os.ErrNotExist
}

func (s *Store) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := slices.IndexFunc(s.data.Events, func(event CalendarEvent) bool {
		return event.ID == id
	})
	if index < 0 {
		return os.ErrNotExist
	}

	s.data.Events = append(s.data.Events[:index], s.data.Events[index+1:]...)
	return s.saveLocked()
}

func NewTokenValidator(issuerURL, clientID string) *TokenValidator {
	return &TokenValidator{
		issuerURL: strings.TrimSuffix(issuerURL, "/"),
		clientID:  clientID,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		keys: make(map[string]*rsa.PublicKey),
	}
}

func (v *TokenValidator) Validate(ctx context.Context, rawToken string) (UserInfo, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return UserInfo{}, errors.New("invalid JWT format")
	}

	headerBytes, err := decodeJWTPart(parts[0])
	if err != nil {
		return UserInfo{}, err
	}

	var header struct {
		Alg string `json:"alg"`
		KID string `json:"kid"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return UserInfo{}, err
	}
	if header.Alg != "RS256" {
		return UserInfo{}, fmt.Errorf("unsupported alg %s", header.Alg)
	}

	key, err := v.publicKey(ctx, header.KID)
	if err != nil {
		return UserInfo{}, err
	}

	signingInput := parts[0] + "." + parts[1]
	signature, err := decodeJWTPart(parts[2])
	if err != nil {
		return UserInfo{}, err
	}
	digest := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature); err != nil {
		return UserInfo{}, fmt.Errorf("signature verification failed: %w", err)
	}

	payloadBytes, err := decodeJWTPart(parts[1])
	if err != nil {
		return UserInfo{}, err
	}
	var claims tokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return UserInfo{}, err
	}
	if err := v.validateClaims(claims); err != nil {
		return UserInfo{}, err
	}

	return UserInfo{
		Username: claims.PreferredUsername,
		Email:    claims.Email,
		Name:     claims.Name,
	}, nil
}

func (v *TokenValidator) validateClaims(claims tokenClaims) error {
	now := time.Now().Unix()
	if claims.Issuer != v.issuerURL {
		return fmt.Errorf("unexpected issuer %s", claims.Issuer)
	}
	if claims.Expiration < now {
		return errors.New("token expired")
	}
	if claims.NotBefore > 0 && claims.NotBefore > now {
		return errors.New("token is not active yet")
	}
	if !audienceContains(claims.Audience, v.clientID) && claims.AuthorizedParty != v.clientID {
		return errors.New("unexpected audience")
	}
	return nil
}

func (v *TokenValidator) publicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	key := v.keys[kid]
	refreshNeeded := time.Since(v.lastSynced) > time.Hour || len(v.keys) == 0
	v.mu.RUnlock()

	if key != nil && !refreshNeeded {
		return key, nil
	}

	if err := v.refreshKeys(ctx); err != nil {
		return nil, err
	}

	v.mu.RLock()
	defer v.mu.RUnlock()
	key = v.keys[kid]
	if key == nil {
		return nil, fmt.Errorf("kid %s not found in JWKS", kid)
	}
	return key, nil
}

func (v *TokenValidator) refreshKeys(ctx context.Context) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if time.Since(v.lastSynced) <= time.Minute && len(v.keys) > 0 {
		return nil
	}

	wellKnownURL := v.issuerURL + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnownURL, nil)
	if err != nil {
		return err
	}
	resp, err := v.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cfg openIDConfiguration
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, cfg.JWKSURI, nil)
	if err != nil {
		return err
	}
	resp, err = v.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var jwks jwksPayload
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, jwk := range jwks.Keys {
		if jwk.KTY != "RSA" || (jwk.Use != "" && jwk.Use != "sig") {
			continue
		}
		key, err := convertJWKToRSA(jwk)
		if err == nil {
			keys[jwk.KID] = key
		}
	}

	if len(keys) == 0 {
		return errors.New("no usable JWKS keys loaded")
	}

	v.jwksURI = cfg.JWKSURI
	v.keys = keys
	v.lastSynced = time.Now()
	return nil
}

func convertJWKToRSA(jwk jwkKey) (*rsa.PublicKey, error) {
	if len(jwk.X5C) > 0 {
		block, _ := pem.Decode([]byte("-----BEGIN CERTIFICATE-----\n" + jwk.X5C[0] + "\n-----END CERTIFICATE-----"))
		if block == nil {
			return nil, errors.New("invalid x5c certificate")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err == nil {
			if key, ok := cert.PublicKey.(*rsa.PublicKey); ok {
				return key, nil
			}
		}
	}

	modulus, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	exponent, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	var expInt int
	for _, b := range exponent {
		expInt = expInt<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulus),
		E: expInt,
	}, nil
}

func (a *App) routes() {
	a.mux.HandleFunc("GET /api/health", a.handleHealth)
	a.mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(filepath.Join(a.cfg.DataDir, "uploads")))))
	a.mux.HandleFunc("GET /app-config.js", a.handleAppConfig)
	a.mux.Handle("GET /api/me", a.withAuth(http.HandlerFunc(a.handleMe)))
	a.mux.Handle("GET /api/dashboard", a.withAuth(http.HandlerFunc(a.handleDashboard)))
	a.mux.Handle("GET /api/plants", a.withAuth(http.HandlerFunc(a.handlePlantsList)))
	a.mux.Handle("POST /api/plants", a.withAuth(http.HandlerFunc(a.handlePlantCreate)))
	a.mux.Handle("PUT /api/plants/{id}", a.withAuth(http.HandlerFunc(a.handlePlantUpdate)))
	a.mux.Handle("DELETE /api/plants/{id}", a.withAuth(http.HandlerFunc(a.handlePlantDelete)))
	a.mux.Handle("POST /api/plants/{id}/photos", a.withAuth(http.HandlerFunc(a.handlePhotoUpload)))
	a.mux.Handle("GET /api/events", a.withAuth(http.HandlerFunc(a.handleEventsList)))
	a.mux.Handle("POST /api/events", a.withAuth(http.HandlerFunc(a.handleEventCreate)))
	a.mux.Handle("PUT /api/events/{id}", a.withAuth(http.HandlerFunc(a.handleEventUpdate)))
	a.mux.Handle("DELETE /api/events/{id}", a.withAuth(http.HandlerFunc(a.handleEventDelete)))

	publicDir := filepath.Join("..", "frontend", "public")
	a.mux.Handle("/", http.FileServer(http.Dir(publicDir)))
}

func (a *App) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

func (a *App) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" || token == r.Header.Get("Authorization") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
			return
		}

		user, err := a.tokenChecker.Validate(r.Context(), token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type userContextKey struct{}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleAppConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	fmt.Fprintf(w, "window.APP_CONFIG = %s;", mustJSON(map[string]string{
		"apiBaseUrl":          a.cfg.AppBaseURL,
		"keycloakIssuerUrl":   a.cfg.KeycloakIssuerURL,
		"keycloakClientId":    a.cfg.KeycloakClientID,
		"keycloakRedirectUri": a.cfg.AppBaseURL,
	}))
}

func (a *App) handleMe(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(userContextKey{}).(UserInfo)
	writeJSON(w, http.StatusOK, user)
}

func (a *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.Snapshot())
}

func (a *App) handlePlantsList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.Snapshot().Plants)
}

func (a *App) handlePlantCreate(w http.ResponseWriter, r *http.Request) {
	var payload PlantPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if strings.TrimSpace(payload.Name) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	plant, err := a.store.CreatePlant(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, plant)
}

func (a *App) handlePlantUpdate(w http.ResponseWriter, r *http.Request) {
	var payload PlantPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	plant, err := a.store.UpdatePlant(r.PathValue("id"), payload)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, plant)
}

func (a *App) handlePlantDelete(w http.ResponseWriter, r *http.Request) {
	if err := a.store.DeletePlant(r.PathValue("id")); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *App) handlePhotoUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart payload"})
		return
	}

	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "photo file is required"})
		return
	}
	defer file.Close()

	capturedAt := time.Now()
	if raw := r.FormValue("capturedAt"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err == nil {
			capturedAt = parsed
		}
	}

	photo, err := a.store.AddPhoto(r.PathValue("id"), file, fileHeader.Filename, r.FormValue("note"), capturedAt)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, photo)
}

func (a *App) handleEventsList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.Snapshot().Events)
}

func (a *App) handleEventCreate(w http.ResponseWriter, r *http.Request) {
	payload, ok := decodeEventPayload(w, r)
	if !ok {
		return
	}
	event, err := a.store.CreateEvent(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

func (a *App) handleEventUpdate(w http.ResponseWriter, r *http.Request) {
	payload, ok := decodeEventPayload(w, r)
	if !ok {
		return
	}
	event, err := a.store.UpdateEvent(r.PathValue("id"), payload)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, event)
}

func (a *App) handleEventDelete(w http.ResponseWriter, r *http.Request) {
	if err := a.store.DeleteEvent(r.PathValue("id")); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func decodeEventPayload(w http.ResponseWriter, r *http.Request) (EventPayload, bool) {
	var payload EventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return EventPayload{}, false
	}
	if payload.Date == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date is required"})
		return EventPayload{}, false
	}
	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date must be YYYY-MM-DD"})
		return EventPayload{}, false
	}
	if strings.TrimSpace(payload.Title) == "" {
		payload.Title = "Уход за растениями"
	}
	return payload, true
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func decodeJWTPart(part string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(part)
}

func audienceContains(raw interface{}, value string) bool {
	switch aud := raw.(type) {
	case string:
		return aud == value
	case []interface{}:
		for _, item := range aud {
			if s, ok := item.(string); ok && s == value {
				return true
			}
		}
	}
	return false
}

func newID(prefix string) string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s_%x", prefix, buf[:])
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func deleteString(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func mustJSON(value any) string {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(body)
}
