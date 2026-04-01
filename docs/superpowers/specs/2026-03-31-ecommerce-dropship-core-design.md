# 电商一键代发工具 — 核心链路设计文档

## 概述

构建一个电商运营工具，核心能力是「一键代发」：从批发平台（1688）选品，编辑后发布到卖货平台（拼多多）。本文档覆盖第一期核心链路：**货源采集 → 筛选选品 → 发品**。

**技术栈**：Go-zero (后端) + Next.js 14 (前端) + MySQL (数据库)  
**架构风格**：DDD + Clean Architecture，经典四层分层，单体服务先行

---

## 领域模型

### 聚合划分

三个核心聚合，各自独立生命周期：

#### 1. SourceItem（货源）

从批发平台采集的原始商品信息，作为选品池。

- **聚合根**: SourceItem
- **值对象**: Price (区间价/阶梯价)、Supplier (名称/评分/地区)、Category、Image
- **状态**: New → Selected / Ignored
- **行为**: Import(), Select(), Ignore(), AddTag(), Refresh()
- **筛选维度**: 平台、分类、价格区间、供应商评分/地区、销量、起订量、标签、状态、导入时间

#### 2. Product（商品）

从货源创建的待发布商品，允许编辑标题/图片/描述/定价。

- **聚合根**: Product
- **实体**: SKU (规格+库存)
- **值对象**: Price (成本价/售价/利润)、Image
- **关联**: SourceItemID
- **状态**: Draft → Ready → Published
- **行为**: CreateFromSource(), EditInfo(), SetPrice(), MarkReady()

#### 3. PublishTask（发品任务）

将商品发布到目标平台的任务记录。

- **聚合根**: PublishTask
- **值对象**: PublishConfig (目标分类/运费模板等)、TargetPlatform
- **关联**: ProductID
- **状态**: Pending → Publishing → Success / Failed
- **行为**: Create(), MarkPublishing(), MarkSuccess(), MarkFailed()

### 领域服务

**PublishDomainService** — 协调 Product 和 PublishTask 的跨聚合操作：
- `PublishProduct(productID, platform, config)`: 校验商品状态 → 创建发品任务 → 调用平台适配器 → 更新状态

### 防腐层（Anti-Corruption Layer）

平台网关接口定义在 Domain 层，实现在 Infrastructure 层：

```go
// domain/platform/gateway.go
type SourcePlatformGateway interface {
    FetchProduct(ctx context.Context, sourceURL string) (*SourceProduct, error)
}

type TargetPlatformGateway interface {
    PublishProduct(ctx context.Context, product *Product, config *PublishConfig) (*PublishResult, error)
}
```

当前阶段使用 Mock 实现，后续对接真实 API。

---

## 数据流

```
货源采集: 用户输入URL → SourcePlatformGateway.FetchProduct(Mock) → 创建SourceItem → 存入MySQL
筛选选品: 用户筛选货源池 → 选中货源 → Product.CreateFromSource() → 编辑信息 → MarkReady
发品:     用户点击发布 → PublishDomainService → TargetPlatformGateway.Publish(Mock) → 更新PublishTask状态
```

---

## 技术栈

### 后端
| 组件 | 选型 |
|------|------|
| 框架 | Go-zero REST |
| 数据库 | MySQL + GORM |
| 依赖注入 | Google Wire |
| 配置 | go-zero 内置 conf |
| 日志 | go-zero logx |
| API 定义 | go-zero .api 文件 |

### 前端
| 组件 | 选型 |
|------|------|
| 框架 | Next.js 14 (App Router) |
| 语言 | TypeScript |
| UI | shadcn/ui + Tailwind CSS |
| 状态 | React Server Components + SWR |
| 表格 | TanStack Table |

---

## 后端目录结构

