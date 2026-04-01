# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

### Backend (dropship-api)
```bash
cd backend && go build ./...                    # Build
cd backend && go run main.go -f etc/config.yaml # Run (requires MySQL + etcd + user-center-rpc)
```

### User Center
```bash
cd user-center && go build ./...                          # Build
cd user-center && go run cmd/rpc/main.go -f etc/rpc.yaml  # Run RPC (start first, registers to etcd)
cd user-center && go run cmd/api/main.go -f etc/api.yaml  # Run API
```

### Frontend
```bash
cd frontend && npm run dev       # Dev server
cd frontend && npx next build    # Production build (also runs TypeScript check)
```

### Startup Order
1. `user-center/cmd/rpc` (registers to etcd, seeds roles + super admin)
2. `user-center/cmd/api`
3. `backend` (discovers user-center-rpc via etcd)
4. `frontend`

### Protobuf Regeneration
```bash
cd user-center
protoc --go_out=proto/pb --go_opt=paths=source_relative \
       --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
       proto/usercenter.proto
```

## Architecture

Three services in one repo, two Go modules (`backend/go.mod`, `user-center/go.mod`). Backend depends on user-center via `go mod replace` directive for the proto/pb package.

### DDD Four-Layer Architecture (both Go services)

```
server/     → HTTP handlers, routes
application/→ Use case orchestration, calls domain + repo
domain/     → Aggregate roots (entities with behavior), domain services, repository interfaces, gateway interfaces
model/      → PO (GORM structs), DTO (request/response), anticorruption (external API types)
repository/ → GORM implementations of domain repository interfaces
gateway/    → Mock implementations of domain gateway interfaces
```

**Dependency direction**: `server → application → domain → model`. Never reverse. Repository/gateway interfaces are defined in `domain/`, implementations are in `repository/` and `gateway/`.

### Cross-Service Communication

`dropship-api` authenticates every request by calling `user-center-rpc.VerifyToken` via gRPC. The RPC client is injected through `internal/wire.go` (manual DI, no codegen). Service discovery uses etcd with key `user-center.rpc`.

### Key Patterns

- **Aggregates**: SourceItem, Product, PublishTask (backend); User, Role (user-center). State transitions are methods on the aggregate root, not external logic.
- **Anti-corruption layer**: `domain/platform/gateway.go` defines interfaces for external platforms (1688, PDD). Currently Mock implementations in `gateway/`. When adding real platform integrations, implement the interface without changing domain code.
- **CQRS-lite**: Write operations go through aggregate roots. Read-heavy queries (list with filters) use `queries/` package that hits DB directly, bypassing the domain layer.

### Frontend

Next.js 14 App Router with shadcn/ui (base-ui variant, NOT Radix). `src/lib/api.ts` has dual base URLs: dropship-api (:8888) for business APIs, user-center-api (:8880) for auth/user APIs. Auth guard in `components/layout/auth-guard.tsx` redirects unauthenticated users to `/login`.

**Important**: This project uses Next.js 16+ with shadcn/ui based on `@base-ui/react`, not `@radix-ui`. The `asChild` prop does not exist — use `render` prop instead. Check `node_modules/next/dist/docs/` for current API docs before assuming Next.js patterns.

## Conventions

- Go naming: snake_case for files, CamelCase for types/functions
- All Go code formatted with `go fmt` and `goimports`
- Config: `etc/*.yaml` files, loaded via `go-zero/core/conf.MustLoad`. etcd config overlays file config.
- Google OAuth in user-center-api routes through `http://127.0.0.1:7890` proxy for local development
- MySQL root password: `root123`, database: `dropship`, all tables auto-migrated via GORM
- Super admin auto-created on first RPC startup: `admin` / `admin123`
- JWT secret shared between `user-center/etc/api.yaml` and `user-center/etc/rpc.yaml` — must match

## Sensitive Files (do not commit)

- `user-center/etc/api.yaml` — contains Google OAuth Client ID and Secret
