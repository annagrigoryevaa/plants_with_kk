# Plant Keeper

Полноценный стартовый проект для ведения домашней коллекции растений:

- Go-бекенд с REST API, хранением каталога и загрузкой фотографий
- фронтенд без сборщика, который работает как SPA
- фотопленка по каждому растению для отслеживания роста
- календарь ухода с мультивыбором действий: `полив`, `пересадка`, `подкормка`
- привязка календарных событий к нескольким растениям
- обязательная авторизация через Keycloak по OIDC Authorization Code + PKCE
- Terraform-конфиг для управления Keycloak realm и клиентом

## Структура

- [backend/main.go](/Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/backend/main.go)
- [frontend/public/index.html](/Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/frontend/public/index.html)
- [infra/terraform/main.tf](/Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/infra/terraform/main.tf)
- [docker-compose.yml](/Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/docker-compose.yml)

## Что уже реализовано

### Бекенд

- каталог растений: создание, обновление, удаление, получение списка
- загрузка фотографий в файловую систему и возврат ссылок для фронтенда
- календарные события ухода с несколькими действиями и несколькими растениями
- файловое JSON-хранилище для локального запуска без внешней БД приложения
- обязательная проверка JWT access token от Keycloak по JWKS

### Фронтенд

- автоматический логин через Keycloak
- карточки растений с фотопленкой
- форма добавления событий в календарь с мультивыбором действий и растений
- помесячный календарь ухода

### Инфраструктура

- Docker Compose для локального Keycloak + PostgreSQL
- Terraform для создания realm `plants`, роли `app-user` и клиента `plant-keeper-web`

## Локальный запуск

### 1. Поднять Keycloak

```bash
cd /Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform
docker compose up -d
```

Keycloak будет доступен по адресу [http://localhost:8081](http://localhost:8081).

### 2. Применить Terraform

На текущей машине `terraform` не установлен, поэтому сначала нужен Terraform CLI.

```bash
cd /Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/infra/terraform
cp terraform.tfvars.example terraform.tfvars
terraform init
terraform apply
```

После `apply` Terraform создаст realm и клиент для приложения.

### 3. Настроить бекенд

```bash
cd /Users/annagrigoreva/Documents/Codex/2026-04-17-go-terraform/backend
cp .env.example .env
set -a
source .env
set +a
go run .
```

Приложение будет доступно по адресу [http://localhost:8080](http://localhost:8080).

## Как работает авторизация

- фронтенд начинает OIDC Authorization Code Flow с PKCE
- Keycloak возвращает `code`
- фронтенд обменивает `code` на токены через `token endpoint`
- Go-бекенд принимает `Bearer` access token и валидирует подпись через JWKS Keycloak
- без токена API не отвечает

## Ограничения текущего стартового проекта

- приложение хранит данные в JSON-файле, а не в PostgreSQL или SQLite
- календарные события пока можно создавать через форму, но отдельный UI для редактирования событий еще не добавлен
- для Terraform я подготовил конфиг, но не смог локально выполнить `terraform init/apply`, потому что CLI отсутствует
