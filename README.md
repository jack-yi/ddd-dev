# DDD-Dev: 电商一键代发工具

基于 DDD + Clean Architecture 的电商运营工具，实现从批发平台（1688）选品，到卖货平台（拼多多）发品的一键代发能力。

## 技术架构

```
                        ┌─────────────────────────────────┐
                        │     Frontend (Next.js :3000)     │
                        │  登录 / 货源 / 商品 / 发品 / 用户  │
                        └──────┬──────────────┬────────────┘
                               │ HTTP          │ HTTP
                               ▼               ▼
                   ┌───────────────────┐ ┌──────────────────┐
                   │ user-center-api   │ │   dropship-api   │
                   │   (go-zero :8880) │ │  (go-zero :8888) │
                   │                   │ │                  │
                   │ Google OAuth 2.0  │ │ 货源/商品/发品    │
                   │ 用户管理/角色分配  │ │ API              │
                   │ JWT 签发          │ │                  │
                   └────────┬──────────┘ └────────┬─────────┘
                            │ gRPC                │ gRPC
                            ▼                     ▼
                   ┌──────────────────────────────────────┐
                   │      user-center-rpc (zRPC :8881)     │
                   │   VerifyToken / CheckRole / GetUser    │
                   │      注册到 etcd 做服务发现             │
                   └──────────────────────────────────────┘
                            │                │
                            ▼                ▼
                   ┌──────────────┐  ┌──────────────┐
                   │  MySQL 8.0   │  │    etcd      │
                   │  (dropship)  │  │  配置中心     │
                   └──────────────┘  │  服务发现     │
                                     └──────────────┘
```

## 技术栈

### 后端

| 组件 | 技术选型 | 说明 |
|------|---------|------|
| 框架 | Go-zero | REST API + zRPC 微服务框架 |
| 数据库 | MySQL 8.0 + GORM | ORM，AutoMigrate 自动建表 |
| 配置中心 | etcd | 启动时从 etcd 加载配置，文件配置兜底 |
| 服务发现 | etcd | user-center-rpc 注册到 etcd，其他服务通过 etcd 发现 |
| 认证 | Google OAuth 2.0 + JWT | OAuth 登录，JWT 鉴权 |
| 密码加密 | bcrypt | 超管账号密码加密存储 |
| 依赖注入 | 手动 DI (wire.go) | 无代码生成，显式构造 |
| 通信协议 | gRPC (Protobuf) | 微服务间通信 |

### 前端

| 组件 | 技术选型 | 说明 |
|------|---------|------|
| 框架 | Next.js 14 (App Router) | React 服务端组件 |
| 语言 | TypeScript | 类型安全 |
| UI 组件 | shadcn/ui + Tailwind CSS | 基于 Radix/base-ui 的组件库 |
| 状态管理 | localStorage + fetch | JWT token 本地存储 |
| 表格 | shadcn Table | 数据展示 |

### 基础设施

| 组件 | 部署方式 | 数据持久化 |
|------|---------|-----------|
| MySQL 8.0 | Docker | `/Users/yangboyi/docker-data/mysql` |
| etcd 3.5 | Docker | `/Users/yangboyi/docker-data/etcd` |

## 项目结构

```
ddd-dev/
├── backend/                    # dropship-api 服务 (Go module)
├── user-center/                # 用户中心微服务 (Go module)
├── frontend/                   # Next.js 前端
└── docs/                       # 设计文档和实施计划
```

### backend/ — 代发核心服务

DDD 四层架构，依赖方向：`Server → Application → Domain → Model`

```
backend/
├── main.go                          # 入口：go-zero REST + etcd 配置 + RPC client
├── etc/config.yaml                  # 配置文件（MySQL、etcd、RPC）
├── infra/
│   ├── config/config.go             # 配置结构体
│   ├── vars/vars.go                 # 全局变量
│   └── consts/consts.go             # 常量（平台、状态枚举）
├── pkg/
│   └── response/response.go         # 统一 HTTP 响应
├── internal/
│   ├── server/                      # 接口层 — HTTP Handlers + 路由
│   ├── middleware/                   # 鉴权中间件（RPC 调 user-center）
│   ├── application/                 # 应用层 — 用例编排
│   ├── domain/                      # 领域层 — 聚合根 + 领域服务
│   │   ├── source_item/entity/      #   货源聚合根
│   │   ├── product/entity/          #   商品聚合根
│   │   ├── publish/entity/          #   发品任务聚合根
│   │   ├── platform/gateway.go      #   平台网关接口（防腐层）
│   │   └── domain_service/          #   发品领域服务
│   ├── model/                       # 模型层
│   │   ├── po/mysql/                #   持久化对象（DB 表映射）
│   │   ├── dto/                     #   数据传输对象
│   │   └── anticorruption/          #   防腐层数据结构
│   ├── repository/                  # 基础设施 — Repository GORM 实现
│   ├── gateway/                     # 基础设施 — Mock 平台网关
│   └── queries/                     # 查询层（CQRS 读侧）
└── go.mod
```

