# 用户中心微服务设计文档

## 概述

新增用户中心微服务，提供 Google OAuth 2.0 第三方登录、JWT 鉴权、RBAC 角色权限体系。与现有 dropship-api 通过 gRPC (go-zero zRPC) 通信，使用 etcd 做服务发现。前端新增登录页和用户管理模块，支持初始化超级管理员。

**技术栈**：Go-zero (REST + zRPC), GORM, JWT, Google OAuth 2.0, Protobuf  
**架构风格**：DDD + Clean Architecture，与现有 backend 服务保持一致

---

## 微服务拓扑

```
Frontend (Next.js :3000)
  │ HTTP                         │ HTTP
  ▼                              ▼
user-center-api (:8880)    dropship-api (:8888)
  │                              │
  │ gRPC                         │ gRPC
  ▼                              ▼
user-center-rpc (:8881, 注册到 etcd)
  │
  ▼
MySQL (dropship DB, 共用)
```

- **user-center-api**：REST 服务，处理 Google OAuth 回调、用户管理、超管初始化
- **user-center-rpc**：gRPC 服务，提供 VerifyToken / CheckRole / GetUserInfo 给其他微服务调用
- **dropship-api**：现有服务，新增 JWT 鉴权中间件，通过 RPC 调用 user-center-rpc 验证 token 和角色

---

## Google OAuth 2.0 登录流程

1. 前端点击「Sign in with Google」→ 跳转 `user-center-api` 的 `/api/auth/google` 接口
2. 后端构造 Google OAuth URL，302 重定向到 Google 授权页
3. 用户授权后，Google 回调 `/api/auth/google/callback?code=xxx`
4. 后端用 code 换取 Google access token，获取用户信息（email, name, avatar）
5. 查找或创建本地用户（以 google_id 为唯一标识）
6. 生成 JWT token（包含 userId, email, roles），返回给前端
7. 前端存储 JWT，后续请求在 Header 中携带 `Authorization: Bearer <JWT>`

---

## RBAC 权限模型

### 数据库表

#### users
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| google_id | VARCHAR(128) UNIQUE | Google 用户ID |
| email | VARCHAR(256) UNIQUE | 邮箱 |
| name | VARCHAR(128) | 用户名 |
| avatar | VARCHAR(512) | 头像URL |
| status | VARCHAR(32) | 状态: active/disabled |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

#### roles
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| name | VARCHAR(64) UNIQUE | 角色名 |
| description | VARCHAR(256) | 描述 |
| is_default | BOOL | 是否默认角色（新用户自动分配） |
| created_at | DATETIME | 创建时间 |

#### user_roles
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| user_id | BIGINT FK | 用户ID |
| role_id | BIGINT FK | 角色ID |
| created_at | DATETIME | 创建时间 |

### 预设角色

| 角色 | 说明 | 用户管理 | 角色分配 | 货源/商品/发品 |
|------|------|---------|---------|--------------|
| super_admin | 超级管理员（不可删除） | ✅ | ✅ | ✅ |
| admin | 管理员 | ✅ | ❌ | ✅ |
| operator | 运营人员 | ❌ | ❌ | ✅ |
| viewer | 只读用户（默认角色） | ❌ | ❌ | 只读 |

角色权限映射在代码中硬编码，不额外建权限点表。

### 超管初始化

- 接口：`POST /api/init/super-admin`，body: `{email: "xxx@gmail.com"}`
- 条件：仅当 users 表中没有 super_admin 角色的用户时可调用
- 效果：将指定邮箱的已登录用户提升为 super_admin
- 前端在首次访问时检测是否需要初始化，展示初始化引导页

---

## RPC 接口定义 (Protobuf)

```protobuf
syntax = "proto3";

package usercenter;
option go_package = "./pb";

message VerifyTokenReq {
  string token = 1;
}

message VerifyTokenResp {
  int64 user_id = 1;
  string email = 2;
  string name = 3;
  repeated string roles = 4;
}

message CheckRoleReq {
  int64 user_id = 1;
  string role = 2;
}

message CheckRoleResp {
  bool has_role = 1;
}

message GetUserInfoReq {
  int64 user_id = 1;
}

message UserInfo {
  int64 id = 1;
  string email = 2;
  string name = 3;
  string avatar = 4;
  string status = 5;
  repeated string roles = 6;
}

service UserCenter {
  rpc VerifyToken(VerifyTokenReq) returns (VerifyTokenResp);
  rpc CheckRole(CheckRoleReq) returns (CheckRoleResp);
  rpc GetUserInfo(GetUserInfoReq) returns (UserInfo);
}
```

