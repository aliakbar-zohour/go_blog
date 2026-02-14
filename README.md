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
- [License](#license)

---

## Features

- **Posts** – CRUD, banner image, optional author and category, extra media (images/videos)
- **Authors** – CRUD with name and avatar image
- **Categories** – CRUD; filter posts by category
- **Comments** – List/create per post; update/delete by comment ID
- **Swagger UI** – Interactive API docs at `/docs/` (generated from code in Docker)
- **File uploads** – Banners, avatars, post media; served under `/uploads/`

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
| `internal/middleware` | Panic recovery, security headers, logging |
| `internal/upload` | File validation and storage (banners, avatars, media) |
| `pkg/response` | Shared JSON response format |
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

---

## API endpoints

### Posts

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/posts` | List posts (`limit`, `offset`, `category_id` to filter) |
| `POST` | `/api/posts` | Create (form: `title`, `body`, `author_id`, `category_id`, `banner`, `files[]`) |
| `GET` | `/api/posts/:id` | Get one (includes author and category) |
| `PUT` | `/api/posts/:id` | Update (form: same as create; all optional) |
| `DELETE` | `/api/posts/:id` | Soft delete |

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
- **Response shape:** `{ "success": true|false, "data": ..., "error": "..." }`.

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