**核心领域模型：**

| 聚合根 | 职责 | 状态流转 |
|--------|------|---------|
| SourceItem（货源） | 从批发平台采集的商品信息 | New → Selected / Ignored |
| Product（商品） | 编辑后待发布的商品 | Draft → Ready → Published |
| PublishTask（发品任务） | 发布到卖货平台的任务 | Pending → Publishing → Success / Failed |

**API 端点（12 个，全部需 JWT 鉴权）：**

| 模块 | 端点 |
|------|------|
| 货源 | POST /api/source-items/import, GET /api/source-items, GET /api/source-items/detail, PUT /api/source-items/status, POST /api/source-items/tag |
| 商品 | POST /api/products/create-from-source, GET /api/products, GET /api/products/detail, PUT /api/products, PUT /api/products/ready |
| 发品 | POST /api/publish-tasks, GET /api/publish-tasks |

### user-center/ — 用户中心微服务

同一 Go module，包含 REST API 和 zRPC 两个入口。

```
user-center/
├── cmd/
│   ├── api/main.go                  # REST API 入口 :8880
│   └── rpc/main.go                  # zRPC 入口 :8881（注册 etcd）
├── etc/
│   ├── api.yaml                     # API 配置（含 Google OAuth）
│   └── rpc.yaml                     # RPC 配置
├── proto/
│   ├── usercenter.proto             # Protobuf 定义
│   └── pb/                          # 生成的 Go 代码
├── internal/
│   ├── config/                      # 配置结构体
│   ├── domain/user/                 # 用户领域（entity + repository 接口）
│   ├── model/po/mysql/              # PO：users, roles, user_roles
│   ├── model/dto/                   # DTO
│   ├── repository/                  # GORM 实现
│   ├── application/                 # JWT 工具 + UserApp 业务逻辑
│   ├── seed/                        # 预设角色 + 超管自动创建
│   ├── middleware/                   # JWT 鉴权中间件
│   └── server/
│       ├── api/                     # REST handlers + 路由
│       └── rpc/                     # zRPC 实现
└── go.mod
```

**RBAC 权限模型：**

| 角色 | 用户管理 | 角色分配 | 业务操作 |
|------|---------|---------|---------|
| super_admin | ✅ | ✅ | ✅ |
| admin | ✅ | ❌ | ✅ |
| operator | ❌ | ❌ | ✅ |
| viewer（默认） | ❌ | ❌ | 只读 |

**RPC 接口：**

| 方法 | 说明 | 调用方 |
|------|------|--------|
| VerifyToken | JWT 验证，返回用户信息和角色 | dropship-api 鉴权中间件 |
| CheckRole | 检查用户是否有指定角色 | dropship-api |
| GetUserInfo | 获取用户详情 | dropship-api |

**REST API 端点（8 个）：**

| 端点 | 鉴权 | 说明 |
|------|------|------|
| GET /api/auth/google/login | 无 | 发起 Google OAuth |
| GET /api/auth/google/callback | 无 | OAuth 回调 |
| POST /api/auth/login | 无 | 账号密码登录（超管） |
| GET /api/auth/me | JWT | 当前用户信息 |
| GET /api/init/check | 无 | 检查是否需要初始化超管 |
| POST /api/init/super-admin | JWT | 初始化超管 |
| GET /api/users | JWT + admin | 用户列表 |
| PUT /api/users/status | JWT + admin | 启用/禁用用户 |
| PUT /api/users/role | JWT + super_admin | 分配角色 |

### frontend/ — 管理面板

