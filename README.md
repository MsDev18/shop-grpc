# Shop API

A backend for an online shop, built in Go with a clean, layered architecture. This is a
learning project focused on writing idiomatic, well-structured backend Go — not a finished
storefront yet, but the foundation (auth, users, categories, products, addresses) is
production-quality. It was migrated from REST to **gRPC-only**, module by module, using the
strangler-fig pattern (REST and gRPC ran side by side until every module had a gRPC
equivalent, then the REST/Gin layer was deleted entirely). The next phase, now starting, is
splitting this single monolith process into independently deployable microservices, one
module at a time.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-only-4285F4?logo=grpc&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?logo=mysql&logoColor=white)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Migration History](#migration-history)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference (gRPC)](#api-reference-grpc)
- [Authentication Flow](#authentication-flow)
- [Roadmap](#roadmap)
- [License](#license)

## Features

- **Clean architecture** — strict separation between `grpc server` → `validator` → `service`
  → `repository` layers, with each layer depending only on interfaces of the layer below it.
- **OTP + JWT authentication** — phone-number login via one-time codes (generated with
  `crypto/rand`, not `math/rand`), backed by short-lived access tokens and long-lived refresh
  tokens.
- **Real session revocation** — refresh tokens are tied to a session record in the database,
  so logout actually invalidates a session instead of just discarding a token client-side.
  Enforced globally by a `grpc.UnaryServerInterceptor` / `grpc.StreamServerInterceptor` pair.
- **Role-based access control** — split into two composable interceptors, chained via
  `grpc.ChainUnaryInterceptor` / `grpc.ChainStreamInterceptor`: an authentication interceptor
  (is this caller logged in?) and a separate authorization interceptor (does this method
  require a specific role, and does the caller have it?), driven by a declarative
  `requiredRoles` map keyed by full gRPC method name.
- **Storage-agnostic image pipeline** — uploaded images are validated by content sniffing
  (not just file extension), auto-oriented, resized, and re-encoded before being handed to a
  pluggable `Storage` interface (local disk today, swappable for S3 or similar later without
  touching business logic). Image bytes travel as raw `bytes` proto fields, and product
  gallery images use client-streaming so multiple images can be sent without one giant
  message.
- **Structured error handling** — a custom `richerror` package carries operation, message,
  kind, and metadata through every layer. A `mapper` package turns a `richerror.Kind` into a
  `google.golang.org/grpc/codes.Code`, so every module reports errors consistently.
- **Configuration via `.env` + YAML** — secrets and per-environment values load from `.env`,
  general settings from `config.yml`, merged with [koanf](https://github.com/knadh/koanf).
- **Migrations built into the binary** — schema migrations run automatically on startup via a
  small wrapper around [golang-migrate](https://github.com/golang-migrate/migrate), with the
  standalone CLI still available for manual up/down.

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.26 |
| RPC framework | [gRPC-Go](https://github.com/grpc/grpc-go) + [Protocol Buffers](https://protobuf.dev/) (`protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`) |
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
 gRPC request
     │
     ▼
 Auth interceptor (authentication — is this caller logged in?)
     │
     ▼
 Role interceptor (authorization — does this method require a role, and does the caller have it?)
     │
     ▼
 gRPC server method (per module)
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

Errors flow back up the same pipeline as a `richerror.RichError` (carrying an operation name,
a message, a `Kind`, and optional field-level metadata). The shared `mapper` package converts
that into a `google.golang.org/grpc/codes.Code` wrapped in a `status.Error`, so a gRPC client
sees a proper gRPC status (e.g. `NotFound`, `PermissionDenied`, `Unauthenticated`) instead of
a generic failure.

## Migration History

This project used to be REST-only (Gin), then went through two migrations:

1. **REST → gRPC** (done), via the strangler fig pattern: a gRPC server ran on `:50051` in a
   goroutine alongside the existing Gin server on `:3000`, and each module was migrated one at
   a time, both protocols sharing the same `validator` / `service` / `repository` layers and
   the same MySQL database. Once every module had a gRPC equivalent, the entire REST/Gin
   layer (`handler`, `middleware`, `router`, `server` packages) was deleted. The project is
   now gRPC-only, on `:50051`, with [server reflection](https://github.com/grpc/grpc-go/tree/master/reflection)
   enabled so tools like `grpcurl` or Postman/SpecHub can discover services without needing
   the `.proto` files locally.
2. **Monolith → microservices** (starting now), gradual and one module at a time — see
   [`plans.md`](./plans.md) for the agreed approach (each extracted service gets its own
   database from day one; every other module's in-process call into it gets rewired into a
   real gRPC client call).

`.proto` files live under `proto/<module>/`, one file per module, and generate into
`internal/pb/<module>/` via `protoc` — see [Getting Started](#getting-started) for the exact
command. Generated files are committed but never hand-edited.

Current per-module status is tracked in [`plans.md`](./plans.md); the short version is in
[API Reference (gRPC)](#api-reference-grpc) below.

## Project Structure

```
shop/
├── cmd/                      # composition root — wires repositories, services, and gRPC
│                              # servers per module, plus the single gRPC server entrypoint
├── proto/                    # .proto source files, one folder per module (source of truth)
├── internal/
│   ├── api/
│   │   └── grpc/              # gRPC servers — one package per module, plus the auth package
│   │                          # (authentication interceptor, authorization/role interceptor)
│   │                          # and the top-level Server that registers every service and
│   │                          # starts listening on :50051
│   ├── config/                # .env + YAML config loading (koanf)
│   ├── dto/                   # request/response shapes per module
│   ├── entity/                # domain models
│   ├── migrator/              # golang-migrate wrapper, runs on startup
│   ├── pb/                    # generated Go code from proto/ — never edited by hand
│   ├── pkg/
│   │   ├── claims/            # JWT access/refresh token creation & parsing
│   │   ├── imageprocessor/    # upload validation, resize/encode, storage-agnostic save
│   │   ├── mapper/             # richerror.Kind -> grpc/codes.Code
│   │   └── richerror/          # structured application error type
│   ├── repository/mysql/      # one package per module, raw SQL
│   ├── service/               # one package per module, business logic
│   └── validator/             # one package per module, input validation
├── migrations/                 # golang-migrate SQL files
├── uploads/                     # locally stored uploaded images
├── config.yml
├── docker-compose.yml
└── plans.md                     # module/endpoint checklist, migration status
```

> **Known gap:** uploaded images (avatars, category/product images) are still saved to
> `uploads/` on local disk, but the REST layer used to serve them back over HTTP at
> `/uploads/*` via Gin's static file route — that route is gone along with REST, and gRPC has
> no equivalent built-in static file serving. This needs a real answer (a small dedicated
> static-file HTTP server, or moving storage to something like S3 with public URLs) before any
> front-end can actually display these images. Tracked in [`plans.md`](./plans.md).

## Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (optional — only needed for
  manual migrations; the app runs pending migrations automatically on startup)
- [`protoc`](https://protobuf.dev/installation/), `protoc-gen-go`, and `protoc-gen-go-grpc`
  (only needed if you're changing a `.proto` file and regenerating `internal/pb/`; running the
  app itself doesn't require them, the generated code is already committed)

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

Create a `.env` file in the project root (see [Configuration](#configuration) for the full
list of required keys) and adjust `config.yml` if needed.

### 4. Run the server

```bash
go run ./cmd
```

Pending migrations run automatically on startup, then the gRPC server starts and blocks on
`0.0.0.0:50051`.

```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 health.HealthService/Check
```

### Regenerating gRPC code after editing a `.proto` file

```bash
protoc --proto_path=proto --go_out=internal/pb --go_opt=paths=source_relative \
       --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
       <module>/<module>.proto
```

Run this from the project root, once per module you changed.

## Configuration

Settings are split between `config.yml` (non-secret, checked into git) and `.env` (secrets,
gitignored). `.env` is loaded first, `config.yml` second; env var names map to config keys by
lowercasing and turning `__` into `.` (e.g. `MYSQL__HOST` → `mysql.host`).

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

**`config.yml`** — non-secret settings (connection pool sizes, token durations, upload
limits); see the file itself for the current values.

### Manual migrations (optional)

The app applies pending migrations automatically on startup, but the CLI is available for
manual control:

```bash
# up
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" up

# down
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" down
```

## API Reference (gRPC)

Server reflection is enabled, so `grpcurl -plaintext localhost:50051 list` (or Postman/SpecHub
with reflection) discovers everything below without needing the `.proto` files locally.
Endpoints marked 🔒 require an `authorization: Bearer <access-token>` gRPC metadata header; 🔒👑
additionally requires the `ADMIN` role (enforced by the authorization interceptor via a
declarative per-method role map, not inside the handler).

### `health.HealthService`

| Method | Description |
|---|---|
| `Check` | Liveness check |

### `auth.AuthService`

| Method | Description |
|---|---|
| `SendOtp` | Request a one-time code for a phone number |
| `CheckOtp` | Verify the code, receive an access token and a refresh token |
| `Me` | 🔒 Get the current session's user id |
| `RefreshToken` | Exchange a valid refresh token for a new access token |
| `Logout` | 🔒 Revoke the current session |

### `user.UserService`

| Method | Description |
|---|---|
| `Profile` | 🔒 Get the current user's profile |
| `UpdateProfile` | 🔒 Update name and/or avatar (avatar sent as raw `bytes`) |
| `ChangePassword` | 🔒 Set a new password |

### `category.CategoryService`

| Method | Description |
|---|---|
| `Create` | 🔒👑 Create a category (root or one level of children) |
| `Update` | 🔒👑 Partially update a category |
| `Delete` | 🔒👑 Soft-delete a category (rejected if a root category still has children) |
| `GetAll` | List all categories as a nested tree |
| `GetOne` | Get a single category by slug, including its children if it's a root category |

### `province.ProvinceService`

| Method | Description |
|---|---|
| `GetAll` | List all provinces |
| `GetOne` | Get a single province by id |

### `address.AddressService`

| Method | Description |
|---|---|
| `Create` | 🔒 Create a new shipping address |
| `GetAll` | 🔒 List the current user's addresses |
| `GetOne` | 🔒 Get a single address by id |
| `Update` | 🔒 Partially update an address |
| `Delete` | 🔒 Soft-delete an address |

### `product.ProductService`

| Method | Description |
|---|---|
| `Create` | 🔒👑 Create a product — **client-streaming**: the client sends one `ProductMetadata` message followed by zero or more `GalleryImage` messages, then the server responds once with the created product |
| `GetAll` | List products with pagination and an optional `category_slug` filter (price filter, search, and sort are not built yet) |
| `GetOneBySlug` | Get a single product by slug, including its gallery images |

`Update` and `Delete` for products are not built yet.

## Authentication Flow

1. **`SendOtp`** (`auth.AuthService/SendOtp`) — client submits a phone number. The server
   finds or creates a user for that number and generates a 5-digit code with `crypto/rand`,
   valid for a short window (`otp_code_duration`).
2. **`CheckOtp`** (`auth.AuthService/CheckOtp`) — client submits the phone number and code. On
   success, the server creates a **session record** in the database and issues a short-lived
   **access token** and a longer-lived **refresh token**, both JWTs carrying the session id.
3. Protected methods go through the shared `grpc.UnaryServerInterceptor` /
   `grpc.StreamServerInterceptor` pair, which parses the access token from the
   `authorization: Bearer <token>` gRPC metadata, loads the referenced session, and rejects
   the call if the session has expired or been revoked.
4. **`RefreshToken`** — when the access token expires, the client exchanges the refresh token
   (sent as an explicit field on `RefreshTokenRequest`) for a new access token, as long as the
   underlying session is still valid.
5. **`Logout`** — revokes the session server-side, so both existing and any future
   access/refresh tokens issued for it stop working immediately. This is the key advantage
   over plain stateless JWTs, where logout can only ever be simulated client-side.

## Roadmap

Tracked in detail in [`plans.md`](./plans.md). The REST → gRPC migration is complete —
health, auth, user, category, province, address, and product all have working gRPC
equivalents, and the REST/Gin layer has been removed entirely. Next up, in rough priority
order:

- [ ] Decide how uploaded images get served now that Gin's static `/uploads` route is gone
- [ ] **Cart** — add/update/remove items, view current cart
- [ ] **Order** — checkout from cart, order history, status transitions
- [ ] **Payment** — payment initiation and gateway callback handling
- [ ] **Product** — listing filters (price, search, sort), update, delete
- [ ] **Microservices migration** — split this monolith into independently deployable
  services, one module at a time, each with its own database; see
  [Migration History](#migration-history) and [`plans.md`](./plans.md)
- [ ] Nice-to-have: reviews, wishlist, coupons

## License

No license has been chosen yet — all rights reserved by default until one is added.
