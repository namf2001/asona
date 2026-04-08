# AGENTS.md

## Project Snapshot
- Monorepo with a Go API (`cmd/`, `internal/`) plus a Next.js frontend (`frontend/`).
- API entrypoint is `cmd/api/main.go`: loads env (`config.Init`), initializes RSA keys (`server.InitRSA`), builds dependencies (`server.NewServer`), then runs Gin with graceful shutdown.
- Runtime wiring is constructor-first in `cmd/api/server/server.go`: `service -> repository.Registry -> controller -> handler -> router`.

## Backend Architecture (What Talks to What)
- HTTP routes are defined centrally in `cmd/api/server/routes.go` and grouped as:
  - public (`/`, `/health`, `/swagger/*any`, `/api/v1/login`)
  - authenticated (`/api/v1/...` with `TokenCheckMiddleware` + `RSAAuthMiddleware`)
  - websocket (`/ws` and `/api/v1/ws`).
- Layering convention (documented in `internal/handler/README.md`):
  - `handler/rest/v1/*` handles transport (bind/parse/write HTTP)
  - `controller/*` handles business orchestration
  - `repository/*` handles SQL.
- Repositories are centralized via `internal/repository/registry.go`; cross-repo transactions use `DoInTx(...)` (nested tx explicitly disallowed).
- API responses are standardized through `internal/handler/response/response.go` and `internal/constants/message_constants.go` (`{ code, message, data }`).

## Security + Request Processing Patterns
- Authenticated routes expect `Authorization: Bearer <jwt>` parsed by `internal/handler/middleware/auth.go`; user info is stored in Gin context as `userID`, `email`.
- Request decryption middleware (`internal/handler/middleware/decode_rsa.go`) is **environment-dependent**:
  - `APP_ENV=local`: RSA decryption is skipped.
  - non-local: request body must be JSON with `data` bytes decryptable by `GlobalRSAKeyPair`.
- RSA keys are loaded/generated at startup in `cmd/api/server/rsa.go` using paths from `config/config.go` (`config/key-pem/*` by default).

## Data + Realtime Integration
- PostgreSQL:
  - connection via `internal/pkg/database/postgres.go` and `internal/repository/db/pg/*`.
  - schema is migration-driven (`migrations/000001_auth.up.sql` .. `000004_chat.up.sql`) and models auth/core/project-management/chat domains.
- Redis:
  - DB 0 for sessions, DB 1 for websocket pub/sub (`internal/constants/redis_constants.go`, `internal/service/redis/new.go`).
- WebSocket:
  - hub/room/client runtime in `internal/service/websocket/{hub.go,room.go,client.go}`.
  - room broadcasts are Redis-backed so message fanout works across instances.
- External integrations are initialized eagerly in server boot: SMTP (`internal/service/mail/service.go`), Google OAuth (`internal/service/oauth/service.go`), S3 (`internal/service/s3/service.go`). Missing env vars can fail startup.

## Developer Workflows (Canonical Commands)
- Primary workflow is Makefile-driven (`Makefile`):
  - `make run` starts Go API then `frontend` install+dev.
  - `make watch` / `make be-dev` uses Air (`.air.toml`) for backend hot reload.
  - `make test`, `make itest`, `make swagger`, `make migrate-up`, `make migrate-down`.
  - `make docker-run` starts app + frontend + postgres + redis via `docker-compose.yml`.
- Env loading behavior (`config/config.go`): app loads `.env`, then optionally `.env.<APP_ENV>` from repo root.

## Frontend Reality Check
- Frontend is currently scaffold-first (`frontend/src/app/page.tsx`), with infra ready (React Query, Zustand, theme provider) but API integration is minimal.
- Existing sample components (`frontend/src/components/user-list.tsx`, `frontend/src/components/add-to-do.tsx`) hit JSONPlaceholder, not this Go API.

## Agent Tips For Safe Changes
- When adding backend features, follow existing vertical slice pattern: add `repository` method -> `controller` use case -> `handler/rest/v1` endpoint -> register route in `cmd/api/server/routes.go`.
- Reuse `response.NewResponse(...)` and constants for error/success payload consistency.
- If touching authenticated endpoints, verify middleware order and Gin context key types (`userID` is assumed `int64` in handlers).
- For schema changes, add new numbered migration files; do not edit old migrations already in use.

