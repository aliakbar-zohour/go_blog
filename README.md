# Go Blog API

A REST API for a blog built with Go, PostgreSQL, and support for image and video uploads.

---

## Table of contents

- [Requirements](#requirements)
- [Project structure](#project-structure)
- [Quick start with Docker](#quick-start-with-docker)
- [Local development (without Docker)](#local-development-without-docker)
- [Environment variables](#environment-variables)
- [API endpoints](#api-endpoints)
- [Testing with Swagger UI](#testing-with-swagger-ui)
- [Regenerating Swagger docs](#regenerating-swagger-docs)
- [Building and running the binary](#building-and-running-the-binary)
- [License](#license)

---

## Requirements

- **Go 1.21+** (for local run and building)
- **Docker and Docker Compose** (optional; for running app + PostgreSQL in containers)
- **PostgreSQL 14+** (if running locally without Docker)

---

## Project structure

| Path | Description |
|------|-------------|
| `cmd/api/main.go` | Application entry point; loads config, DB, services, starts HTTP server |
| `internal/config` | Loads settings from environment and defaults |
| `internal/database` | PostgreSQL connection and auto-migration |
| `internal/model` | Post, Media, Author, Category, Comment models |
| `internal/repository` | Data access layer (CRUD for posts, media, authors, categories, comments) |
| `internal/service` | Business logic and validation |
| `internal/handler` | HTTP handlers for API routes |
| `internal/router` | Route definitions and middleware wiring |
| `internal/middleware` | Panic recovery, security headers, request logging |
| `internal/upload` | File type/size validation and safe storage |
| `pkg/response` | Shared JSON response helpers |
| `docs/` | Generated Swagger/OpenAPI docs (do not edit by hand) |

---

## Quick start with Docker

This is the simplest way to run the API and database together.

### 1. Clone and enter the project

```bash
git clone https://github.com/aliakbar-zohour/go_blog.git
cd go_blog
```

### 2. Build and start services

```bash
docker compose up --build
```

- **API** will be available at **http://localhost:8080**
- **PostgreSQL** runs in a container; no local install needed.
- Uploaded files are stored in a Docker volume.

### 3. Stop services

```bash
docker compose down
```

To also remove the database volume:

```bash
docker compose down -v
```

---

## Local development (without Docker)

Use this when you want to run the API on your machine and use a local (or remote) PostgreSQL instance.

### 1. Install and run PostgreSQL

- Install PostgreSQL 14+ and ensure it is running.
- Create a database, for example: `createdb go_blog`

### 2. Clone the project

```bash
git clone https://github.com/aliakbar-zohour/go_blog.git
cd go_blog
```

### 3. Copy environment file and set variables

```bash
cp .env.example .env
```

Edit `.env` and set at least:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

### 4. Install dependencies and run the API

```bash
go mod download
go run ./cmd/api
```

The server listens on **http://localhost:8080** (or the port set in `PORT`).

### 5. (Optional) Create uploads directory

If you use file uploads and keep the default `UPLOAD_DIR=uploads`, ensure the directory exists or the app will create it on first upload.

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
| `GET` | `/api/posts` | List posts (query: `limit`, `offset`, `category_id` to filter by category) |
| `POST` | `/api/posts` | Create post (form: `title`, `body`, `author_id`, `category_id`, `banner`, `files[]`) |
| `GET` | `/api/posts/:id` | Get one post by ID (includes author and category) |
| `PUT` | `/api/posts/:id` | Update post (form: `title`, `body`, `author_id`, `category_id`, `banner`, `files[]`; all optional) |
| `DELETE` | `/api/posts/:id` | Delete post (soft delete) |

### Authors

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/authors` | List all authors |
| `POST` | `/api/authors` | Create author (form: `name`, `avatar`) |
| `GET` | `/api/authors/:id` | Get one author |
| `PUT` | `/api/authors/:id` | Update author (form: `name`, `avatar`) |
| `DELETE` | `/api/authors/:id` | Delete author |

### Categories

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/categories` | List all categories |
| `POST` | `/api/categories` | Create category (form: `name`) |
| `GET` | `/api/categories/:id` | Get one category |
| `PUT` | `/api/categories/:id` | Update category (form: `name`) |
| `DELETE` | `/api/categories/:id` | Delete category |

### Comments

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/posts/:postId/comments` | List comments for a post |
| `POST` | `/api/posts/:postId/comments` | Create comment (form: `body`, `author_name`) |
| `PUT` | `/api/comments/:id` | Update comment (form: `body`) |
| `DELETE` | `/api/comments/:id` | Delete comment |

- Uploaded files are served under **`/uploads/<path>`** (e.g. `/uploads/posts/1/xyz.jpg`, `/uploads/banners/...`, `/uploads/avatars/...`).
- **Allowed image extensions:** jpg, jpeg, png, gif, webp  
- **Allowed video extensions:** mp4, webm, mov  

All JSON responses use a common shape: `{ "success": true|false, "data": ..., "error": "..." }`.

---

## Testing with Swagger UI

After the server is running (Docker or local):

1. Open in a browser: **http://localhost:8080/docs/index.html**
2. Use **Try it out** on any endpoint, fill parameters, then **Execute**.
3. The response (including status and body) is shown on the same page; if the API returns an error, the `error` field and status code are visible there.

---

## Regenerating Swagger docs

If you change Swagger comments in `cmd/api/main.go` or `internal/handler/post_handler.go`, regenerate the docs:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -d . -o docs
```

Then rebuild or restart the app so the new docs are served.

---

## Building and running the binary

Build:

```bash
go build -o api ./cmd/api
```

Run (set env vars or use `.env`):

```bash
./api
```

On Windows:

```bash
go build -o api.exe ./cmd/api
api.exe
```

---

## License

MIT
