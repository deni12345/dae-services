# 🏗️ Clean Architecture - Project Structure Review & Improvements

## ✅ Current Structure Analysis

### **What's Good:**

```
internal/
├── domain/          ✅ Pure business entities (no dependencies)
├── port/            ✅ Repository interfaces (contracts)
├── infra/firestore/ ✅ Firestore implementations
├── app/             ✅ Use cases & DTOs
└── grpc/            ✅ gRPC handlers
```

### **What's Missing:**

```
internal/
├── grpc/
│   ├── converter/   ❌ MISSING: Proto ↔ DTO converters
│   ├── errors/      ❌ MISSING: Error mapping layer
│   └── interceptor/ ❌ MISSING: Middleware (logging, recovery, idempotency)
│
└── infra/
    ├── redis/       ❌ MISSING: Redis implementation for idempotency
    └── firestore/
        ├── sheet_repo.go (incomplete)          ❌ NEEDS: Full CRUD
        ├── sheet_repo_create.go                ❌ MISSING
        ├── sheet_repo_update.go                ❌ MISSING
        └── sheet_repo_list.go                  ❌ MISSING
```

---

## 🎯 Recommended Structure (Clean Architecture)

```
dae-core/
├── cmd/
│   └── main.go                          # Application entry point
│
├── internal/
│   ├── domain/                          # 🔵 ENTITY LAYER
│   │   ├── common.go                    # Money, Status enums
│   │   ├── order.go                     # Order entity
│   │   ├── user.go                      # User entity
│   │   ├── sheet.go                     # Sheet entity
│   │   └── menu.go                      # MenuItem entity
│   │
│   ├── port/                            # 🟢 PORT LAYER (Interfaces)
│   │   ├── order_repo.go                # OrdersRepo interface
│   │   ├── user_repo.go                 # UsersRepo interface
│   │   ├── sheet_repo.go                # SheetRepo interface
│   │   ├── menu_repo.go                 # MenuRepo interface
│   │   └── idempotency_store.go         # IdempotencyStore interface
│   │
│   ├── infra/                           # 🟡 ADAPTER LAYER (Implementations)
│   │   ├── firestore/
│   │   │   ├── client.go                # Firestore client initialization
│   │   │   │
│   │   │   ├── order_repo.go            # Order repository base
│   │   │   ├── order_repo_create.go     # Create order
│   │   │   ├── order_repo_update.go     # Update order (with callback fn)
│   │   │   ├── order_repo_get.go        # Get order by ID
│   │   │   ├── order_repo_list.go       # List orders
│   │   │   │
│   │   │   ├── user_repo.go             # User repository base
│   │   │   ├── user_repo_update.go      # Update user
│   │   │   ├── user_repo_get.go         # Get user
│   │   │   ├── user_repo_list.go        # List users
│   │   │   │
│   │   │   ├── sheet_repo.go            # Sheet repository base
│   │   │   ├── sheet_repo_create.go     # ❌ TODO
│   │   │   ├── sheet_repo_update.go     # ❌ TODO
│   │   │   ├── sheet_repo_get.go        # ❌ TODO (partially done)
│   │   │   ├── sheet_repo_list.go       # ❌ TODO
│   │   │   │
│   │   │   ├── menu_repo.go             # Menu repository base
│   │   │   └── menu_repo_get.go         # Get menu item
│   │   │
│   │   └── redis/
│   │       ├── client.go                # ❌ TODO: Redis client init
│   │       └── idempotency_store.go     # ✅ CREATED: Idempotency impl
│   │
│   ├── app/                             # 🟣 USE CASE LAYER
│   │   ├── order_dto.go                 # Order DTOs/Commands
│   │   ├── order_usecase.go             # Order use cases
│   │   │
│   │   ├── user_dto.go                  # User DTOs/Commands
│   │   ├── user_usecase.go              # User use cases
│   │   │
│   │   ├── sheet_dto.go                 # ❌ TODO: Sheet DTOs
│   │   ├── sheet_usecase.go             # ❌ TODO: Sheet use cases
│   │   │
│   │   └── pricing_usecase.go           # ⚠️  CONSIDER: Merge into order_usecase
│   │
│   ├── grpc/                            # 🔴 INTERFACE LAYER
│   │   ├── handler/
│   │   │   ├── order_handler.go         # Order gRPC handlers
│   │   │   ├── user_handler.go          # User gRPC handlers
│   │   │   └── sheet_handler.go         # ❌ TODO
│   │   │
│   │   ├── converter/
│   │   │   ├── order_converter.go       # ✅ CREATED: Order proto ↔ DTO
│   │   │   ├── user_converter.go        # ✅ CREATED: User proto ↔ DTO
│   │   │   └── sheet_converter.go       # ❌ TODO
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go                # ✅ CREATED: Error mapping
│   │   │
│   │   └── interceptor/
│   │       ├── logging.go               # ✅ CREATED: Logging interceptor
│   │       ├── recovery.go              # ✅ CREATED: Panic recovery
│   │       ├── idempotency.go           # ❌ TODO: Idempotency interceptor
│   │       ├── auth.go                  # ❌ TODO: Authentication
│   │       └── request_id.go            # ❌ TODO: Request ID injection
│   │
│   └── configs/
│       └── configs.go                   # Configuration loader
│
├── common/
│   └── logx/
│       └── logx.go                      # Structured logging
│
├── proto/                               # Protocol buffer definitions
│   ├── common.proto
│   ├── orders.proto
│   ├── sheets.proto
│   ├── users.proto
│   └── gen/                             # Generated code
│
├── configs.yml                          # Configuration file
├── docker-compose.yml                   # Dev environment
├── Makefile                             # Build commands
└── go.mod
```

