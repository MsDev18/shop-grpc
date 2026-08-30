# Shop API

A backend for an online shop, built in Go with a clean, layered architecture. This is a
learning project focused on writing idiomatic, well-structured backend Go — not a finished
storefront yet, but the foundation (auth, users, categories, addresses) is production-quality.
It is being migrated from REST-only to **gRPC**, module by module, using the strangler-fig
pattern: both protocols run side by side against the same service layer until the migration
is complete.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-008ECF)
![gRPC](https://img.shields.io/badge/gRPC-migrating-4285F4?logo=grpc&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.4-4479A1?logo=mysql&logoColor=white)
![Status](https://img.shields.io/badge/status-in%20development-yellow)

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [gRPC Migration](#grpc-migration)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Reference (REST)](#api-reference-rest)
- [API Reference (gRPC)](#api-reference-grpc)
- [Authentication Flow](#authentication-flow)
- [Roadmap](#roadmap)
- [License](#license)

## Features

- **Clean architecture** — strict separation between `handler` → `validator` → `service` →
  `repository` layers, with each layer depending only on interfaces of the layer below it.
  Both the REST handlers and the gRPC servers sit on top of the exact same `validator` /
  `service` / `repository` layers — only the transport-facing layer differs per protocol.
- **OTP + JWT authentication** — phone-number login via one-time codes (generated with
  `crypto/rand`, not `math/rand`), backed by short-lived access tokens and long-lived refresh
  tokens.
- **Real session revocation** — refresh tokens are tied to a session record in the database,
  so logout actually invalidates a session instead of just discarding a token client-side.
  This is enforced identically for REST (Gin middleware) and gRPC (a unary interceptor).
- **Role-based access control** — only admins can create/update/delete categories, for
  example — enforced via REST middleware and, on the gRPC side, an equivalent
  `auth.RequireRole` check inside each handler.
- **Storage-agnostic image pipeline** — uploaded images are validated by content sniffing
  (not just file extension), auto-oriented, resized, and re-encoded before being handed to a
  pluggable `Storage` interface (local disk today, swappable for S3 or similar later without
  touching business logic). Image bytes travel as `multipart.FileHeader` on the REST side and
  as raw `bytes` fields on the gRPC side, converging on the same in-memory representation
  before reaching the service layer.
- **Structured error handling** — a custom `richerror` package carries operation, message,
  kind, and metadata through every layer. A `mapper` package turns a `richerror.Kind` into
  either an HTTP status code (REST) or a `google.golang.org/grpc/codes.Code` (gRPC), so both
  protocols report the same logical error consistently.
- **Configuration via `.env` + YAML** — secrets and per-environment values load from `.env`,
  general settings from `config.yml`, merged with [koanf](https://github.com/knadh/koanf).
- **Migrations built into the binary** — schema migrations run automatically on startup via a
  small wrapper around [golang-migrate](https://github.com/golang-migrate/migrate), with the
  standalone CLI still available for manual up/down.

## Tech Stack

| Layer | Choice |
|---|---|
| Language | Go 1.26 |
| REST framework | [Gin](https://github.com/gin-gonic/gin) |
| RPC framework | [gRPC-Go](https://github.com/grpc/grpc-go) + [Protocol Buffers](https://protobuf.dev/) (`protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`) |
| Database | MySQL 8.4 (`database/sql` + `go-sql-driver/mysql`, no ORM) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Config | [koanf](https://github.com/knadh/koanf) (`.env` + YAML) |
| Auth | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) + bcrypt (`golang.org/x/crypto`) |
| Validation | [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) |
| Image processing | [disintegration/imaging](https://github.com/disintegration/imaging) |
| Local dev infra | Docker Compose (MySQL + phpMyAdmin) |

## Architecture

Every request flows through the same pipeline, regardless of module — the only thing that
changes between REST and gRPC is the top two layers:

```
 REST request                          gRPC request
     │                                      │
     ▼                                      ▼
 Router (gin route groups)          Unary interceptor (auth.Interceptor)
     │                                      │
     ▼                                      ▼
 Middleware (auth / role checks)     gRPC server method (per module)
     │                                      │
     └──────────────────┬───────────────────┘
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
that into the right shape for whichever protocol is asking:

- REST: an HTTP status code, serialized by the shared `response` package into one consistent
  JSON envelope:

  ```json
  {
    "code": 200,
    "message": "human readable message",
    "data": { },
    "errors": { }
  }
  ```

- gRPC: a `google.golang.org/grpc/codes.Code` wrapped in a `status.Error`, so a gRPC client
  sees a proper gRPC status (e.g. `NotFound`, `PermissionDenied`, `Unauthenticated`) instead of
  a generic failure.

## gRPC Migration

This project is being converted from REST-only to gRPC using the **strangler fig pattern**:
instead of a big-bang rewrite, a gRPC server is started in a goroutine alongside the existing
Gin server, and each module is migrated one at a time. Both servers stay live throughout the
migration, share the same MySQL database, and run inside the same process — there is no
microservices split, that is an explicit separate future decision, not something this
migration does implicitly.

- **REST** stays on `:3000` (unchanged, module by module, until it's fully replaced).
- **gRPC** runs on `:50051`, registered with [server reflection](https://github.com/grpc/grpc-go/tree/master/reflection) enabled so tools like `grpcurl` or Postman/SpecHub can discover services without needing the `.proto` files locally.
- Authentication/authorization is enforced once, globally, via a single `grpc.UnaryServerInterceptor` (`internal/api/grpc/auth`). Methods that don't require a login are explicitly listed in a `publicMethods` allow-list, keyed by full method name (`/<package>.<Service>/<Method>`); everything else requires a valid `authorization: Bearer <token>` gRPC metadata header. Role-restricted methods (admin-only, for example) additionally call `auth.RequireRole(ctx, entity.AdminRole)` at the top of the handler.
- `.proto` files live under `proto/<module>/`, one file per module, and generate into `internal/pb/<module>/` via `protoc` — see [Getting Started](#getting-started) for the exact command. Generated files are committed but never hand-edited.
- The last step of the migration is deleting the Gin/REST layer entirely, once every module (including ones not yet built in REST, like cart/order/payment) has a gRPC equivalent.

Current per-module status is tracked in [`plans.md`](./plans.md); the short version is in
[API Reference (gRPC)](#api-reference-grpc) below.

## Project Structure

```
shop/
├── cmd/                      # composition root — wires repositories, services, REST handlers
│                              # and gRPC servers per module
├── proto/                    # .proto source files, one folder per module (source of truth)
├── internal/
│   ├── api/
│   │   ├── grpc/              # gRPC servers — one package per module, plus the shared auth
│   │   │                      # interceptor (internal/api/grpc/auth) and the top-level Server
│   │   │                      # that registers every service and starts listening on :50051
│   │   ├── handler/          # REST (Gin) — one package per module
│   │   ├── middleware/       # REST auth & role-based middleware
│   │   ├── router/           # REST route registration
│   │   └── server/           # REST HTTP server setup
│   ├── config/                # .env + YAML config loading (koanf)
│   ├── dto/                   # request/response shapes per module (shared by REST & gRPC)
│   ├── entity/                # domain models
│   ├── migrator/              # golang-migrate wrapper, runs on startup
│   ├── pb/                    # generated Go code from proto/ — never edited by hand
│   ├── pkg/
│   │   ├── claims/            # JWT access/refresh token creation & parsing
│   │   ├── imageprocessor/    # upload validation, resize/encode, storage-agnostic save
│   │   ├── mapper/             # richerror.Kind -> HTTP status code, and -> grpc/codes.Code
│   │   ├── response/           # consistent JSON response envelope (REST only)
│   │   └── richerror/          # structured application error type
│   ├── repository/mysql/      # one package per module, raw SQL
│   ├── service/               # one package per module, business logic (shared by REST & gRPC)
│   └── validator/             # one package per module, input validation (shared by REST & gRPC)
├── migrations/                 # golang-migrate SQL files
├── uploads/                     # locally stored uploaded images (served at /uploads over REST)
├── config.yml
├── docker-compose.yml
└── plans.md                     # module/endpoint checklist, REST + gRPC migration status
```

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

Pending migrations run automatically on startup. Two servers come up from this one command:

- REST at `http://localhost:3000` (or whatever `server.host`/`server.port` you configured)
- gRPC at `0.0.0.0:50051`

```bash
curl http://localhost:3000/health-check
grpcurl -plaintext localhost:50051 list
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

**`config.yml`** — non-secret settings (connection pool sizes, server timeouts, token
durations, upload limits); see the file itself for the current values.

### Manual migrations (optional)

The app applies pending migrations automatically on startup, but the CLI is available for
manual control:

```bash
# up
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" up

# down
migrate -source file://migrations -database "mysql://shop:shop-pass@(localhost:3308)/shop?parseTime=true&x-migrations-table=migrations" down
```

## API Reference (REST)

All responses use the envelope shown in [Architecture](#architecture). Endpoints marked 🔒
require a valid `Authorization: Bearer <access-token>` header; 🔒👑 additionally requires the
`ADMIN` role.

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

## API Reference (gRPC)

Server reflection is enabled, so `grpcurl -plaintext localhost:50051 list` (or Postman/SpecHub
with reflection) discovers everything below without needing the `.proto` files locally.
Endpoints marked 🔒 require an `authorization: Bearer <access-token>` gRPC metadata header; 🔒👑
additionally requires the `ADMIN` role (checked inside the handler via `auth.RequireRole`, not
by the global interceptor).

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

Not started yet — will be the first module to use gRPC client-streaming, for gallery image
uploads.

## Authentication Flow

1. **`send-otp`** (REST: `POST /auth/send-otp`, gRPC: `auth.AuthService/SendOtp`) — client
   submits a phone number. The server finds or creates a user for that number and generates a
   5-digit code with `crypto/rand`, valid for a short window (`otp_code_duration`).
2. **`check-otp`** (REST: `POST /auth/check-otp`, gRPC: `auth.AuthService/CheckOtp`) — client
   submits the phone number and code. On success, the server creates a **session record** in
   the database and issues a short-lived **access token** and a longer-lived **refresh
   token**, both JWTs carrying the session id.
3. Protected routes/methods go through the auth check for their protocol — Gin's
   `AuthRequired` middleware for REST, or the shared `grpc.UnaryServerInterceptor` for gRPC —
   which parses the access token, loads the referenced session, and rejects the request if the
   session has expired or been revoked. Both read the token the same way: `Authorization:
   Bearer <token>` as an HTTP header for REST, `authorization: Bearer <token>` as gRPC metadata.
4. **`refresh-token`** — when the access token expires, the client exchanges the refresh token
   for a new access token, as long as the underlying session is still valid. (REST reads it
   from an httpOnly cookie; gRPC has no cookie mechanism, so it's an explicit field on
   `RefreshTokenRequest` instead.)
5. **`logout`** — revokes the session server-side, so both existing and any future
   access/refresh tokens issued for it stop working immediately. This is the key advantage
   over plain stateless JWTs, where logout can only ever be simulated client-side — and it's
   enforced identically for both protocols, since both check session revocation in the
   database rather than trusting the JWT signature alone.

## Roadmap

Tracked in detail in [`plans.md`](./plans.md). Foundational modules (config, error handling,
migrations, HTTP server) and the auth/user/category/province/address modules are done on both
REST and gRPC. Still to build/migrate, in rough priority order:

- [ ] **Product** — REST create + get-by-slug are done; REST listing/update/delete still
  pending, and the gRPC equivalent hasn't started yet (this will introduce client-streaming
  for gallery image uploads)
- [ ] **Cart** — add/update/remove items, view current cart (REST first, gRPC to follow)
- [ ] **Order** — checkout from cart, order history, status transitions
- [ ] **Payment** — payment initiation and gateway callback handling
- [ ] **Remove the REST/Gin layer** once every module above has a working gRPC equivalent
- [ ] Nice-to-have: reviews, wishlist, coupons

## License

No license has been chosen yet — all rights reserved by default until one is added.
