# Go Blog API

REST API for a blog built with **Go**, **PostgreSQL**, and **Swagger UI**. Supports posts (with banner, author, category), authors, categories, comments, and image/video uploads.

---

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Project structure](#project-structure)
- [Quick start with Docker](#quick-start-with-docker)
- [Local development (without Docker)](#local-development-without-docker)
- [Environment variables](#environment-variables)
- [API endpoints](#api-endpoints)
- [Swagger UI](#swagger-ui)
- [Regenerating Swagger docs (local)](#regenerating-swagger-docs-local)
- [Building and running the binary](#building-and-running-the-binary)
- [Testing](#testing)
- [License](#license)

---

## Features

- **Auth** – Register with email (verification code), verify & complete profile, login with email/password; JWT for protected routes
- **Posts** – CRUD with banner, category, media; **create/update/delete require JWT** (author = logged-in writer)
- **Authors** – List/create/update/delete (name, avatar); registered writers have email and can log in
- **Categories** – CRUD; filter posts by category
- **Comments** – List/create per post; update/delete by comment ID
- **Swagger UI** – Interactive API docs at `/docs/` (generated from code in Docker)
- **File uploads** – Banners, avatars, post media; served under `/uploads/`
- **Health check** – `GET /health` for load balancers and orchestration (checks DB when configured)
- **Pagination** – List posts returns `{ "items": [...], "total": N }` for proper paging
- **Request ID** – Every response includes `X-Request-ID`; structured logging uses it
- **Error codes** – Auth and protected routes return a `code` field (e.g. `invalid_token`, `auth_required`) for clients

---

## Requirements

- **Go 1.23+** (for local run and building)
- **Docker and Docker Compose** (optional; runs API + PostgreSQL)
- **PostgreSQL 14+** (if running without Docker)

---

## Project structure

| Path | Description |
|------|-------------|
| `cmd/api/main.go` | Entry point; config, DB, services, HTTP server |
| `internal/config` | Settings from environment and defaults |
| `internal/database` | PostgreSQL connection and auto-migration |
| `internal/model` | Post, Media, Author, Category, Comment |
| `internal/repository` | Data access (CRUD for all entities) |
| `internal/service` | Business logic and validation |
| `internal/handler` | HTTP handlers and Swagger annotations |
| `internal/router` | Routes and middleware |
| `internal/middleware` | Panic recovery, security headers, logging, JWT auth |
| `internal/mail` | Sends verification code email (HTML template) |
| `internal/upload` | File validation and storage (banners, avatars, media) |
| `pkg/response` | Shared JSON response format |
| `pkg/auth` | Password hashing (bcrypt), JWT create/parse |
| `docs/` | Generated Swagger (by `swag init` or inside Docker) |

---

## Quick start with Docker

Run API and PostgreSQL with one command.

### 1. Clone and enter the project

```bash
git clone https://github.com/aliakbar-zohour/go_blog.git
cd go_blog
```

### 2. Build and start

```bash
docker compose up --build -d
```

- **API:** **http://localhost:8080**
- **Swagger UI:** **http://localhost:8080/docs/** or **http://localhost:8080/docs/index.html**
- PostgreSQL and uploads use Docker volumes.

Swagger is generated **inside the image** from the handler source (the `docs/` folder on your machine is not used). After code or Swagger comment changes, run again:

```bash
docker compose up --build -d
```

### 3. Stop

```bash
docker compose down
```

Remove database and uploads volumes:

```bash
docker compose down -v
```

---

## Local development (without Docker)

### 1. PostgreSQL

Install PostgreSQL 14+, create a database:

```bash
createdb go_blog
```

### 2. Clone and env

```bash
git clone https://github.com/aliakbar-zohour/go_blog.git
cd go_blog
cp .env.example .env
```

Edit `.env`: set `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`.

### 3. Run

```bash
go mod download
go run ./cmd/api
```

Server: **http://localhost:8080**. Create an `uploads` directory if you use file uploads (or rely on auto-creation).

---

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `go_blog` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `UPLOAD_DIR` | `uploads` | Directory for uploaded files |
| `MAX_UPLOAD_MB` | `50` | Max file size per upload (MB) |
| `JWT_SECRET` | `change-me-in-production` | Secret for signing JWTs (set in production) |
| `JWT_EXPIRY_HOURS` | `72` | JWT expiry in hours |
| `SMTP_HOST` | (empty) | SMTP server for verification emails; if empty, codes are not sent |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_USER` | (empty) | SMTP username |
| `SMTP_PASS` | (empty) | SMTP password |
| `SMTP_FROM` | `noreply@go-blog.local` | From address for emails |
| `CORS_ORIGINS` | `*` | Comma-separated allowed origins (e.g. `https://app.example.com`) |
| `BODY_LIMIT_BYTES` | `33554432` (32MB) | Max request body size; 413 if exceeded |
| `AUTH_RATE_PER_MIN` | `10` | Max auth requests per IP per minute (login/register) |

---

## Security & performance

- **Rate limiting** – Auth endpoints (`/api/auth/*`) are limited per IP to reduce brute-force and abuse.
- **Body limit** – Request body size is capped; oversized requests get 413.
- **CORS** – Configure `CORS_ORIGINS` in production to allow only your front-end origin(s).
- **Headers** – `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Referrer-Policy` are set.
- **DB pool** – Connection pool (max open 25, idle 5) and prepared statements for faster queries.
- **Gzip** – JSON responses are compressed when the client sends `Accept-Encoding: gzip`.
- **Input validation** – Email format, field lengths (title, category name, comment body), and file type/size are validated.
- **Request ID** – Each request gets a unique `X-Request-ID` header in the response; logs are structured (slog) with `request_id`, `method`, `path`, `status`, `duration_ms`.

---

## API endpoints

### Health

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check; returns `{ "status": "ok" }` and 503 if DB is unavailable |

---

### Auth (no JWT required)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/auth/register/request` | Request verification code (body: `{"email":"..."}`); sends code to email if SMTP configured |
| `POST` | `/api/auth/register/verify` | Verify code and complete registration (body: `email`, `code`, `name`, `password`); returns `author` + `token` |
| `POST` | `/api/auth/login` | Login (body: `email`, `password`); returns `author` + `token` |

Use the `token` in the **Authorization** header: `Authorization: Bearer <token>` for protected routes.

**Registration flow:**  
1. `POST /api/auth/register/request` with `{"email":"writer@example.com"}` → a 6-digit code is generated. If **SMTP is not configured**, the response includes `dev_code` (use it in step 2). If SMTP is set, the code is sent by email.  
2. `POST /api/auth/register/verify` with `{"email":"...", "code":"<dev_code or from email>", "name":"Jane", "password":"secret123"}` → account is created and a JWT is returned.  
3. Use the JWT in `Authorization: Bearer <token>` when creating or editing posts.

**Sending real emails:** Set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, and `SMTP_FROM` in your env (or `.env`). For Gmail use an [App Password](https://support.google.com/accounts/answer/185833) and `SMTP_HOST=smtp.gmail.com`, `SMTP_PORT=587`. For testing, you can use [Mailtrap](https://mailtrap.io) or similar.

### Posts (create/update/delete require JWT; author = logged-in writer)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/posts` | List posts; returns `{ "items": [...], "total": N }`. Query: `limit`, `offset`, `category_id` |
| `POST` | `/api/posts` | **Auth.** Create (form: `title`, `body`, `category_id`, `banner`, `files[]`); author set from JWT |
| `GET` | `/api/posts/:id` | Get one (includes author and category) |
| `PUT` | `/api/posts/:id` | **Auth.** Update own post (form: `title`, `body`, `category_id`, `banner`, `files[]`) |
| `DELETE` | `/api/posts/:id` | **Auth.** Delete own post |

### Authors

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/authors` | List all |
| `POST` | `/api/authors` | Create (form: `name`, `avatar`) |
| `GET` | `/api/authors/:id` | Get one |
| `PUT` | `/api/authors/:id` | Update (form: `name`, `avatar`) |
| `DELETE` | `/api/authors/:id` | Delete |

### Categories

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/categories` | List all |
| `POST` | `/api/categories` | Create (form: `name`) |
| `GET` | `/api/categories/:id` | Get one |
| `PUT` | `/api/categories/:id` | Update (form: `name`) |
| `DELETE` | `/api/categories/:id` | Delete |

### Comments

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/posts/:postId/comments` | List comments for a post |
| `POST` | `/api/posts/:postId/comments` | Create (form: `body`, `author_name`) |
| `PUT` | `/api/comments/:id` | Update (form: `body`) |
| `DELETE` | `/api/comments/:id` | Delete |

- **Static files:** `/uploads/<path>` (e.g. `/uploads/posts/1/xyz.jpg`, `/uploads/banners/...`, `/uploads/avatars/...`).
- **Images:** jpg, jpeg, png, gif, webp. **Videos:** mp4, webm, mov.
- **Response shape:** `{ "success": true|false, "data": ..., "error": "...", "code": "..." }`. The `code` field is set on errors (e.g. `invalid_credentials`, `auth_required`) for machine-readable handling.

---

## Testing

Run unit tests:

```bash
go test ./pkg/... ./internal/... -v
```

- **pkg/auth** – Password hashing and JWT (no DB).
- **internal/handler** – Health handler (no DB).
- **internal/service** – Post list/count tests use SQLite in-memory and **skip when CGO is disabled** (e.g. default Windows build).

Integration test (requires PostgreSQL; skips if DB is not available):

```bash
go test -v -run TestIntegration .
```

---

## Swagger UI

With the server running (Docker or local):

1. Open **http://localhost:8080/docs/** or **http://localhost:8080/docs/index.html**
2. Use **Try it out** on any endpoint, set parameters, then **Execute**
3. Response body and status (including errors) are shown on the page

In Docker, Swagger is generated at build time from the handler code; no need to run `swag init` on your machine for the container.

---

## Regenerating Swagger docs (local)

When you change Swagger comments in `cmd/api/main.go` or any handler under `internal/handler/`, regenerate docs locally:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -d . -o docs
```

Then restart the app (or rebuild the Docker image with `docker compose up --build -d`).

---

## Building and running the binary

```bash
go build -o api ./cmd/api
./api
```

Windows:

```bash
go build -o api.exe ./cmd/api
api.exe
```

---

## License

MIT
