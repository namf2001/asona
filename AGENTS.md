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
  - `APP_ENV=dev`: RSA decryption is skipped.
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

# Asona Backend Coding Conventions & Architecture

Tài liệu này quy định cấu trúc và cách triển khai các Layer trong dự án Asona, tuân thủ theo pattern của dự án **beta-workplace**.

---

## 0. Quy tắc Documentation (Comments)

- **Mọi hàm/struct/interface public** (viết hoa chữ cái đầu) đều PHẢI có comment mô tả mục đích và chức năng.
- Comment phải bắt đầu bằng tên của thực thể đó.
- Ví dụ: `// Create inserts a new record into the database.`
- Quy tắc này áp dụng cho TẤT CẢ các layer (Handler, Controller, Repository).

---

## 0.1. Quy tắc Import

Imports phải được chia thành 3 nhóm rõ rệt, ngăn cách bởi một dòng trống:
1. **Standard Library**: Các package có sẵn của Go (ví dụ: `context`, `time`, `fmt`).
2. **Third-party Libraries**: Các thư viện bên ngoài (ví dụ: `github.com/gin-gonic/gin`, `github.com/google/uuid`).
3. **Internal Packages**: Các package nội bộ trong dự án (ví dụ: `asona/internal/...`).

Ví dụ:
```go
import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"asona/internal/repository"
)
```

---

## 1. Repository Layer (Tầng dữ liệu)

Nằm tại `internal/repository/[domain]`.

### Quy tắc chung:
- **Tách file**: Mỗi phương thức CRUD/Logic phải nằm trong một file riêng biệt (ví dụ: `create.go`, `get_by_id.go`).
- **No-Pointer Rule**: Interface và Struct implementation (`impl`) sử dụng **value receiver** để tránh lỗi nil pointer dereference.
- **pg.ContextExecutor**: Sử dụng interface này để thực thi SQL, cho phép hỗ trợ cả `sql.DB` và `sql.Tx`.

### Cấu trúc file `new.go`:
```go
package domain

import (
	"asona/internal/repository/db/pg"
)

type Repository interface {
	Create(ctx context.Context, data model.Data) (model.Data, error)
}

type impl struct {
	db pg.ContextExecutor
}

func New(db pg.ContextExecutor) Repository {
	return impl{db: db}
}
```

### Triển khai phương thức (ví dụ `create.go`):
```go
func (i impl) Create(ctx context.Context, data model.Data) (model.Data, error) {
	// Thực thi SQL qua i.db
	return data, nil
}
```

---

## 2. Controller Layer (Tầng nghiệp vụ - Business Logic)

Nằm tại `internal/controller/[domain]`.

### Quy tắc chung:
- **Độc lập**: Controller không làm việc trực tiếp với database, chỉ gọi qua Repo Registry.
- **Mapping**: Chuyển đổi dữ liệu từ `model` (DB) sang `Response` struct (Dùng cho API) tại đây.
- **Transaction**: Sử dụng `repo.DoInTx` nếu cần thực hiện nhiều hành động trong một transaction.
- **Local Structs**: Định nghĩa các struct `Input` hoặc `Output` phục vụ cho logic ngay trong file controller (thường đặt ở cuối file hoặc ngay trên hàm sử dụng). **KHÔNG** tạo thêm các file riêng lẻ chỉ để chứa struct.

### Cấu trúc file `new.go`:
```go
type Controller interface {
	DoSomething(ctx context.Context, input Input) (Response, error)
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}
```

---

## 3. Handler Layer (Tầng Giao tiếp - REST API)

Nằm tại `internal/handler/rest/v1/[domain]`.

### Quy tắc chung:
- **Unified Response**: Luôn sử dụng package `response.NewResponse` để trả về kết quả.
- **Constants**: Không bao giờ hardcode chuỗi tin nhắn hoặc mã lỗi. Sử dụng `constants.ResponseCode` từ `message_constants.go`.
- **Local Structs**: Định nghĩa các struct `request` và `response` ngay trong file handler. **KHÔNG** dùng chung (reuse) struct từ tầng Controller để Bind JSON hoặc làm Output. Điều này giúp tách biệt hoàn toàn layer giao tiếp và layer nghiệp vụ.
- **Logging**: Luôn thực hiện ghi log thông tin (Error log cho các lỗi, Info log cho các phản hồi thành công) cho mỗi kết quả trả về cho người dùng.
    - Pattern: `logger.ERROR.Printf("[TênHàm] message: %+v", err)`
    - Ví dụ: `logger.ERROR.Printf("[AuthRegister] failed request param: %+v", err)`