```
frontend/src/
├── app/
│   ├── layout.tsx                   # 根布局（AuthGuard 全局鉴权）
│   ├── page.tsx                     # 首页（重定向到 /sources）
│   ├── login/page.tsx               # 登录页（账号密码 + Google OAuth）
│   ├── login/callback/page.tsx      # OAuth 回调处理
│   ├── sources/page.tsx             # 货源管理
│   ├── products/page.tsx            # 商品管理
│   ├── publish/page.tsx             # 发品任务
│   ├── users/page.tsx               # 用户管理（admin+）
│   ├── init/page.tsx                # 超管初始化
│   └── profile/page.tsx             # 个人中心
├── components/
│   ├── layout/
│   │   ├── sidebar.tsx              # 侧边栏（角色权限控制菜单可见性）
│   │   └── auth-guard.tsx           # 鉴权守卫（未登录跳 /login）
│   ├── sources/                     # 货源组件（导入对话框、列表表格）
│   ├── products/                    # 商品组件（编辑对话框、列表表格）
│   ├── publish/                     # 发品组件（任务表格）
│   └── users/                       # 用户组件（用户表格 + 角色分配）
└── lib/
    ├── api.ts                       # API 请求封装（双后端、自动带 JWT）
    ├── auth.ts                      # Token/User 本地存储管理
    └── utils.ts                     # 工具函数
```

## 架构设计原则

### DDD 分层与依赖方向

```
Server (HTTP Handler)
    ↓ 调用
Application (用例编排)
    ↓ 调用
Domain (聚合根 + 领域服务)
    ↓ 引用
Model (PO / DTO / 值对象)

Infrastructure (Repository 实现 / Gateway 实现)
    ↑ 实现 Domain 层定义的接口
```

- **禁止反向依赖**：下层不能引用上层
- **依赖倒置**：Domain 层定义 Repository/Gateway 接口，Infrastructure 层实现
- **防腐层**：外部平台 API 通过 Gateway 接口隔离，当前用 Mock 实现

### CQRS 思路

- **写操作**：通过聚合根进行，保证业务规则
- **读操作**：查询层（queries/）直接查 DB，不走聚合根，提高查询灵活性

### 微服务通信

- 服务间通过 **gRPC (zRPC)** 通信
- 使用 **etcd** 做服务注册与发现
- dropship-api 的鉴权中间件通过 RPC 调用 user-center-rpc 验证 JWT

## 快速启动

### 前置依赖

- Go 1.25+
- Node.js 20+
- Docker（MySQL + etcd）

### 启动基础设施

```bash
# MySQL（数据持久化到本地）
docker run -d --name mysql8 \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -v ~/docker-data/mysql:/var/lib/mysql \
  mysql:8.0 --lower-case-table-names=2

# etcd（数据持久化到本地）
docker run -d --name etcd \
  -p 2379:2379 -p 2380:2380 \
  -v ~/docker-data/etcd:/etcd-data \
  quay.io/coreos/etcd:v3.5.17 \
  /usr/local/bin/etcd \
  --data-dir /etcd-data \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-client-urls http://0.0.0.0:2379

# 创建数据库
docker exec mysql8 mysql -uroot -proot123 \
  -e "CREATE DATABASE IF NOT EXISTS dropship CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 写入 etcd 配置
docker exec etcd etcdctl put /dropship/config \
  '{"Name":"dropship-api","Host":"0.0.0.0","Port":8888,"MySQL":{"DataSource":"root:root123@tcp(127.0.0.1:3306)/dropship?charset=utf8mb4&parseTime=True&loc=Local"}}'
```

### 启动后端服务

```bash
# 1. user-center-rpc（先启动，注册 etcd）
cd user-center && go run cmd/rpc/main.go -f etc/rpc.yaml
# 首次启动自动创建超管：admin / admin123

# 2. user-center-api
cd user-center && go run cmd/api/main.go -f etc/api.yaml

# 3. dropship-api
cd backend && go run main.go -f etc/config.yaml
```

### 启动前端

```bash
cd frontend && npm install && npm run dev
```

### 访问

- 前端：http://localhost:3000
- 超管登录：admin / admin123
- Google OAuth：需在 Google Cloud Console 配置 OAuth Client ID

## 数据库表

| 表名 | 所属服务 | 说明 |
|------|---------|------|
| source_items | dropship-api | 货源池 |
| products | dropship-api | 商品 |
| product_skus | dropship-api | 商品 SKU |
| publish_tasks | dropship-api | 发品任务 |
| users | user-center | 用户（支持 OAuth + 密码） |
| roles | user-center | 角色（4 个预设） |
| user_roles | user-center | 用户角色关联 |

所有表共用 `dropship` 数据库，通过 GORM AutoMigrate 自动建表。

## 后续规划

- [ ] 对接真实 1688 / 拼多多 API（替换 Mock Gateway）
- [ ] 订单监控：监控卖货平台下单情况
- [ ] 自动采购：自动/手动触发批发平台下单
- [ ] 物流同步：发货信息 + 退货地址同步到卖货平台
- [ ] 微服务拆分：业务复杂后按领域拆分独立服务
