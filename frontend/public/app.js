const config = window.APP_CONFIG;
const state = {
  accessToken: localStorage.getItem("pk_access_token"),
  refreshToken: localStorage.getItem("pk_refresh_token"),
  tokenExpiry: Number(localStorage.getItem("pk_token_expiry") || "0"),
  user: null,
  plants: [],
  events: [],
  monthCursor: new Date(),
};

const profileName = document.querySelector("#profileName");
const profileEmail = document.querySelector("#profileEmail");
const logoutButton = document.querySelector("#logoutButton");
const plantForm = document.querySelector("#plantForm");
const plantsGrid = document.querySelector("#plantsGrid");
const eventForm = document.querySelector("#eventForm");
const eventPlantsSelect = document.querySelector("#eventPlantsSelect");
const monthLabel = document.querySelector("#monthLabel");
const calendarGrid = document.querySelector("#calendarGrid");
const prevMonth = document.querySelector("#prevMonth");
const nextMonth = document.querySelector("#nextMonth");
const plantsCount = document.querySelector("#plantsCount");
const photosCount = document.querySelector("#photosCount");
const eventsCount = document.querySelector("#eventsCount");
const plantTemplate = document.querySelector("#plantCardTemplate");

const actionLabels = {
  "полив": "Полив",
  "пересадка": "Пересадка",
  "подкормка": "Подкормка",
};

boot().catch((error) => {
  console.error(error);
  alert(`Не удалось загрузить приложение: ${error.message}`);
});

async function boot() {
  bindEvents();
  await ensureAuthenticated();
  await loadDashboard();
}

function bindEvents() {
  plantForm.addEventListener("submit", createPlant);
  eventForm.addEventListener("submit", createEvent);
  logoutButton.addEventListener("click", logout);
  prevMonth.addEventListener("click", () => {
    state.monthCursor = addMonths(state.monthCursor, -1);
    renderCalendar();
  });
  nextMonth.addEventListener("click", () => {
    state.monthCursor = addMonths(state.monthCursor, 1);
    renderCalendar();
  });
}

async function ensureAuthenticated() {
  if (config.authMode === "oauth2-proxy") {
    state.user = await api("/api/me");
    if (!state.user) {
      window.location.assign(config.oauth2ProxySignInUrl || "/oauth2/sign_in");
      return;
    }
    profileName.textContent = state.user.name || state.user.username;
    profileEmail.textContent = state.user.email || "";
    return;
  }

  const code = new URLSearchParams(window.location.search).get("code");
  if (code) {
    await exchangeCodeForTokens(code);
    window.history.replaceState({}, "", "/");
  }

  if (!state.accessToken || Date.now() >= state.tokenExpiry - 15_000) {
    if (state.refreshToken) {
      const refreshed = await refreshSession();
      if (!refreshed) {
        await startLogin();
        return;
      }
    } else {
      await startLogin();
      return;
    }
  }

  state.user = await api("/api/me");
  profileName.textContent = state.user.name || state.user.username;
  profileEmail.textContent = state.user.email || "";
}

async function startLogin() {
  const verifier = randomString(64);
  const challenge = await base64UrlEncodeSha256(verifier);
  localStorage.setItem("pk_code_verifier", verifier);

  const params = new URLSearchParams({
    client_id: config.keycloakClientId,
    response_type: "code",
    scope: "openid profile email",
    redirect_uri: config.keycloakRedirectUri,
    code_challenge: challenge,
    code_challenge_method: "S256",
  });

  window.location.assign(`${config.keycloakIssuerUrl}/protocol/openid-connect/auth?${params.toString()}`);
}