- **Abort**: Trong middleware, sử dụng `c.AbortWithStatusJSON`.

### Ví dụ Handler:
```go
type createRequest struct {
	Name string `json:"name" binding:"required"`
}

type dataResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// handle error...
	}
	// 1. Map req -> Controller Input
	// 2. Call Controller
	// 3. Map Controller Output -> dataResponse
	// 4. Trả về response chuẩn:
	c.JSON(http.StatusOK, response.NewResponse(
		constants.Success.Code,
		constants.Success.Message,
		dataResponse{...},
	))
}
```

---

## 4. Repository Registry (Khởi tạo & Transaction)

Nằm tại `internal/repository/registry.go`.

- **Mục đích**: Là đầu mối duy nhất để truy cập tất cả các repository con.
- **DoInTx**: Một phương thức đặc biệt để bọc các hành động vào transaction với cơ chế **Backoff Retry** tự động.

### Cách sử dụng trong Controller:
```go
err := i.repo.DoInTx(ctx, func(ctx context.Context, txRepo repository.Registry) error {
	// Sử dụng txRepo để đảm bảo các câu lệnh chạy trong cùng 1 transaction
	created, _ := txRepo.User().Create(...)
	_ = txRepo.Account().Create(created.ID, ...)
	return nil
}, nil)
```

---

## 5. Message Constants (Quản lý lỗi tập trung)

Nằm tại `internal/constants/message_constants.go`.

Mọi mã lỗi (Code) và tin nhắn (Message) phải được định nghĩa tại đây. Ví dụ:
- `Success = ResponseCode{"00", "Success"}`
- `InvalidToken = ResponseCode{"ERR009", "Invalid token"}`

---

## 6. Wiring (Khởi tạo hệ thống)

Thực hiện tại `cmd/api/server/server.go`:
1. Khởi tạo `database.Service`.
2. Truyền `db.DB()` vào `repository.New`.
3. Truyền `repo` vào các `Controller.New`.
4. Truyền `controller` vào các `Handler.New`.
5. Đăng ký route trong `router.go`---

## 7. Error Handling (Domain-Specific Errors)

Để tránh lỗi vòng lặp (Circular Dependency) và đảm bảo API trả về thông báo lỗi sạch:
- **Repository Errors**: Định nghĩa lỗi tại `internal/repository/[domain]/errors.go`. Ví dụ: `ErrUserNotFound`.
- **Controller Errors**: Định nghĩa lỗi tại `internal/controller/[domain]/errors.go`. Ví dụ: `ErrUserNotFound`.
- **Mapping**: Tầng Controller sử dụng `errors.Is(err, repo.Err...)` để chuyển đổi lỗi từ Repository sang lỗi Controller tương ứng trước khi trả về cho Handler. Handler sẽ trả về `err.Error()` cho Client.

## 8. Stack Trace Usage (Error Tracing)

Để hỗ trợ gỡ lỗi và truy vết vị trí lỗi phát sinh:
- **Usage**: Sử dụng `github.com/pkg/errors` (bí danh `pkgerrors`) để bọc lỗi bằng `pkgerrors.WithStack(err)`.
- **Repository**: Tất cả lỗi trả về từ tầng Repository (lỗi SQL hoặc lỗi Domain) PHẢI được bọc bằng `WithStack`.
- **Controller**: Trả về lỗi trực tiếp nếu nó đã được bọc từ Repository. Nếu Controller tạo lỗi mới, hãy bọc nó.
- **Log**: Hệ thống log sẽ tự động in stack trace nếu lỗi được bọc đúng cách, giúp xác định chính xác tệp và dòng code gây lỗi.

Ví dụ trong Repository:
```go
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return model.Data{}, pkgerrors.WithStack(ErrDataNotFound)
    }
    return model.Data{}, pkgerrors.WithStack(fmt.Errorf("db error: %w", err))
}
```
