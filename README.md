# Shop API

A REST API for an online shop, built in Go with a clean, layered architecture. This is a learning project focused on writing idiomatic, well-structured backend Go — not a finished storefront yet, but the foundation (auth, users, categories) is production-quality.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-008ECF)
![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?logo=mysql&logoColor=white)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Authentication Flow](#authentication-flow)
- [Roadmap](#roadmap)
- [License](#license)

## Features

- **Clean architecture** — strict separation between `handler` → `validator` → `service` → `repository` layers, with each layer depending only on interfaces of the layer below it.
- **OTP + JWT authentication** — phone-number login via one-time codes (generated with `crypto/rand`, not `math/rand`), backed by short-lived access tokens and long-lived refresh tokens.
- **Real session revocation** — refresh tokens are tied to a session record in the database, so logout actually invalidates a session instead of just discarding a token client-side.
- **Role-based access control** — middleware-enforced roles (e.g. only admins can create categories).
- **Storage-agnostic image pipeline** — uploaded images are validated by content sniffing (not just file extension), auto-oriented, resized, and re-encoded before being handed to a pluggable `Storage` interface (local disk today, swappable for S3 or similar later without touching business logic).
- **Structured error handling** — a custom `richerror` package carries operation, message, kind, and metadata through every layer, and maps cleanly to HTTP status codes and a consistent JSON response envelope.
- **Configuration via `.env` + YAML** — secrets and per-environment values load from `.env`, general settings from `config.yml`, merged with [koanf](https://github.com/knadh/koanf).
- **Migrations built into the binary** — schema migrations run automatically on startup via a small wrapper around [golang-migrate](https://github.com/golang-migrate/migrate), with the standalone CLI still available for manual up/down.

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.26 |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | MySQL 8.4 (`database/sql` + `go-sql-driver/mysql`, no ORM) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Config | [koanf](https://github.com/knadh/koanf) (`.env` + YAML) |
| Auth | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) + bcrypt (`golang.org/x/crypto`) |
| Validation | [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) |
| Image processing | [disintegration/imaging](https://github.com/disintegration/imaging) |
| Local dev infra | Docker Compose (MySQL + phpMyAdmin) |

## Architecture

Every request flows through the same pipeline, regardless of module:

```
HTTP request
   │
   ▼
Router (gin route groups)
   │
   ▼
Middleware (auth / role checks)
   │
   ▼
Handler        — binds & parses the request, delegates everything else
   │
   ▼
Validator      — business-rule and input validation (ozzo-validation)
   │
   ▼
Service        — business logic, orchestrates repositories
   │
   ▼
Repository     — raw SQL against MySQL, maps rows to entities
```

Errors flow back up the same pipeline as a `richerror.RichError` (carrying an operation name, a message, a `Kind`, and optional field-level metadata), which a shared `mapper` package converts to the right HTTP status code, and a shared `response` package serializes into a single consistent JSON envelope:

```json
{
  "code": 200,
  "message": "human readable message",
  "data": { },
  "errors": { }
}
```

## Project Structure

```
shop/
├── cmd/                      # composition root — wires repositories, services, handlers per module
├── internal/
│   ├── api/
│   │   ├── handler/          # one package per module (auth, user, category, health)
│   │   ├── middleware/       # auth & role-based middleware
│   │   ├── router/           # route registration
│   │   └── server/           # HTTP server setup
│   ├── config/                # .env + YAML config loading (koanf)
│   ├── dto/                   # request/response shapes per module
│   ├── entity/                # domain models
│   ├── migrator/              # golang-migrate wrapper, runs on startup
│   ├── pkg/
│   │   ├── claims/            # JWT access/refresh token creation & parsing
│   │   ├── imageprocessor/    # upload validation, resize/encode, storage-agnostic save
│   │   ├── mapper/             # richerror.Kind -> HTTP status code
│   │   ├── response/           # consistent JSON response envelope
│   │   └── richerror/          # structured application error type
│   ├── repository/mysql/      # one package per module, raw SQL
│   ├── service/               # one package per module, business logic
│   └── validator/             # one package per module, input validation
├── migrations/                 # golang-migrate SQL files
├── uploads/                     # locally stored uploaded images (served at /uploads)
├── config.yml
├── docker-compose.yml
└── plans.md                     # module/endpoint checklist
```

## Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (optional — only needed for manual migrations; the app runs pending migrations automatically on startup)

### 1. Clone and install dependencies

```bash
git clone <this-repo-url>
cd shop
go mod download
```

### 2. Start MySQL (and phpMyAdmin)

```bash
docker compose up -d
```

This starts MySQL 8.4 on `localhost:3308` and phpMyAdmin on `localhost:8080`.

### 3. Configure environment

Create a `.env` file in the project root (see [Configuration](#configuration) for the full list of required keys) and adjust `config.yml` if needed.

### 4. Run the server

```bash
go run ./cmd
```

Pending migrations run automatically on startup. The API is now available at `http://localhost:3000` (or whatever `server.host`/`server.port` you configured).

```bash
curl http://localhost:3000/health-check
```

## Configuration

Settings are split between `config.yml` (non-secret, checked into git) and `.env` (secrets, gitignored). `.env` is loaded first, `config.yml` second; env var names map to config keys by lowercasing and turning `__` into `.` (e.g. `MYSQL__HOST` → `mysql.host`).

**`.env`** — required keys:

| Key | Description |
|---|---|
| `MYSQL__HOST` | MySQL host (`localhost` for the Docker Compose setup) |
| `MYSQL__PORT` | MySQL port (`3308` for the Docker Compose setup) |
| `MYSQL__USERNAME` | MySQL user |
| `MYSQL__PASSWORD` | MySQL password |
| `MYSQL__DATABASE` | Database name |
| `AUTH_SERVICE__ACCESS_TOKEN_SECRET` | Secret used to sign JWT access tokens |
| `AUTH_SERVICE__REFRESH_TOKEN_SECRET` | Secret used to sign JWT refresh tokens |

**`config.yml`** — non-secret settings (connection pool sizes, server timeouts, token durations, upload limits); see the file itself for the current values.

### Manual migrations (optional)

The app applies pending migrations automatically on startup, but the CLI is available for manual control:

```bash
# up
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" up

# down
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" down
```

## API Reference

All responses use the envelope shown in [Architecture](#architecture). Endpoints marked 🔒 require a valid `Authorization: Bearer <access-token>` header; 🔒👑 additionally requires the `ADMIN` role.

### Health

| Method | Path | Description |
|---|---|---|
| GET | `/health-check` | Liveness check |

### Auth

| Method | Path | Description |
|---|---|---|
| POST | `/auth/send-otp` | Request a one-time code for a phone number |
| POST | `/auth/check-otp` | Verify the code, receive an access token (and a refresh token cookie) |
| GET | `/auth/me` | 🔒 Get the current session's user id |
| POST | `/auth/refresh-token` | Exchange a valid refresh token for a new access token |
| POST | `/auth/logout` | 🔒 Revoke the current session |

### User

| Method | Path | Description |
|---|---|---|
| GET | `/user/profile` | 🔒 Get the current user's profile |
| PATCH | `/user/update-profile` | 🔒 Update name and/or avatar |
| PATCH | `/user/change-password` | 🔒 Set a new password |

### Category

| Method | Path | Description |
|---|---|---|
| POST | `/category` | 🔒👑 Create a category (root or one level of children) |
| GET | `/category` | List all categories as a nested tree |
| GET | `/category/:slug` | Get a single category by slug, including its children if it's a root category |
| PATCH | `/category/:slug` | 🔒👑 Partially update a category (title/slug/image — only fields provided are changed) |
| DELETE | `/category/:slug` | 🔒👑 Soft-delete a category (rejected if a root category still has children) |

### Province

| Method | Path | Description |
|---|---|---|
| GET | `/province` | List all provinces |
| GET | `/province/:id` | Get a single province by id |

### Address

| Method | Path | Description |
|---|---|---|
| POST | `/address` | 🔒 Create a new shipping address |
| GET | `/address` | 🔒 List the current user's addresses |
| GET | `/address/:id` | 🔒 Get a single address by id |
| PATCH | `/address/:id` | 🔒 Partially update an address (only fields provided are changed) |
| DELETE | `/address/:id` | 🔒 Soft-delete an address |

### Product

| Method | Path | Description |
|---|---|---|
| POST | `/product` | 🔒👑 Create a product (main image + gallery images, price, optional stock, category) |
| GET | `/product/:slug` | Get a single product by slug, including its gallery images |

### Static files

| Method | Path | Description |
|---|---|---|
| GET | `/uploads/*` | Serves uploaded images (avatars, category images) |

## Authentication Flow

1. **`POST /auth/send-otp`** — client submits a phone number. The server finds or creates a user for that number and generates a 5-digit code with `crypto/rand`, valid for a short window (`otp_code_duration`).
2. **`POST /auth/check-otp`** — client submits the phone number and code. On success, the server creates a **session record** in the database and issues a short-lived **access token** and a longer-lived **refresh token**, both JWTs carrying the session id.
3. Protected routes go through `AuthRequired` middleware, which parses the access token, loads the referenced session, and rejects the request if the session has expired or been revoked.
4. **`POST /auth/refresh-token`** — when the access token expires, the client uses the refresh token (sent as an httpOnly cookie) to get a new access token, as long as the underlying session is still valid.
5. **`POST /auth/logout`** — revokes the session server-side, so both existing and any future access/refresh tokens issued for it stop working immediately. This is the key advantage over plain stateless JWTs, where logout can only ever be simulated client-side.

## Roadmap

Tracked in detail in [`plans.md`](./plans.md). Foundational modules (config, error handling, migrations, HTTP server) and the auth/user/category/province/address modules above are done. Still to build, in rough priority order:

- [ ] **Product** (in progress — create and get-by-slug are done) — remaining: listing with filters/search/pagination, update, delete
- [ ] **Cart** — add/update/remove items, view current cart
- [ ] **Order** — checkout from cart, order history, status transitions
- [ ] **Payment** — payment initiation and gateway callback handling
- [ ] Nice-to-have: reviews, wishlist, coupons

## License

No license has been chosen yet — all rights reserved by default until one is added.