async function exchangeCodeForTokens(code) {
  const verifier = localStorage.getItem("pk_code_verifier");
  if (!verifier) {
    throw new Error("PKCE verifier не найден");
  }

  const response = await fetch(`${config.keycloakIssuerUrl}/protocol/openid-connect/token`, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "authorization_code",
      client_id: config.keycloakClientId,
      code,
      redirect_uri: config.keycloakRedirectUri,
      code_verifier: verifier,
    }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Не удалось обменять code на токен: ${text}`);
  }

  saveSession(await response.json());
  localStorage.removeItem("pk_code_verifier");
}

async function refreshSession() {
  const response = await fetch(`${config.keycloakIssuerUrl}/protocol/openid-connect/token`, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "refresh_token",
      client_id: config.keycloakClientId,
      refresh_token: state.refreshToken,
    }),
  });

  if (!response.ok) {
    clearSession();
    return false;
  }

  saveSession(await response.json());
  return true;
}

function saveSession(tokens) {
  state.accessToken = tokens.access_token;
  state.refreshToken = tokens.refresh_token;
  state.tokenExpiry = Date.now() + tokens.expires_in * 1000;
  localStorage.setItem("pk_access_token", state.accessToken);
  localStorage.setItem("pk_refresh_token", state.refreshToken);
  localStorage.setItem("pk_token_expiry", String(state.tokenExpiry));
}

function clearSession() {
  state.accessToken = "";
  state.refreshToken = "";
  state.tokenExpiry = 0;
  localStorage.removeItem("pk_access_token");
  localStorage.removeItem("pk_refresh_token");
  localStorage.removeItem("pk_token_expiry");
}

async function logout() {
  if (config.authMode === "oauth2-proxy") {
    window.location.assign(config.oauth2ProxySignOutUrl || "/oauth2/sign_out");
    return;
  }

  clearSession();
  window.location.assign(
    `${config.keycloakIssuerUrl}/protocol/openid-connect/logout?post_logout_redirect_uri=${encodeURIComponent(config.keycloakRedirectUri)}`
  );
}

async function loadDashboard() {
  const data = await api("/api/dashboard");
  state.plants = data.plants.sort((a, b) => a.name.localeCompare(b.name, "ru"));
  state.events = data.events.sort((a, b) => a.date.localeCompare(b.date));
  renderAll();
}

function renderAll() {
  renderStats();
  renderPlantOptions();
  renderPlants();
  renderCalendar();
}

function renderStats() {
  plantsCount.textContent = String(state.plants.length);
  photosCount.textContent = String(state.plants.reduce((total, plant) => total + plant.photos.length, 0));
  eventsCount.textContent = String(state.events.length);
}

function renderPlantOptions() {
  eventPlantsSelect.innerHTML = "";
  for (const plant of state.plants) {
    const option = document.createElement("option");
    option.value = plant.id;
    option.textContent = `${plant.name}${plant.room ? ` (${plant.room})` : ""}`;
    eventPlantsSelect.append(option);
  }
}

function renderPlants() {
  plantsGrid.innerHTML = "";
  if (state.plants.length === 0) {
    plantsGrid.innerHTML = `<div class="empty-state">Коллекция пока пустая. Добавьте первое растение слева.</div>`;
    return;
  }

  for (const plant of state.plants) {
    const node = plantTemplate.content.firstElementChild.cloneNode(true);
    node.querySelector(".plant-name").textContent = plant.name;
    node.querySelector(".plant-species").textContent = plant.species || "Вид не указан";
    node.querySelector(".plant-room").textContent = plant.room || "Локация не указана";
    node.querySelector(".plant-notes").textContent = plant.notes || "Без заметок";
    node.querySelector(".delete-plant").addEventListener("click", () => deletePlant(plant.id));

    const photoStrip = node.querySelector(".photo-strip");
    if (plant.photos.length === 0) {
      photoStrip.innerHTML = `<div class="empty-state">Фотопленка пока пустая</div>`;
    } else {
      for (const photo of plant.photos) {
        const frame = document.createElement("figure");
        frame.className = "photo-frame";
        frame.innerHTML = `
          <img src="${photo.url}" alt="${escapeHTML(plant.name)}" />
          <figcaption>
            <strong>${new Date(photo.capturedAt).toLocaleDateString("ru-RU")}</strong><br />
            ${escapeHTML(photo.note || "Без комментария")}
          </figcaption>
        `;
        photoStrip.append(frame);
      }
    }

    node.querySelector(".photo-form").addEventListener("submit", (event) => uploadPhoto(event, plant.id));
    plantsGrid.append(node);
  }
}

function renderCalendar() {
  const monthStart = new Date(state.monthCursor.getFullYear(), state.monthCursor.getMonth(), 1);
  const monthEnd = new Date(state.monthCursor.getFullYear(), state.monthCursor.getMonth() + 1, 0);
  const gridStart = new Date(monthStart);
  const weekday = (monthStart.getDay() + 6) % 7;
  gridStart.setDate(monthStart.getDate() - weekday);

  monthLabel.textContent = monthStart.toLocaleDateString("ru-RU", {
    month: "long",
    year: "numeric",
  });

  calendarGrid.innerHTML = "";
  for (let index = 0; index < 42; index += 1) {
    const current = new Date(gridStart);
    current.setDate(gridStart.getDate() + index);
    const isoDate = current.toISOString().slice(0, 10);
    const events = state.events.filter((event) => event.date === isoDate);
    const outsideMonth = current < monthStart || current > monthEnd;

    const day = document.createElement("article");
    day.className = `calendar-day${outsideMonth ? " outside" : ""}`;
    const items = events
      .map((event) => {
        const plantNames = event.plantIds
          .map((plantId) => state.plants.find((plant) => plant.id === plantId)?.name)
          .filter(Boolean)
          .join(", ");
        const actions = event.actions.map((item) => actionLabels[item] || item).join(", ");
        return `<div class="event-pill" title="${escapeHTML(event.notes || "")}">
          <strong>${escapeHTML(event.title)}</strong><br />
          ${escapeHTML(actions || "Без действий")}<br />
          ${escapeHTML(plantNames || "Без привязки")}
        </div>`;
      })
      .join("");

    day.innerHTML = `
      <header>
        <strong>${current.getDate()}</strong>
        <span class="caption">${events.length || ""}</span>
      </header>
      ${items || '<span class="muted">Нет событий</span>'}
    `;
    calendarGrid.append(day);
  }
}

async function createPlant(event) {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  await api("/api/plants", {
    method: "POST",
    body: JSON.stringify({
      name: form.get("name"),
      species: form.get("species"),
      room: form.get("room"),
      notes: form.get("notes"),
    }),
  });
  event.currentTarget.reset();
  await loadDashboard();
}

async function createEvent(event) {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  const payload = {
    title: form.get("title"),
    date: form.get("date"),
    actions: selectedValues(event.currentTarget.querySelector('select[name="actions"]')),
    plantIds: selectedValues(event.currentTarget.querySelector('select[name="plantIds"]')),
    notes: form.get("notes"),
  };
  await api("/api/events", {
    method: "POST",
    body: JSON.stringify(payload),
  });
  event.currentTarget.reset();
  await loadDashboard();
}

async function uploadPhoto(event, plantId) {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  const capturedAt = form.get("capturedAt");
  if (capturedAt) {
    form.set("capturedAt", new Date(capturedAt).toISOString());
  }
  await api(`/api/plants/${plantId}/photos`, {
    method: "POST",
    body: form,
    isMultipart: true,
  });
  event.currentTarget.reset();
  await loadDashboard();
}

async function deletePlant(plantId) {
  if (!window.confirm("Удалить растение из коллекции?")) {
    return;
  }
  await api(`/api/plants/${plantId}`, { method: "DELETE" });
  await loadDashboard();
}

async function api(path, options = {}) {
  if (config.authMode === "jwt" && state.refreshToken && Date.now() >= state.tokenExpiry - 15_000) {
    await refreshSession();
  }

  const headers =
    config.authMode === "oauth2-proxy"
      ? options.isMultipart
        ? {}
        : { "Content-Type": "application/json" }
      : options.isMultipart
        ? { Authorization: `Bearer ${state.accessToken}` }
        : {
            "Content-Type": "application/json",
            Authorization: `Bearer ${state.accessToken}`,
          };

  const response = await fetch(path, {
    ...options,
    headers: {
      ...headers,
      ...(options.headers || {}),
    },
  });

  if (response.status === 401) {
    if (config.authMode === "oauth2-proxy") {
      window.location.assign(config.oauth2ProxySignInUrl || "/oauth2/sign_in");
      return null;
    }

    clearSession();
    await startLogin();
    return null;
  }

  if (!response.ok) {
    const body = await response.text();
    throw new Error(body || "Ошибка запроса");
  }

  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    return response.json();
  }
  return response.text();
}

function selectedValues(select) {
  return Array.from(select.selectedOptions).map((item) => item.value);
}

function addMonths(date, delta) {
  return new Date(date.getFullYear(), date.getMonth() + delta, 1);
}

function randomString(length) {
  const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const bytes = crypto.getRandomValues(new Uint8Array(length));
  return Array.from(bytes, (value) => alphabet[value % alphabet.length]).join("");
}

async function base64UrlEncodeSha256(value) {
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(value));
  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/g, "");
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}
