# Plant Keeper

Локальный учебный проект для учета домашних растений:

- Go backend с REST API и файловым хранением данных
- frontend без сборщика, который backend отдает как статику
- отдельный Keycloak для OIDC
- `oauth2-proxy` как sidecar перед приложением
- Terraform для настройки realm и OIDC-клиента в Keycloak

## Архитектура

- `docker-compose.keycloak.yml` поднимает отдельный `Keycloak + PostgreSQL`
- `docker-compose.yml` поднимает приложение:
  - `backend` слушает только внутри docker-сети
  - `oauth2-proxy` публикуется наружу на `http://localhost:8080`
  - все запросы в приложение проходят через `oauth2-proxy`
- `oauth2-proxy` ходит в Keycloak по issuer `http://host.docker.internal:8081/realms/plants`

## Что меняется по сравнению с прямым OIDC во frontend

- frontend больше не обменивает `code` на токены сам
- браузер работает с приложением через cookie-сессию `oauth2-proxy`
- backend доверяет identity headers от `oauth2-proxy`
- Keycloak-клиент теперь `CONFIDENTIAL`, а не public SPA client

## Локальный запуск

### 1. Поднять отдельный Keycloak

```powershell
docker compose -f docker-compose.keycloak.yml up -d
```

Keycloak будет доступен на `http://localhost:8081`.

### 2. Настроить Keycloak через Terraform

```powershell
cd infra/terraform
terraform init
terraform apply
```

Terraform создаст:

- realm `plants`
- роль `app-user`
- confidential client `plant-keeper-proxy`

Для локального dev-сценария client secret по умолчанию:

```text
dev-oauth2-proxy-secret
```

### 3. Поднять приложение с sidecar proxy

```powershell
cd ..
docker compose up --build -d
```

После этого приложение будет открываться на `http://localhost:8080`.

## Compose-файлы

- [docker-compose.keycloak.yml](./docker-compose.keycloak.yml): отдельный Keycloak и PostgreSQL
- [docker-compose.yml](./docker-compose.yml): backend и `oauth2-proxy`

## Основные переменные

В backend:

- `AUTH_MODE=oauth2-proxy`
- `OAUTH2_PROXY_SIGN_IN_URL=/oauth2/sign_in`
- `OAUTH2_PROXY_SIGN_OUT_URL=/oauth2/sign_out?rd=%2F`

В Terraform:

- `frontend_base_url = "http://localhost:8080"`
- `oauth2_proxy_client_id = "plant-keeper-proxy"`
- `oauth2_proxy_client_secret = "dev-oauth2-proxy-secret"`

## Если нужно вернуться к прямой JWT-проверке

Backend все еще поддерживает старый режим:

```env
AUTH_MODE=jwt
KEYCLOAK_ISSUER_URL=http://localhost:8081/realms/plants
KEYCLOAK_CLIENT_ID=plant-keeper-proxy
```

Но в текущей локальной схеме основным считается режим через `oauth2-proxy`.