---

## 用户中心目录结构

```
user-center/
├── api/                              # REST API 服务
│   ├── main.go
│   ├── etc/config.yaml
│   └── internal/
│       ├── server/
│       │   ├── routes.go
│       │   ├── auth_handler.go       # Google OAuth 登录/回调
│       │   ├── user_handler.go       # 用户列表/详情/禁用
│       │   └── init_handler.go       # 超管初始化
│       ├── application/
│       │   ├── auth_app.go           # OAuth + JWT 逻辑
│       │   └── user_app.go           # 用户管理逻辑
│       ├── middleware/
│       │   └── auth_middleware.go     # JWT 验证中间件
│       └── wire.go
├── rpc/                              # gRPC 服务
│   ├── main.go
│   ├── etc/config.yaml
│   ├── pb/
│   │   ├── usercenter.proto
│   │   └── (generated files)
│   └── internal/
│       ├── server/
│       │   └── usercenter_server.go
│       └── wire.go
├── internal/                         # API 和 RPC 共享层
│   ├── domain/
│   │   ├── user/
│   │   │   ├── entity/user.go
│   │   │   └── repository/repository.go
│   │   └── role/
│   │       ├── entity/role.go
│   │       └── repository/repository.go
│   ├── model/
│   │   ├── po/mysql/user.go
│   │   ├── po/mysql/role.go
│   │   ├── po/mysql/user_role.go
│   │   └── dto/user_dto.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   └── role_repo.go
│   └── pkg/
│       └── jwt/jwt.go
└── go.mod
```

依赖方向与 backend 保持一致：`server → application → domain → model`

---

## 对现有服务的改造

### dropship-api (backend/)

1. **新增 RPC client 配置**：`etc/config.yaml` 加入 UserCenter RPC 地址
2. **新增鉴权中间件**：`internal/middleware/auth_middleware.go`，拦截所有 `/api/*` 请求，通过 RPC 调用 `VerifyToken` 验证 JWT
3. **路由改造**：需鉴权的路由加中间件，`/api/auth/*` 相关路由不鉴权

### 前端 (frontend/)

1. **新增登录页** `/login`：Google Sign-in 按钮
2. **新增用户管理页** `/users`：用户列表、角色分配（仅 super_admin 可见）
3. **新增超管初始化页** `/init`：首次使用引导
4. **全局鉴权**：未登录跳转 `/login`，JWT 存 localStorage，请求自动带 Authorization header
5. **侧边栏权限控制**：根据角色隐藏/显示菜单项

---

## API 端点汇总

### user-center-api (:8880)

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | /api/auth/google | 发起 Google OAuth | 无 |
| GET | /api/auth/google/callback | Google 回调 | 无 |
| GET | /api/auth/me | 获取当前用户信息 | JWT |
| POST | /api/init/super-admin | 初始化超管 | JWT |
| GET | /api/users | 用户列表 | JWT + admin |
| PUT | /api/users/status | 启用/禁用用户 | JWT + admin |
| PUT | /api/users/role | 分配角色 | JWT + super_admin |
| GET | /api/init/check | 检查是否需要初始化 | 无 |

---

## 验证方案

1. **RPC 服务**：启动 user-center-rpc，用 grpcurl 测试 VerifyToken / CheckRole
2. **OAuth 流程**：启动 user-center-api，浏览器访问 `/api/auth/google`，完成 Google 登录，拿到 JWT
3. **鉴权集成**：dropship-api 加中间件后，无 token 请求返回 401，有效 token 正常通过
4. **超管初始化**：首次访问前端 → 自动跳转初始化页 → 指定邮箱为超管
5. **权限控制**：不同角色用户登录，验证菜单可见性和 API 访问权限