---

## 📊 Layer Dependencies (Clean Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│  🔴 Interface Layer (grpc/)                                 │
│  - Handlers, Converters, Interceptors                       │
│  - Depends on: App, Domain                                  │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  🟣 Use Case Layer (app/)                                   │
│  - Business logic, DTOs, Orchestration                      │
│  - Depends on: Port, Domain                                 │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  🟢 Port Layer (port/)                                      │
│  - Repository interfaces                                     │
│  - Depends on: Domain only                                  │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  🔵 Entity Layer (domain/)                                  │
│  - Pure business entities                                    │
│  - No dependencies                                           │
└─────────────────────────────────────────────────────────────┘
                          ↑
┌─────────────────────────────────────────────────────────────┐
│  🟡 Adapter Layer (infra/)                                  │
│  - Firestore, Redis implementations                         │
│  - Depends on: Port, Domain                                 │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Action Items

### **High Priority:**

1. ✅ **Created**: `internal/grpc/converter/` - Proto ↔ DTO converters
2. ✅ **Created**: `internal/grpc/errors/` - Error mapping layer
3. ✅ **Created**: `internal/grpc/interceptor/` - Middleware layer
4. ✅ **Created**: `internal/infra/redis/` - Idempotency store
5. ✅ **Created**: `internal/infra/firestore/sheet_repo.go` - Base implementation

### **Medium Priority:**

6. ❌ **TODO**: Complete `sheet_repo` CRUD operations
7. ❌ **TODO**: Create `sheet_usecase.go` + `sheet_dto.go`
8. ❌ **TODO**: Implement idempotency interceptor
9. ❌ **TODO**: Move existing handlers to use new converter layer

### **Low Priority:**

10. ❌ **TODO**: Add authentication interceptor
11. ❌ **TODO**: Add request ID interceptor
12. ❌ **TODO**: Consider merging `pricing_usecase` into `order_usecase`

---

## 🔧 Refactoring Guide

### **Step 1: Move converters out of handlers**

**Before:**

```go
// internal/grpc/user_handler.go
func (h *UserHandler) UpdateUser(ctx context.Context, req *corev1.UpdateUserReq) (*corev1.UpdateUserResp, error) {
    patch := &app.UpdateUserReq{
        ID:         req.Id,
        UserName:   req.DisplayName,
        // ... mapping logic here
    }
    // ...
}
```

**After:**

```go
// internal/grpc/converter/user_converter.go
func UpdateUserReqFromProto(req *corev1.UpdateUserReq) *app.UpdateUserReq {
    return &app.UpdateUserReq{
        ID:         req.Id,
        UserName:   req.DisplayName,
        // ...
    }
}

// internal/grpc/handler/user_handler.go
func (h *UserHandler) UpdateUser(ctx context.Context, req *corev1.UpdateUserReq) (*corev1.UpdateUserResp, error) {
    dto := converter.UpdateUserReqFromProto(req)
    user, err := h.uc.UpdateUser(ctx, dto)
    if err != nil {
        return nil, grpcerrors.ToGRPCStatus(err)
    }
    return &corev1.UpdateUserResp{
        User: converter.UserToProto(user),
    }, nil
}
```

### **Step 2: Use error mapping layer**

**Before:**

```go
if err != nil {
    if errors.Is(err, app.ErrUserNotFound) {
        return nil, status.Error(codes.NotFound, err.Error())
    }
    return nil, status.Error(codes.Internal, err.Error())
}
```

**After:**

```go
if err != nil {
    return nil, grpcerrors.ToGRPCStatus(err)
}
```

### **Step 3: Register interceptors in main.go**

```go
// cmd/main.go
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        interceptor.RecoveryInterceptor(),
        interceptor.LoggingInterceptor(logger),
        // TODO: interceptor.IdempotencyInterceptor(idemStore),
        // TODO: interceptor.AuthInterceptor(),
    ),
)
```

---

## 🎓 Clean Architecture Principles Applied

1. ✅ **Dependency Rule**: Dependencies point inward (grpc → app → port → domain)
2. ✅ **Interface Segregation**: Small, focused interfaces in `port/`
3. ✅ **Single Responsibility**: Each layer has one reason to change
4. ✅ **Dependency Inversion**: High-level modules don't depend on low-level (use interfaces)
5. ✅ **Separation of Concerns**: Business logic separate from infrastructure

---

## 📝 Summary

**Current Status**: 70% Clean Architecture compliant
**Missing**: Converter layer, Error handling layer, Interceptors, Redis impl, Complete Sheet repo

**Next Steps**:

1. Move existing handlers to use new converter layer
2. Implement remaining interceptors
3. Complete Sheet repository CRUD
4. Create Sheet use case layer

This structure follows **Robert C. Martin's Clean Architecture** and **Hexagonal Architecture (Ports & Adapters)** best practices! 🎯