```
backend/
├── main.go
├── etc/config.yaml
├── infra/
│   ├── config/config.go
│   ├── vars/vars.go
│   └── consts/consts.go
├── pkg/
│   └── response/response.go
├── internal/
│   ├── server/                      # 接口层
│   │   ├── routes.go
│   │   ├── source_item_handler.go
│   │   ├── product_handler.go
│   │   └── publish_handler.go
│   ├── application/                 # 应用层
│   │   ├── source_item_app.go
│   │   ├── product_app.go
│   │   └── publish_app.go
│   ├── domain/                      # 领域层
│   │   ├── source_item/
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── product/
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── publish/
│   │   │   ├── entity/
│   │   │   └── repository/
│   │   ├── platform/
│   │   │   └── gateway.go
│   │   └── domain_service/
│   │       └── publish_service.go
│   ├── model/                       # 模型层
│   │   ├── po/mysql/
│   │   ├── dto/
│   │   └── anticorruption/
│   ├── repository/                  # 基础设施 - Repository实现
│   │   ├── source_item_repo.go
│   │   ├── product_repo.go
│   │   └── publish_task_repo.go
│   ├── gateway/                     # 基础设施 - Gateway实现
│   │   ├── mock_source_gateway.go
│   │   └── mock_target_gateway.go
│   ├── queries/                     # 查询层(CQRS读侧)
│   │   └── source_item_query.go
│   ├── wire.go
│   └── wire_gen.go
└── go.mod
```

依赖方向: `server → application → domain → model`，禁止反向依赖。Repository/Gateway 接口定义在 domain 层，实现在 infrastructure 层（repository/、gateway/ 目录）。

---

## 数据库表设计

### source_items
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| platform | VARCHAR(32) | 货源平台 (Ali1688等) |
| source_url | VARCHAR(512) | 原始链接 |
| external_id | VARCHAR(128) | 平台商品ID |
| title | VARCHAR(256) | 标题 |
| description | TEXT | 描述 |
| images | JSON | 图片列表 |
| price_min | DECIMAL(10,2) | 最低价 |
| price_max | DECIMAL(10,2) | 最高价 |
| supplier | JSON | 供应商信息 |
| category | VARCHAR(128) | 分类 |
| tags | JSON | 用户标签 |
| sales_volume | INT | 销量 |
| min_order | INT | 起订量 |
| status | VARCHAR(32) | 状态 |
| fetched_at | DATETIME | 采集时间 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### products
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| source_item_id | BIGINT FK | 关联货源 |
| name | VARCHAR(256) | 商品名称 |
| description | TEXT | 描述 |
| images | JSON | 图片列表 |
| cost_price | DECIMAL(10,2) | 成本价 |
| sell_price | DECIMAL(10,2) | 售价 |
| category_id | VARCHAR(128) | 分类 |
| status | VARCHAR(32) | 状态 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### product_skus
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| product_id | BIGINT FK | 关联商品 |
| spec_name | VARCHAR(128) | 规格名称 (如: 颜色/尺码) |
| spec_value | VARCHAR(128) | 规格值 (如: 红色/XL) |
| price | DECIMAL(10,2) | SKU价格 |
| stock | INT | 库存 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### publish_tasks
| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT PK | 主键 |
| product_id | BIGINT FK | 关联商品 |
| target_platform | VARCHAR(32) | 目标平台 |
| platform_product_id | VARCHAR(128) | 平台商品ID |
| publish_config | JSON | 发布配置 |
| status | VARCHAR(32) | 状态 |
| error_message | TEXT | 错误信息 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

---

## 前端页面

### 货源管理页
- 导入货源（输入 URL）
- 货源列表（筛选/搜索/标签）
- 货源详情
- 批量选品操作

### 商品管理页
- 商品列表（按状态筛选）
- 商品编辑（标题/图片/描述/定价）
- SKU 管理
- 一键发布操作

### 发品任务页
- 任务列表（状态/平台筛选）
- 任务详情（成功/失败信息）
- 重试失败任务

---

## 验证方案

1. **后端**: `go build` 编译通过 → 启动服务 → 通过 curl/Postman 调用 API 验证完整链路
2. **数据库**: 验证 3 张表的 CRUD 操作，检查数据一致性
3. **前端**: `npm run dev` 启动 → 手动操作 3 个页面的核心流程
4. **端到端**: 导入货源 → 筛选选品 → 创建商品 → 编辑 → 发品 → 查看任务状态
