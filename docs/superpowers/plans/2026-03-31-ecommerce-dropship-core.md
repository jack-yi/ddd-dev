# 电商一键代发 — 核心链路实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建「货源采集 → 筛选选品 → 发品」核心链路的完整后端 API + 前端管理面板

**Architecture:** Go-zero 单体 REST 服务，经典 DDD 四层分层（Server → Application → Domain → Model），依赖方向严格单向。Repository/Gateway 接口定义在 Domain 层，实现在 Infrastructure 层。前端 Next.js 14 App Router + shadcn/ui 管理面板。

**Tech Stack:** Go 1.22+, Go-zero, GORM, Google Wire, MySQL | Next.js 14, TypeScript, shadcn/ui, Tailwind CSS, TanStack Table, SWR

**Spec:** `docs/superpowers/specs/2026-03-31-ecommerce-dropship-core-design.md`

---

## File Structure

### Backend (`backend/`)

```
backend/
├── main.go                                    # 入口：初始化配置、Wire注入、启动go-zero server
├── etc/config.yaml                            # go-zero 配置(端口、数据库连接等)
├── go.mod                                     # Go module
├── infra/
│   ├── config/config.go                       # 配置结构体(映射config.yaml)
│   ├── vars/vars.go                           # 全局变量(DB实例)
│   └── consts/consts.go                       # 全局常量
├── pkg/
│   └── response/response.go                   # 统一 HTTP 响应格式
├── internal/
│   ├── model/
│   │   ├── po/mysql/source_item.go            # SourceItem PO
│   │   ├── po/mysql/product.go                # Product PO
│   │   ├── po/mysql/product_sku.go            # ProductSKU PO
│   │   ├── po/mysql/publish_task.go           # PublishTask PO
│   │   ├── dto/source_item_dto.go             # 货源请求/响应 DTO
│   │   ├── dto/product_dto.go                 # 商品请求/响应 DTO
│   │   ├── dto/publish_dto.go                 # 发品请求/响应 DTO
│   │   └── anticorruption/platform_types.go   # 平台API数据结构
│   ├── domain/
│   │   ├── source_item/
│   │   │   ├── entity/source_item.go          # SourceItem 聚合根+值对象
│   │   │   └── repository/repository.go       # Repository 接口
│   │   ├── product/
│   │   │   ├── entity/product.go              # Product 聚合根+SKU实体+值对象
│   │   │   └── repository/repository.go       # Repository 接口
│   │   ├── publish/
│   │   │   ├── entity/publish_task.go         # PublishTask 聚合根
│   │   │   └── repository/repository.go       # Repository 接口
│   │   ├── platform/gateway.go                # 平台网关接口(防腐层)
│   │   └── domain_service/publish_service.go  # 发品领域服务
│   ├── repository/
│   │   ├── source_item_repo.go                # SourceItem GORM 实现
│   │   ├── product_repo.go                    # Product GORM 实现
│   │   └── publish_task_repo.go               # PublishTask GORM 实现
│   ├── gateway/
│   │   ├── mock_source_gateway.go             # Mock 货源平台
│   │   └── mock_target_gateway.go             # Mock 目标平台
│   ├── queries/
│   │   └── source_item_query.go               # 货源筛选查询
│   ├── application/
│   │   ├── source_item_app.go                 # 货源用例
│   │   ├── product_app.go                     # 商品用例
│   │   └── publish_app.go                     # 发品用例
│   ├── server/
│   │   ├── routes.go                          # 路由注册
│   │   ├── source_item_handler.go             # 货源 handler
│   │   ├── product_handler.go                 # 商品 handler
│   │   └── publish_handler.go                 # 发品 handler
│   ├── wire.go                                # Wire 依赖注入定义
│   └── wire_gen.go                            # Wire 生成代码
```

### Frontend (`frontend/`)

```
frontend/
├── package.json
├── next.config.ts
├── tsconfig.json
├── tailwind.config.ts
├── postcss.config.js
├── src/
│   ├── app/
│   │   ├── layout.tsx                         # 根布局(侧边栏导航)
│   │   ├── page.tsx                           # 首页(重定向到货源)
│   │   ├── sources/
│   │   │   └── page.tsx                       # 货源管理页
│   │   ├── products/
│   │   │   └── page.tsx                       # 商品管理页
│   │   └── publish/
│   │       └── page.tsx                       # 发品任务页
│   ├── components/
│   │   ├── layout/sidebar.tsx                 # 侧边栏导航
│   │   ├── sources/source-table.tsx           # 货源列表表格
│   │   ├── sources/import-dialog.tsx          # 导入货源对话框
│   │   ├── products/product-table.tsx         # 商品列表表格
│   │   ├── products/product-edit-dialog.tsx   # 商品编辑对话框
│   │   └── publish/publish-table.tsx          # 发品任务表格
│   └── lib/
│       └── api.ts                             # API 请求封装
```

---

## Task 1: 后端项目初始化

**Files:**
- Create: `backend/go.mod`
- Create: `backend/main.go`
- Create: `backend/etc/config.yaml`
- Create: `backend/infra/config/config.go`
- Create: `backend/infra/vars/vars.go`
- Create: `backend/infra/consts/consts.go`
- Create: `backend/pkg/response/response.go`

- [ ] **Step 1: 初始化 Go module**

```bash
cd backend
go mod init github.com/yangboyi/ddd-dev/backend
```

- [ ] **Step 2: 创建配置文件 `backend/etc/config.yaml`**

```yaml
Name: dropship-api
Host: 0.0.0.0
Port: 8888

MySQL:
  DataSource: "root:123456@tcp(127.0.0.1:3306)/dropship?charset=utf8mb4&parseTime=True&loc=Local"
```

- [ ] **Step 3: 创建配置结构体 `backend/infra/config/config.go`**

```go
package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQL MySQLConfig
}

type MySQLConfig struct {
	DataSource string
}
```

- [ ] **Step 4: 创建全局变量 `backend/infra/vars/vars.go`**

```go
package vars

import "gorm.io/gorm"

var DB *gorm.DB
```

- [ ] **Step 5: 创建全局常量 `backend/infra/consts/consts.go`**

```go
package consts

const (
	PlatformAli1688 = "ali1688"
	PlatformPDD     = "pdd"

	SourceItemStatusNew      = "new"
	SourceItemStatusSelected = "selected"
	SourceItemStatusIgnored  = "ignored"

	ProductStatusDraft     = "draft"
	ProductStatusReady     = "ready"
	ProductStatusPublished = "published"

	PublishTaskStatusPending    = "pending"
	PublishTaskStatusPublishing = "publishing"
	PublishTaskStatusSuccess    = "success"
	PublishTaskStatusFailed     = "failed"
)
```

- [ ] **Step 6: 创建统一响应 `backend/pkg/response/response.go`**

```go
package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, &Body{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, msg string) {
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code:    code,
		Message: msg,
	})
}
```

- [ ] **Step 7: 创建入口 `backend/main.go`**

```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/backend/infra/config"
	"github.com/yangboyi/ddd-dev/backend/infra/vars"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}
	vars.DB = db

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// TODO: 注册路由（Task 8 完成后替换）

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```

- [ ] **Step 8: 安装依赖**

```bash
cd backend
go get github.com/zeromicro/go-zero@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/mysql@latest
go mod tidy
```

- [ ] **Step 9: 验证编译通过**

```bash
cd backend && go build ./...
```

Expected: 编译成功，无报错

- [ ] **Step 10: Commit**

```bash
git add backend/
git commit -m "feat: init backend project with go-zero, gorm, config"
```

---

## Task 2: 模型层 — PO 和 DTO

**Files:**
- Create: `backend/internal/model/po/mysql/source_item.go`
- Create: `backend/internal/model/po/mysql/product.go`
- Create: `backend/internal/model/po/mysql/product_sku.go`
- Create: `backend/internal/model/po/mysql/publish_task.go`
- Create: `backend/internal/model/dto/source_item_dto.go`
- Create: `backend/internal/model/dto/product_dto.go`
- Create: `backend/internal/model/dto/publish_dto.go`
- Create: `backend/internal/model/anticorruption/platform_types.go`

- [ ] **Step 1: 创建 SourceItem PO `backend/internal/model/po/mysql/source_item.go`**

```go
package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type SupplierInfo struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	Region string  `json:"region"`
}

func (s SupplierInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SupplierInfo) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type SourceItem struct {
	ID          int64        `gorm:"primaryKey;autoIncrement"`
	Platform    string       `gorm:"type:varchar(32);not null"`
	SourceURL   string       `gorm:"type:varchar(512);not null"`
	ExternalID  string       `gorm:"type:varchar(128)"`
	Title       string       `gorm:"type:varchar(256)"`
	Description string       `gorm:"type:text"`
	Images      StringSlice  `gorm:"type:json"`
	PriceMin    float64      `gorm:"type:decimal(10,2)"`
	PriceMax    float64      `gorm:"type:decimal(10,2)"`
	Supplier    SupplierInfo `gorm:"type:json"`
	Category    string       `gorm:"type:varchar(128)"`
	Tags        StringSlice  `gorm:"type:json"`
	SalesVolume int          `gorm:"type:int"`
	MinOrder    int          `gorm:"type:int"`
	Status      string       `gorm:"type:varchar(32);not null;default:'new'"`
	FetchedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SourceItem) TableName() string {
	return "source_items"
}
```

- [ ] **Step 2: 创建 Product PO `backend/internal/model/po/mysql/product.go`**

```go
package mysql

import "time"

type Product struct {
	ID           int64       `gorm:"primaryKey;autoIncrement"`
	SourceItemID int64       `gorm:"type:bigint"`
	Name         string      `gorm:"type:varchar(256);not null"`
	Description  string      `gorm:"type:text"`
	Images       StringSlice `gorm:"type:json"`
	CostPrice    float64     `gorm:"type:decimal(10,2)"`
	SellPrice    float64     `gorm:"type:decimal(10,2)"`
	CategoryID   string      `gorm:"type:varchar(128)"`
	Status       string      `gorm:"type:varchar(32);not null;default:'draft'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Product) TableName() string {
	return "products"
}
```

- [ ] **Step 3: 创建 ProductSKU PO `backend/internal/model/po/mysql/product_sku.go`**

```go
package mysql

import "time"

type ProductSKU struct {
	ID        int64   `gorm:"primaryKey;autoIncrement"`
	ProductID int64   `gorm:"type:bigint;not null"`
	SpecName  string  `gorm:"type:varchar(128)"`
	SpecValue string  `gorm:"type:varchar(128)"`
	Price     float64 `gorm:"type:decimal(10,2)"`
	Stock     int     `gorm:"type:int"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ProductSKU) TableName() string {
	return "product_skus"
}
```

- [ ] **Step 4: 创建 PublishTask PO `backend/internal/model/po/mysql/publish_task.go`**

```go
package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PublishConfig struct {
	CategoryID      string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}

func (p PublishConfig) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PublishConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

type PublishTask struct {
	ID                int64         `gorm:"primaryKey;autoIncrement"`
	ProductID         int64         `gorm:"type:bigint;not null"`
	TargetPlatform    string        `gorm:"type:varchar(32);not null"`
	PlatformProductID string        `gorm:"type:varchar(128)"`
	PublishConfig     PublishConfig `gorm:"type:json"`
	Status            string        `gorm:"type:varchar(32);not null;default:'pending'"`
	ErrorMessage      string        `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (PublishTask) TableName() string {
	return "publish_tasks"
}
```

- [ ] **Step 5: 创建防腐层类型 `backend/internal/model/anticorruption/platform_types.go`**

```go
package anticorruption

// SourceProduct 货源平台返回的商品结构（防腐层，隔离外部API变化）
type SourceProduct struct {
	ExternalID  string   `json:"externalId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	PriceMin    float64  `json:"priceMin"`
	PriceMax    float64  `json:"priceMax"`
	Supplier    Supplier `json:"supplier"`
	Category    string   `json:"category"`
	SalesVolume int      `json:"salesVolume"`
	MinOrder    int      `json:"minOrder"`
}

type Supplier struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	Region string  `json:"region"`
}

// PublishResult 目标平台发布结果
type PublishResult struct {
	PlatformProductID string `json:"platformProductId"`
	Success           bool   `json:"success"`
	ErrorMessage      string `json:"errorMessage"`
}
```

- [ ] **Step 6: 创建 SourceItem DTO `backend/internal/model/dto/source_item_dto.go`**

```go
package dto

type ImportSourceItemReq struct {
	Platform  string `json:"platform"`
	SourceURL string `json:"sourceUrl"`
}

type SourceItemResp struct {
	ID          int64    `json:"id"`
	Platform    string   `json:"platform"`
	SourceURL   string   `json:"sourceUrl"`
	ExternalID  string   `json:"externalId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	PriceMin    float64  `json:"priceMin"`
	PriceMax    float64  `json:"priceMax"`
	Supplier    struct {
		Name   string  `json:"name"`
		Rating float64 `json:"rating"`
		Region string  `json:"region"`
	} `json:"supplier"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	SalesVolume int      `json:"salesVolume"`
	MinOrder    int      `json:"minOrder"`
	Status      string   `json:"status"`
	FetchedAt   string   `json:"fetchedAt"`
	CreatedAt   string   `json:"createdAt"`
}

type SourceItemFilter struct {
	Platform    *string  `json:"platform,omitempty"`
	Category    *string  `json:"category,omitempty"`
	PriceMin    *float64 `json:"priceMin,omitempty"`
	PriceMax    *float64 `json:"priceMax,omitempty"`
	SupplierMin *float64 `json:"supplierMin,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Keyword     *string  `json:"keyword,omitempty"`
	Page        int      `json:"page"`
	PageSize    int      `json:"pageSize"`
}

type UpdateSourceItemStatusReq struct {
	Status string `json:"status"`
}

type AddTagReq struct {
	Tag string `json:"tag"`
}
```

- [ ] **Step 7: 创建 Product DTO `backend/internal/model/dto/product_dto.go`**

```go
package dto

type CreateProductFromSourceReq struct {
	SourceItemID int64 `json:"sourceItemId"`
}

type UpdateProductReq struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Images      []string    `json:"images,omitempty"`
	CostPrice   *float64    `json:"costPrice,omitempty"`
	SellPrice   *float64    `json:"sellPrice,omitempty"`
	CategoryID  *string     `json:"categoryId,omitempty"`
	SKUs        []SKUItem   `json:"skus,omitempty"`
}

type SKUItem struct {
	ID        int64   `json:"id,omitempty"`
	SpecName  string  `json:"specName"`
	SpecValue string  `json:"specValue"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type ProductResp struct {
	ID           int64     `json:"id"`
	SourceItemID int64     `json:"sourceItemId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Images       []string  `json:"images"`
	CostPrice    float64   `json:"costPrice"`
	SellPrice    float64   `json:"sellPrice"`
	CategoryID   string    `json:"categoryId"`
	Status       string    `json:"status"`
	SKUs         []SKUItem `json:"skus"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

type ProductFilter struct {
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}
```

- [ ] **Step 8: 创建 Publish DTO `backend/internal/model/dto/publish_dto.go`**

```go
package dto

type CreatePublishTaskReq struct {
	ProductID      int64  `json:"productId"`
	TargetPlatform string `json:"targetPlatform"`
	CategoryID     string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}

type PublishTaskResp struct {
	ID                int64  `json:"id"`
	ProductID         int64  `json:"productId"`
	TargetPlatform    string `json:"targetPlatform"`
	PlatformProductID string `json:"platformProductId"`
	Status            string `json:"status"`
	ErrorMessage      string `json:"errorMessage"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type PublishTaskFilter struct {
	Status         *string `json:"status,omitempty"`
	TargetPlatform *string `json:"targetPlatform,omitempty"`
	Page           int     `json:"page"`
	PageSize       int     `json:"pageSize"`
}
```

- [ ] **Step 9: 验证编译**

```bash
cd backend && go build ./...
```

Expected: 编译成功

- [ ] **Step 10: Commit**

```bash
git add backend/internal/model/
git commit -m "feat: add PO, DTO, and anticorruption layer models"
```

---

## Task 3: 领域层 — 聚合根、值对象、Repository 接口、Gateway 接口

**Files:**
- Create: `backend/internal/domain/source_item/entity/source_item.go`
- Create: `backend/internal/domain/source_item/repository/repository.go`
- Create: `backend/internal/domain/product/entity/product.go`
- Create: `backend/internal/domain/product/repository/repository.go`
- Create: `backend/internal/domain/publish/entity/publish_task.go`
- Create: `backend/internal/domain/publish/repository/repository.go`
- Create: `backend/internal/domain/platform/gateway.go`
- Create: `backend/internal/domain/domain_service/publish_service.go`

- [ ] **Step 1: 创建 SourceItem 聚合根 `backend/internal/domain/source_item/entity/source_item.go`**

```go
package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
)

type Price struct {
	Min float64
	Max float64
}

type Supplier struct {
	Name   string
	Rating float64
	Region string
}

type SourceItem struct {
	ID          int64
	Platform    string
	SourceURL   string
	ExternalID  string
	Title       string
	Description string
	Images      []string
	Price       Price
	Supplier    Supplier
	Category    string
	Tags        []string
	SalesVolume int
	MinOrder    int
	Status      string
	FetchedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSourceItem(platform, sourceURL, externalID, title, description string,
	images []string, price Price, supplier Supplier, category string,
	salesVolume, minOrder int) *SourceItem {
	now := time.Now()
	return &SourceItem{
		Platform:    platform,
		SourceURL:   sourceURL,
		ExternalID:  externalID,
		Title:       title,
		Description: description,
		Images:      images,
		Price:       price,
		Supplier:    supplier,
		Category:    category,
		Tags:        []string{},
		SalesVolume: salesVolume,
		MinOrder:    minOrder,
		Status:      consts.SourceItemStatusNew,
		FetchedAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *SourceItem) Select() error {
	if s.Status == consts.SourceItemStatusSelected {
		return errors.New("source item already selected")
	}
	s.Status = consts.SourceItemStatusSelected
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SourceItem) Ignore() error {
	if s.Status == consts.SourceItemStatusIgnored {
		return errors.New("source item already ignored")
	}
	s.Status = consts.SourceItemStatusIgnored
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SourceItem) AddTag(tag string) {
	for _, t := range s.Tags {
		if t == tag {
			return
		}
	}
	s.Tags = append(s.Tags, tag)
	s.UpdatedAt = time.Now()
}
```

- [ ] **Step 2: 创建 SourceItem Repository 接口 `backend/internal/domain/source_item/repository/repository.go`**

```go
package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
)

type SourceItemRepository interface {
	Save(ctx context.Context, item *entity.SourceItem) error
	FindByID(ctx context.Context, id int64) (*entity.SourceItem, error)
	Update(ctx context.Context, item *entity.SourceItem) error
}
```

- [ ] **Step 3: 创建 Product 聚合根 `backend/internal/domain/product/entity/product.go`**

```go
package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
	sourceEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
)

type SKU struct {
	ID        int64
	SpecName  string
	SpecValue string
	Price     float64
	Stock     int
}

type Product struct {
	ID           int64
	SourceItemID int64
	Name         string
	Description  string
	Images       []string
	CostPrice    float64
	SellPrice    float64
	CategoryID   string
	Status       string
	SKUs         []SKU
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateFromSource(source *sourceEntity.SourceItem) *Product {
	now := time.Now()
	return &Product{
		SourceItemID: source.ID,
		Name:         source.Title,
		Description:  source.Description,
		Images:       source.Images,
		CostPrice:    source.Price.Min,
		SellPrice:    0,
		CategoryID:   source.Category,
		Status:       consts.ProductStatusDraft,
		SKUs:         []SKU{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (p *Product) EditInfo(name, description *string, images []string, categoryID *string) {
	if name != nil {
		p.Name = *name
	}
	if description != nil {
		p.Description = *description
	}
	if images != nil {
		p.Images = images
	}
	if categoryID != nil {
		p.CategoryID = *categoryID
	}
	p.UpdatedAt = time.Now()
}

func (p *Product) SetPrice(costPrice, sellPrice *float64) {
	if costPrice != nil {
		p.CostPrice = *costPrice
	}
	if sellPrice != nil {
		p.SellPrice = *sellPrice
	}
	p.UpdatedAt = time.Now()
}

func (p *Product) SetSKUs(skus []SKU) {
	p.SKUs = skus
	p.UpdatedAt = time.Now()
}

func (p *Product) MarkReady() error {
	if p.SellPrice <= 0 {
		return errors.New("sell price must be set before marking ready")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	p.Status = consts.ProductStatusReady
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) MarkPublished() {
	p.Status = consts.ProductStatusPublished
	p.UpdatedAt = time.Now()
}

func (p *Product) IsReady() bool {
	return p.Status == consts.ProductStatusReady
}
```

- [ ] **Step 4: 创建 Product Repository 接口 `backend/internal/domain/product/repository/repository.go`**

```go
package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id int64) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
}
```

- [ ] **Step 5: 创建 PublishTask 聚合根 `backend/internal/domain/publish/entity/publish_task.go`**

```go
package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
)

type PublishConfig struct {
	CategoryID      string
	FreightTemplate string
}

type PublishTask struct {
	ID                int64
	ProductID         int64
	TargetPlatform    string
	PlatformProductID string
	PublishConfig     PublishConfig
	Status            string
	ErrorMessage      string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewPublishTask(productID int64, targetPlatform string, config PublishConfig) *PublishTask {
	now := time.Now()
	return &PublishTask{
		ProductID:      productID,
		TargetPlatform: targetPlatform,
		PublishConfig:  config,
		Status:         consts.PublishTaskStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (t *PublishTask) MarkPublishing() error {
	if t.Status != consts.PublishTaskStatusPending {
		return errors.New("task must be pending to start publishing")
	}
	t.Status = consts.PublishTaskStatusPublishing
	t.UpdatedAt = time.Now()
	return nil
}

func (t *PublishTask) MarkSuccess(platformProductID string) {
	t.Status = consts.PublishTaskStatusSuccess
	t.PlatformProductID = platformProductID
	t.UpdatedAt = time.Now()
}

func (t *PublishTask) MarkFailed(errMsg string) {
	t.Status = consts.PublishTaskStatusFailed
	t.ErrorMessage = errMsg
	t.UpdatedAt = time.Now()
}
```

- [ ] **Step 6: 创建 PublishTask Repository 接口 `backend/internal/domain/publish/repository/repository.go`**

```go
package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
)

type PublishTaskRepository interface {
	Save(ctx context.Context, task *entity.PublishTask) error
	FindByID(ctx context.Context, id int64) (*entity.PublishTask, error)
	Update(ctx context.Context, task *entity.PublishTask) error
}
```

- [ ] **Step 7: 创建平台网关接口 `backend/internal/domain/platform/gateway.go`**

```go
package platform

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type SourcePlatformGateway interface {
	FetchProduct(ctx context.Context, sourceURL string) (*anticorruption.SourceProduct, error)
}

type TargetPlatformGateway interface {
	PublishProduct(ctx context.Context, name, description string, images []string,
		sellPrice float64, config anticorruption.PublishConfig) (*anticorruption.PublishResult, error)
}
```

- [ ] **Step 8: 更新防腐层类型，增加 PublishConfig `backend/internal/model/anticorruption/platform_types.go`**

在已有文件末尾追加：

```go
// PublishConfig 发布配置（传给目标平台的参数）
type PublishConfig struct {
	CategoryID      string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}
```

- [ ] **Step 9: 创建发品领域服务 `backend/internal/domain/domain_service/publish_service.go`**

```go
package domainservice

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/platform"
	productEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	publishEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	publishRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type PublishDomainService struct {
	productRepo   productRepo.ProductRepository
	publishRepo   publishRepo.PublishTaskRepository
	targetGateway platform.TargetPlatformGateway
}

func NewPublishDomainService(
	pr productRepo.ProductRepository,
	ptr publishRepo.PublishTaskRepository,
	tg platform.TargetPlatformGateway,
) *PublishDomainService {
	return &PublishDomainService{
		productRepo:   pr,
		publishRepo:   ptr,
		targetGateway: tg,
	}
}

func (s *PublishDomainService) PublishProduct(ctx context.Context, product *productEntity.Product,
	targetPlatform string, config publishEntity.PublishConfig) (*publishEntity.PublishTask, error) {

	if !product.IsReady() {
		return nil, fmt.Errorf("product %d is not ready for publishing", product.ID)
	}

	task := publishEntity.NewPublishTask(product.ID, targetPlatform, config)
	if err := s.publishRepo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save publish task: %w", err)
	}

	if err := task.MarkPublishing(); err != nil {
		return nil, fmt.Errorf("mark publishing: %w", err)
	}

	acConfig := anticorruption.PublishConfig{
		CategoryID:      config.CategoryID,
		FreightTemplate: config.FreightTemplate,
	}
	result, err := s.targetGateway.PublishProduct(ctx, product.Name, product.Description,
		product.Images, product.SellPrice, acConfig)
	if err != nil {
		task.MarkFailed(err.Error())
		_ = s.publishRepo.Update(ctx, task)
		return task, nil
	}

	if result.Success {
		task.MarkSuccess(result.PlatformProductID)
		product.MarkPublished()
		_ = s.productRepo.Update(ctx, product)
	} else {
		task.MarkFailed(result.ErrorMessage)
	}

	_ = s.publishRepo.Update(ctx, task)
	return task, nil
}
```

- [ ] **Step 10: 验证编译**

```bash
cd backend && go build ./...
```

Expected: 编译成功

- [ ] **Step 11: Commit**

```bash
git add backend/internal/domain/ backend/internal/model/anticorruption/
git commit -m "feat: add domain layer - aggregates, repositories, gateway interfaces, publish service"
```

---

## Task 4: 基础设施层 — Repository 实现 + Mock Gateway

**Files:**
- Create: `backend/internal/repository/source_item_repo.go`
- Create: `backend/internal/repository/product_repo.go`
- Create: `backend/internal/repository/publish_task_repo.go`
- Create: `backend/internal/gateway/mock_source_gateway.go`
- Create: `backend/internal/gateway/mock_target_gateway.go`
- Create: `backend/internal/queries/source_item_query.go`

- [ ] **Step 1: 创建 SourceItem Repository 实现 `backend/internal/repository/source_item_repo.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type SourceItemRepoImpl struct {
	db *gorm.DB
}

func NewSourceItemRepoImpl(db *gorm.DB) *SourceItemRepoImpl {
	return &SourceItemRepoImpl{db: db}
}

func (r *SourceItemRepoImpl) Save(ctx context.Context, item *entity.SourceItem) error {
	record := toSourceItemPO(item)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create source item: %w", err)
	}
	item.ID = record.ID
	return nil
}

func (r *SourceItemRepoImpl) FindByID(ctx context.Context, id int64) (*entity.SourceItem, error) {
	var record po.SourceItem
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find source item by id: %w", err)
	}
	return toSourceItemEntity(&record), nil
}

func (r *SourceItemRepoImpl) Update(ctx context.Context, item *entity.SourceItem) error {
	record := toSourceItemPO(item)
	record.ID = item.ID
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update source item: %w", err)
	}
	return nil
}

func toSourceItemPO(e *entity.SourceItem) *po.SourceItem {
	return &po.SourceItem{
		ID:          e.ID,
		Platform:    e.Platform,
		SourceURL:   e.SourceURL,
		ExternalID:  e.ExternalID,
		Title:       e.Title,
		Description: e.Description,
		Images:      po.StringSlice(e.Images),
		PriceMin:    e.Price.Min,
		PriceMax:    e.Price.Max,
		Supplier: po.SupplierInfo{
			Name:   e.Supplier.Name,
			Rating: e.Supplier.Rating,
			Region: e.Supplier.Region,
		},
		Category:    e.Category,
		Tags:        po.StringSlice(e.Tags),
		SalesVolume: e.SalesVolume,
		MinOrder:    e.MinOrder,
		Status:      e.Status,
		FetchedAt:   e.FetchedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toSourceItemEntity(p *po.SourceItem) *entity.SourceItem {
	return &entity.SourceItem{
		ID:          p.ID,
		Platform:    p.Platform,
		SourceURL:   p.SourceURL,
		ExternalID:  p.ExternalID,
		Title:       p.Title,
		Description: p.Description,
		Images:      []string(p.Images),
		Price:       entity.Price{Min: p.PriceMin, Max: p.PriceMax},
		Supplier: entity.Supplier{
			Name:   p.Supplier.Name,
			Rating: p.Supplier.Rating,
			Region: p.Supplier.Region,
		},
		Category:    p.Category,
		Tags:        []string(p.Tags),
		SalesVolume: p.SalesVolume,
		MinOrder:    p.MinOrder,
		Status:      p.Status,
		FetchedAt:   p.FetchedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
```

- [ ] **Step 2: 创建 Product Repository 实现 `backend/internal/repository/product_repo.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type ProductRepoImpl struct {
	db *gorm.DB
}

func NewProductRepoImpl(db *gorm.DB) *ProductRepoImpl {
	return &ProductRepoImpl{db: db}
}

func (r *ProductRepoImpl) Save(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := toProductPO(product)
		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("create product: %w", err)
		}
		product.ID = record.ID

		for i := range product.SKUs {
			sku := &po.ProductSKU{
				ProductID: product.ID,
				SpecName:  product.SKUs[i].SpecName,
				SpecValue: product.SKUs[i].SpecValue,
				Price:     product.SKUs[i].Price,
				Stock:     product.SKUs[i].Stock,
			}
			if err := tx.Create(sku).Error; err != nil {
				return fmt.Errorf("create product sku: %w", err)
			}
			product.SKUs[i].ID = sku.ID
		}
		return nil
	})
}

func (r *ProductRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	var record po.Product
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}

	var skuRecords []po.ProductSKU
	if err := r.db.WithContext(ctx).Where("product_id = ?", id).Find(&skuRecords).Error; err != nil {
		return nil, fmt.Errorf("find product skus: %w", err)
	}

	product := toProductEntity(&record)
	for _, s := range skuRecords {
		product.SKUs = append(product.SKUs, entity.SKU{
			ID:        s.ID,
			SpecName:  s.SpecName,
			SpecValue: s.SpecValue,
			Price:     s.Price,
			Stock:     s.Stock,
		})
	}
	return product, nil
}

func (r *ProductRepoImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := toProductPO(product)
		record.ID = product.ID
		if err := tx.Save(record).Error; err != nil {
			return fmt.Errorf("update product: %w", err)
		}

		if err := tx.Where("product_id = ?", product.ID).Delete(&po.ProductSKU{}).Error; err != nil {
			return fmt.Errorf("delete old skus: %w", err)
		}
		for i := range product.SKUs {
			sku := &po.ProductSKU{
				ProductID: product.ID,
				SpecName:  product.SKUs[i].SpecName,
				SpecValue: product.SKUs[i].SpecValue,
				Price:     product.SKUs[i].Price,
				Stock:     product.SKUs[i].Stock,
			}
			if err := tx.Create(sku).Error; err != nil {
				return fmt.Errorf("create product sku: %w", err)
			}
			product.SKUs[i].ID = sku.ID
		}
		return nil
	})
}

func toProductPO(e *entity.Product) *po.Product {
	return &po.Product{
		ID:           e.ID,
		SourceItemID: e.SourceItemID,
		Name:         e.Name,
		Description:  e.Description,
		Images:       po.StringSlice(e.Images),
		CostPrice:    e.CostPrice,
		SellPrice:    e.SellPrice,
		CategoryID:   e.CategoryID,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toProductEntity(p *po.Product) *entity.Product {
	return &entity.Product{
		ID:           p.ID,
		SourceItemID: p.SourceItemID,
		Name:         p.Name,
		Description:  p.Description,
		Images:       []string(p.Images),
		CostPrice:    p.CostPrice,
		SellPrice:    p.SellPrice,
		CategoryID:   p.CategoryID,
		Status:       p.Status,
		SKUs:         []entity.SKU{},
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
```

- [ ] **Step 3: 创建 PublishTask Repository 实现 `backend/internal/repository/publish_task_repo.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type PublishTaskRepoImpl struct {
	db *gorm.DB
}

func NewPublishTaskRepoImpl(db *gorm.DB) *PublishTaskRepoImpl {
	return &PublishTaskRepoImpl{db: db}
}

func (r *PublishTaskRepoImpl) Save(ctx context.Context, task *entity.PublishTask) error {
	record := toPublishTaskPO(task)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create publish task: %w", err)
	}
	task.ID = record.ID
	return nil
}

func (r *PublishTaskRepoImpl) FindByID(ctx context.Context, id int64) (*entity.PublishTask, error) {
	var record po.PublishTask
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find publish task by id: %w", err)
	}
	return toPublishTaskEntity(&record), nil
}

func (r *PublishTaskRepoImpl) Update(ctx context.Context, task *entity.PublishTask) error {
	record := toPublishTaskPO(task)
	record.ID = task.ID
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update publish task: %w", err)
	}
	return nil
}

func toPublishTaskPO(e *entity.PublishTask) *po.PublishTask {
	return &po.PublishTask{
		ID:                e.ID,
		ProductID:         e.ProductID,
		TargetPlatform:    e.TargetPlatform,
		PlatformProductID: e.PlatformProductID,
		PublishConfig: po.PublishConfig{
			CategoryID:      e.PublishConfig.CategoryID,
			FreightTemplate: e.PublishConfig.FreightTemplate,
		},
		Status:       e.Status,
		ErrorMessage: e.ErrorMessage,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toPublishTaskEntity(p *po.PublishTask) *entity.PublishTask {
	return &entity.PublishTask{
		ID:                p.ID,
		ProductID:         p.ProductID,
		TargetPlatform:    p.TargetPlatform,
		PlatformProductID: p.PlatformProductID,
		PublishConfig: entity.PublishConfig{
			CategoryID:      p.PublishConfig.CategoryID,
			FreightTemplate: p.PublishConfig.FreightTemplate,
		},
		Status:       p.Status,
		ErrorMessage: p.ErrorMessage,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
```

- [ ] **Step 4: 创建 Mock 货源网关 `backend/internal/gateway/mock_source_gateway.go`**

```go
package gateway

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type MockSourceGateway struct{}

func NewMockSourceGateway() *MockSourceGateway {
	return &MockSourceGateway{}
}

func (g *MockSourceGateway) FetchProduct(ctx context.Context, sourceURL string) (*anticorruption.SourceProduct, error) {
	return &anticorruption.SourceProduct{
		ExternalID:  "mock-ext-001",
		Title:       "Mock商品 - 高品质T恤",
		Description: "优质纯棉T恤，多色可选，批发价格优惠",
		Images:      []string{"https://via.placeholder.com/800x800?text=Mock+Image+1", "https://via.placeholder.com/800x800?text=Mock+Image+2"},
		PriceMin:    15.00,
		PriceMax:    25.00,
		Supplier: anticorruption.Supplier{
			Name:   "广州优品服饰有限公司",
			Rating: 4.8,
			Region: "广东广州",
		},
		Category:    "服装/T恤",
		SalesVolume: 10000,
		MinOrder:    2,
	}, nil
}
```

- [ ] **Step 5: 创建 Mock 目标平台网关 `backend/internal/gateway/mock_target_gateway.go`**

```go
package gateway

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type MockTargetGateway struct{}

func NewMockTargetGateway() *MockTargetGateway {
	return &MockTargetGateway{}
}

func (g *MockTargetGateway) PublishProduct(ctx context.Context, name, description string,
	images []string, sellPrice float64, config anticorruption.PublishConfig) (*anticorruption.PublishResult, error) {
	return &anticorruption.PublishResult{
		PlatformProductID: fmt.Sprintf("pdd-mock-%d", rand.Intn(1000000)),
		Success:           true,
		ErrorMessage:      "",
	}, nil
}
```

- [ ] **Step 6: 创建货源筛选查询 `backend/internal/queries/source_item_query.go`**

```go
package queries

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type SourceItemQuery struct {
	db *gorm.DB
}

func NewSourceItemQuery(db *gorm.DB) *SourceItemQuery {
	return &SourceItemQuery{db: db}
}

type ListResult struct {
	Items []po.SourceItem
	Total int64
}

func (q *SourceItemQuery) List(ctx context.Context, filter *dto.SourceItemFilter) (*ListResult, error) {
	query := q.db.WithContext(ctx).Model(&po.SourceItem{})

	if filter.Platform != nil {
		query = query.Where("platform = ?", *filter.Platform)
	}
	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.PriceMin != nil {
		query = query.Where("price_min >= ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		query = query.Where("price_max <= ?", *filter.PriceMax)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", *filter.Keyword))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count source items: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.SourceItem
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list source items: %w", err)
	}

	return &ListResult{Items: items, Total: total}, nil
}
```

- [ ] **Step 7: 验证编译**

```bash
cd backend && go build ./...
```

Expected: 编译成功

- [ ] **Step 8: Commit**

```bash
git add backend/internal/repository/ backend/internal/gateway/ backend/internal/queries/
git commit -m "feat: add infrastructure layer - repository impls, mock gateways, query service"
```

---

## Task 5: 应用层

**Files:**
- Create: `backend/internal/application/source_item_app.go`
- Create: `backend/internal/application/product_app.go`
- Create: `backend/internal/application/publish_app.go`

- [ ] **Step 1: 创建货源应用服务 `backend/internal/application/source_item_app.go`**

```go
package application

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/platform"
	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
)

type SourceItemApp struct {
	repo          repository.SourceItemRepository
	sourceGateway platform.SourcePlatformGateway
	query         *queries.SourceItemQuery
}

func NewSourceItemApp(
	repo repository.SourceItemRepository,
	gw platform.SourcePlatformGateway,
	query *queries.SourceItemQuery,
) *SourceItemApp {
	return &SourceItemApp{repo: repo, sourceGateway: gw, query: query}
}

func (a *SourceItemApp) Import(ctx context.Context, req *dto.ImportSourceItemReq) (*entity.SourceItem, error) {
	product, err := a.sourceGateway.FetchProduct(ctx, req.SourceURL)
	if err != nil {
		return nil, fmt.Errorf("fetch product from source: %w", err)
	}

	item := entity.NewSourceItem(
		req.Platform, req.SourceURL, product.ExternalID,
		product.Title, product.Description, product.Images,
		entity.Price{Min: product.PriceMin, Max: product.PriceMax},
		entity.Supplier{Name: product.Supplier.Name, Rating: product.Supplier.Rating, Region: product.Supplier.Region},
		product.Category, product.SalesVolume, product.MinOrder,
	)

	if err := a.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("save source item: %w", err)
	}
	return item, nil
}

func (a *SourceItemApp) GetByID(ctx context.Context, id int64) (*entity.SourceItem, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *SourceItemApp) List(ctx context.Context, filter *dto.SourceItemFilter) (*queries.ListResult, error) {
	return a.query.List(ctx, filter)
}

func (a *SourceItemApp) UpdateStatus(ctx context.Context, id int64, status string) error {
	item, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find source item: %w", err)
	}

	switch status {
	case "selected":
		if err := item.Select(); err != nil {
			return err
		}
	case "ignored":
		if err := item.Ignore(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	return a.repo.Update(ctx, item)
}

func (a *SourceItemApp) AddTag(ctx context.Context, id int64, tag string) error {
	item, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find source item: %w", err)
	}
	item.AddTag(tag)
	return a.repo.Update(ctx, item)
}
```

- [ ] **Step 2: 创建商品应用服务 `backend/internal/application/product_app.go`**

```go
package application

import (
	"context"
	"fmt"

	productEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	sourceRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type ProductApp struct {
	productRepo productRepo.ProductRepository
	sourceRepo  sourceRepo.SourceItemRepository
	db          *gorm.DB
}

func NewProductApp(pr productRepo.ProductRepository, sr sourceRepo.SourceItemRepository, db *gorm.DB) *ProductApp {
	return &ProductApp{productRepo: pr, sourceRepo: sr, db: db}
}

func (a *ProductApp) CreateFromSource(ctx context.Context, sourceItemID int64) (*productEntity.Product, error) {
	source, err := a.sourceRepo.FindByID(ctx, sourceItemID)
	if err != nil {
		return nil, fmt.Errorf("find source item: %w", err)
	}

	product := productEntity.CreateFromSource(source)
	if err := a.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}
	return product, nil
}

func (a *ProductApp) GetByID(ctx context.Context, id int64) (*productEntity.Product, error) {
	return a.productRepo.FindByID(ctx, id)
}

func (a *ProductApp) List(ctx context.Context, filter *dto.ProductFilter) ([]po.Product, int64, error) {
	query := a.db.WithContext(ctx).Model(&po.Product{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *filter.Keyword))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.Product
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}

	return items, total, nil
}

func (a *ProductApp) Update(ctx context.Context, id int64, req *dto.UpdateProductReq) (*productEntity.Product, error) {
	product, err := a.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	product.EditInfo(req.Name, req.Description, req.Images, req.CategoryID)

	if req.CostPrice != nil || req.SellPrice != nil {
		product.SetPrice(req.CostPrice, req.SellPrice)
	}

	if req.SKUs != nil {
		skus := make([]productEntity.SKU, len(req.SKUs))
		for i, s := range req.SKUs {
			skus[i] = productEntity.SKU{
				ID:        s.ID,
				SpecName:  s.SpecName,
				SpecValue: s.SpecValue,
				Price:     s.Price,
				Stock:     s.Stock,
			}
		}
		product.SetSKUs(skus)
	}

	if err := a.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	return product, nil
}

func (a *ProductApp) MarkReady(ctx context.Context, id int64) error {
	product, err := a.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}
	if err := product.MarkReady(); err != nil {
		return err
	}
	return a.productRepo.Update(ctx, product)
}
```

- [ ] **Step 3: 创建发品应用服务 `backend/internal/application/publish_app.go`**

```go
package application

import (
	"context"
	"fmt"

	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	publishEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type PublishApp struct {
	productRepo    productRepo.ProductRepository
	publishService *domainservice.PublishDomainService
	db             *gorm.DB
}

func NewPublishApp(
	pr productRepo.ProductRepository,
	ps *domainservice.PublishDomainService,
	db *gorm.DB,
) *PublishApp {
	return &PublishApp{productRepo: pr, publishService: ps, db: db}
}

func (a *PublishApp) CreateTask(ctx context.Context, req *dto.CreatePublishTaskReq) (*publishEntity.PublishTask, error) {
	product, err := a.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	config := publishEntity.PublishConfig{
		CategoryID:      req.CategoryID,
		FreightTemplate: req.FreightTemplate,
	}

	task, err := a.publishService.PublishProduct(ctx, product, req.TargetPlatform, config)
	if err != nil {
		return nil, fmt.Errorf("publish product: %w", err)
	}
	return task, nil
}

func (a *PublishApp) List(ctx context.Context, filter *dto.PublishTaskFilter) ([]po.PublishTask, int64, error) {
	query := a.db.WithContext(ctx).Model(&po.PublishTask{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.TargetPlatform != nil {
		query = query.Where("target_platform = ?", *filter.TargetPlatform)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count publish tasks: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.PublishTask
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list publish tasks: %w", err)
	}

	return items, total, nil
}
```

- [ ] **Step 4: 验证编译**

```bash
cd backend && go build ./...
```

Expected: 编译成功

- [ ] **Step 5: Commit**

```bash
git add backend/internal/application/
git commit -m "feat: add application layer - source item, product, publish use cases"
```

---

## Task 6: 接口层 — HTTP Handlers + 路由 + Wire 依赖注入

**Files:**
- Create: `backend/internal/server/routes.go`
- Create: `backend/internal/server/source_item_handler.go`
- Create: `backend/internal/server/product_handler.go`
- Create: `backend/internal/server/publish_handler.go`
- Create: `backend/internal/wire.go`
- Modify: `backend/main.go`

- [ ] **Step 1: 创建货源 Handler `backend/internal/server/source_item_handler.go`**

```go
package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type SourceItemHandler struct {
	app *application.SourceItemApp
}

func NewSourceItemHandler(app *application.SourceItemApp) *SourceItemHandler {
	return &SourceItemHandler{app: app}
}

func (h *SourceItemHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req dto.ImportSourceItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	item, err := h.app.Import(r.Context(), &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *SourceItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	item, err := h.app.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *SourceItemHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.SourceItemFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("platform"); v != "" {
		filter.Platform = &v
	}
	if v := r.URL.Query().Get("category"); v != "" {
		filter.Category = &v
	}
	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("keyword"); v != "" {
		filter.Keyword = &v
	}

	result, err := h.app.List(r.Context(), &filter)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"items": result.Items,
		"total": result.Total,
	})
}

func (h *SourceItemHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.UpdateSourceItemStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	if err := h.app.UpdateStatus(r.Context(), id, req.Status); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}

func (h *SourceItemHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.AddTagReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	if err := h.app.AddTag(r.Context(), id, req.Tag); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}
```

- [ ] **Step 2: 创建商品 Handler `backend/internal/server/product_handler.go`**

```go
package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type ProductHandler struct {
	app *application.ProductApp
}

func NewProductHandler(app *application.ProductApp) *ProductHandler {
	return &ProductHandler{app: app}
}

func (h *ProductHandler) CreateFromSource(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductFromSourceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	product, err := h.app.CreateFromSource(r.Context(), req.SourceItemID)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	product, err := h.app.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.ProductFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("keyword"); v != "" {
		filter.Keyword = &v
	}

	items, total, err := h.app.List(r.Context(), &filter)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"items": items,
		"total": total,
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.UpdateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	product, err := h.app.Update(r.Context(), id, &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) MarkReady(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	if err := h.app.MarkReady(r.Context(), id); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}
```

- [ ] **Step 3: 创建发品 Handler `backend/internal/server/publish_handler.go`**

```go
package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type PublishHandler struct {
	app *application.PublishApp
}

func NewPublishHandler(app *application.PublishApp) *PublishHandler {
	return &PublishHandler{app: app}
}

func (h *PublishHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePublishTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	task, err := h.app.CreateTask(r.Context(), &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, task)
}

func (h *PublishHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.PublishTaskFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("targetPlatform"); v != "" {
		filter.TargetPlatform = &v
	}

	items, total, err := h.app.List(r.Context(), &filter)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"items": items,
		"total": total,
	})
}
```

- [ ] **Step 4: 创建路由 `backend/internal/server/routes.go`**

```go
package server

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(engine *rest.Server, si *SourceItemHandler, p *ProductHandler, pub *PublishHandler) {
	engine.AddRoutes([]rest.Route{
		{Method: http.MethodPost, Path: "/api/source-items/import", Handler: si.Import},
		{Method: http.MethodGet, Path: "/api/source-items", Handler: si.List},
		{Method: http.MethodGet, Path: "/api/source-items/detail", Handler: si.GetByID},
		{Method: http.MethodPut, Path: "/api/source-items/status", Handler: si.UpdateStatus},
		{Method: http.MethodPost, Path: "/api/source-items/tag", Handler: si.AddTag},

		{Method: http.MethodPost, Path: "/api/products/create-from-source", Handler: p.CreateFromSource},
		{Method: http.MethodGet, Path: "/api/products", Handler: p.List},
		{Method: http.MethodGet, Path: "/api/products/detail", Handler: p.GetByID},
		{Method: http.MethodPut, Path: "/api/products", Handler: p.Update},
		{Method: http.MethodPut, Path: "/api/products/ready", Handler: p.MarkReady},

		{Method: http.MethodPost, Path: "/api/publish-tasks", Handler: pub.CreateTask},
		{Method: http.MethodGet, Path: "/api/publish-tasks", Handler: pub.List},
	})
}
```

- [ ] **Step 5: 创建 Wire 依赖注入 `backend/internal/wire.go`**

由于 Google Wire 需要代码生成，为了简化先手写工厂函数 `backend/internal/wire.go`：

```go
package internal

import (
	"github.com/yangboyi/ddd-dev/backend/internal/application"
	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	"github.com/yangboyi/ddd-dev/backend/internal/gateway"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
	repo "github.com/yangboyi/ddd-dev/backend/internal/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"gorm.io/gorm"
)

type Handlers struct {
	SourceItem *server.SourceItemHandler
	Product    *server.ProductHandler
	Publish    *server.PublishHandler
}

func InitHandlers(db *gorm.DB) *Handlers {
	// Repositories
	sourceItemRepo := repo.NewSourceItemRepoImpl(db)
	productRepo := repo.NewProductRepoImpl(db)
	publishTaskRepo := repo.NewPublishTaskRepoImpl(db)

	// Gateways (Mock)
	sourceGateway := gateway.NewMockSourceGateway()
	targetGateway := gateway.NewMockTargetGateway()

	// Queries
	sourceItemQuery := queries.NewSourceItemQuery(db)

	// Domain Services
	publishDomainService := domainservice.NewPublishDomainService(productRepo, publishTaskRepo, targetGateway)

	// Application Services
	sourceItemApp := application.NewSourceItemApp(sourceItemRepo, sourceGateway, sourceItemQuery)
	productApp := application.NewProductApp(productRepo, sourceItemRepo, db)
	publishApp := application.NewPublishApp(productRepo, publishDomainService, db)

	// Handlers
	return &Handlers{
		SourceItem: server.NewSourceItemHandler(sourceItemApp),
		Product:    server.NewProductHandler(productApp),
		Publish:    server.NewPublishHandler(publishApp),
	}
}
```

- [ ] **Step 6: 更新 `backend/main.go` 注册路由**

替换整个 main.go：

```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/backend/infra/config"
	"github.com/yangboyi/ddd-dev/backend/infra/vars"
	"github.com/yangboyi/ddd-dev/backend/internal"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}
	vars.DB = db

	// Auto migrate tables
	if err := db.AutoMigrate(&po.SourceItem{}, &po.Product{}, &po.ProductSKU{}, &po.PublishTask{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	srv := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer srv.Stop()

	handlers := internal.InitHandlers(db)
	server.RegisterRoutes(srv, handlers.SourceItem, handlers.Product, handlers.Publish)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	srv.Start()
}
```

- [ ] **Step 7: 验证编译**

```bash
cd backend && go build ./...
```

Expected: 编译成功

- [ ] **Step 8: Commit**

```bash
git add backend/internal/server/ backend/internal/wire.go backend/main.go
git commit -m "feat: add server layer - handlers, routes, dependency injection, auto-migrate"
```

---

## Task 7: 后端端到端验证

- [ ] **Step 1: 确保 MySQL 可用，创建数据库**

```bash
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS dropship CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

- [ ] **Step 2: 启动后端服务**

```bash
cd backend && go run main.go -f etc/config.yaml
```

Expected: 输出 `Starting server at 0.0.0.0:8888...`

- [ ] **Step 3: 测试导入货源**

```bash
curl -X POST http://localhost:8888/api/source-items/import \
  -H "Content-Type: application/json" \
  -d '{"platform":"ali1688","sourceUrl":"https://detail.1688.com/offer/123456.html"}'
```

Expected: 返回 `{"code":0,"message":"ok","data":{...}}` 包含 mock 商品信息

- [ ] **Step 4: 测试货源列表**

```bash
curl "http://localhost:8888/api/source-items?page=1&pageSize=10"
```

Expected: 返回列表，包含刚导入的货源

- [ ] **Step 5: 测试从货源创建商品**

```bash
curl -X POST http://localhost:8888/api/products/create-from-source \
  -H "Content-Type: application/json" \
  -d '{"sourceItemId":1}'
```

Expected: 返回新创建的商品，status=draft

- [ ] **Step 6: 测试更新商品 + 标记就绪**

```bash
curl -X PUT "http://localhost:8888/api/products?id=1" \
  -H "Content-Type: application/json" \
  -d '{"sellPrice":39.9,"name":"高品质纯棉T恤 多色可选"}'

curl -X PUT "http://localhost:8888/api/products/ready?id=1"
```

Expected: 两次都返回成功

- [ ] **Step 7: 测试发品**

```bash
curl -X POST http://localhost:8888/api/publish-tasks \
  -H "Content-Type: application/json" \
  -d '{"productId":1,"targetPlatform":"pdd","categoryId":"cat-001","freightTemplate":"tpl-001"}'
```

Expected: 返回 publish task，status=success（Mock 网关直接成功）

- [ ] **Step 8: Commit（如有修复）**

```bash
git add -A && git commit -m "fix: resolve issues found during e2e testing"
```

---

## Task 8: 前端项目初始化

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/next.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tailwind.config.ts`
- Create: `frontend/postcss.config.js`
- Create: `frontend/src/app/layout.tsx`
- Create: `frontend/src/app/page.tsx`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/components/layout/sidebar.tsx`

- [ ] **Step 1: 创建 Next.js 项目**

```bash
cd /Users/yangboyi/github/ddd-dev
npx create-next-app@latest frontend --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --no-turbopack
```

- [ ] **Step 2: 安装 shadcn/ui**

```bash
cd frontend
npx shadcn@latest init -d
```

- [ ] **Step 3: 安装依赖组件**

```bash
cd frontend
npx shadcn@latest add button card dialog input label select table badge
npm install swr @tanstack/react-table
```

- [ ] **Step 4: 创建 API 请求封装 `frontend/src/lib/api.ts`**

```typescript
const API_BASE = "http://localhost:8888/api";

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  const json: ApiResponse<T> = await res.json();
  if (json.code !== 0) {
    throw new Error(json.message);
  }
  return json.data;
}

export const api = {
  sourceItems: {
    import: (data: { platform: string; sourceUrl: string }) =>
      request("/source-items/import", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/source-items?${new URLSearchParams(params)}`
      ),
    updateStatus: (id: number, status: string) =>
      request(`/source-items/status?id=${id}`, {
        method: "PUT",
        body: JSON.stringify({ status }),
      }),
    addTag: (id: number, tag: string) =>
      request(`/source-items/tag?id=${id}`, {
        method: "POST",
        body: JSON.stringify({ tag }),
      }),
  },
  products: {
    createFromSource: (sourceItemId: number) =>
      request("/products/create-from-source", {
        method: "POST",
        body: JSON.stringify({ sourceItemId }),
      }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/products?${new URLSearchParams(params)}`
      ),
    get: (id: number) => request(`/products/detail?id=${id}`),
    update: (id: number, data: any) =>
      request(`/products?id=${id}`, { method: "PUT", body: JSON.stringify(data) }),
    markReady: (id: number) =>
      request(`/products/ready?id=${id}`, { method: "PUT" }),
  },
  publishTasks: {
    create: (data: {
      productId: number;
      targetPlatform: string;
      categoryId: string;
      freightTemplate: string;
    }) => request("/publish-tasks", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(
        `/publish-tasks?${new URLSearchParams(params)}`
      ),
  },
};
```

- [ ] **Step 5: 创建侧边栏 `frontend/src/components/layout/sidebar.tsx`**

```tsx
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/sources", label: "货源管理", icon: "📦" },
  { href: "/products", label: "商品管理", icon: "🏷️" },
  { href: "/publish", label: "发品任务", icon: "🚀" },
];

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-56 border-r bg-gray-50 p-4 min-h-screen">
      <h1 className="text-lg font-bold mb-6 px-2">代发工具</h1>
      <nav className="space-y-1">
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-center gap-2 px-3 py-2 rounded-md text-sm transition-colors",
              pathname.startsWith(item.href)
                ? "bg-white shadow-sm font-medium"
                : "text-gray-600 hover:bg-white/60"
            )}
          >
            <span>{item.icon}</span>
            {item.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
```

- [ ] **Step 6: 更新根布局 `frontend/src/app/layout.tsx`**

```tsx
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Sidebar } from "@/components/layout/sidebar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "代发工具",
  description: "电商一键代发运营工具",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        <div className="flex">
          <Sidebar />
          <main className="flex-1 p-6">{children}</main>
        </div>
      </body>
    </html>
  );
}
```

- [ ] **Step 7: 更新首页重定向 `frontend/src/app/page.tsx`**

```tsx
import { redirect } from "next/navigation";

export default function Home() {
  redirect("/sources");
}
```

- [ ] **Step 8: 验证前端启动**

```bash
cd frontend && npm run dev
```

Expected: 打开 http://localhost:3000 看到侧边栏和空白页面

- [ ] **Step 9: Commit**

```bash
git add frontend/
git commit -m "feat: init frontend with Next.js, shadcn/ui, sidebar layout, API client"
```

---

## Task 9: 前端 — 货源管理页

**Files:**
- Create: `frontend/src/app/sources/page.tsx`
- Create: `frontend/src/components/sources/source-table.tsx`
- Create: `frontend/src/components/sources/import-dialog.tsx`

- [ ] **Step 1: 创建导入对话框 `frontend/src/components/sources/import-dialog.tsx`**

```tsx
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";

export function ImportDialog({ onSuccess }: { onSuccess: () => void }) {
  const [open, setOpen] = useState(false);
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);

  const handleImport = async () => {
    setLoading(true);
    try {
      await api.sourceItems.import({ platform: "ali1688", sourceUrl: url });
      setUrl("");
      setOpen(false);
      onSuccess();
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>导入货源</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>导入货源</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div>
            <Label>货源链接</Label>
            <Input
              placeholder="粘贴1688商品链接"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
            />
          </div>
          <Button onClick={handleImport} disabled={!url || loading} className="w-full">
            {loading ? "导入中..." : "导入"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 2: 创建货源列表表格 `frontend/src/components/sources/source-table.tsx`**

```tsx
"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";

interface SourceItem {
  ID: number;
  Platform: string;
  Title: string;
  PriceMin: number;
  PriceMax: number;
  Supplier: { Name: string; Rating: number; Region: string };
  Category: string;
  SalesVolume: number;
  Status: string;
}

const statusMap: Record<string, { label: string; variant: "default" | "secondary" | "destructive" }> = {
  new: { label: "新导入", variant: "secondary" },
  selected: { label: "已选品", variant: "default" },
  ignored: { label: "已忽略", variant: "destructive" },
};

export function SourceTable({
  items,
  onRefresh,
}: {
  items: SourceItem[];
  onRefresh: () => void;
}) {
  const handleSelect = async (id: number) => {
    await api.sourceItems.updateStatus(id, "selected");
    onRefresh();
  };

  const handleCreateProduct = async (id: number) => {
    await api.products.createFromSource(id);
    alert("商品创建成功，请到商品管理页编辑");
    onRefresh();
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>标题</TableHead>
          <TableHead>平台</TableHead>
          <TableHead>价格区间</TableHead>
          <TableHead>供应商</TableHead>
          <TableHead>销量</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>操作</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => {
          const status = statusMap[item.Status] || statusMap.new;
          return (
            <TableRow key={item.ID}>
              <TableCell className="max-w-[200px] truncate">{item.Title}</TableCell>
              <TableCell>{item.Platform}</TableCell>
              <TableCell>
                ¥{item.PriceMin} - ¥{item.PriceMax}
              </TableCell>
              <TableCell>
                {item.Supplier.Name} ({item.Supplier.Rating})
              </TableCell>
              <TableCell>{item.SalesVolume}</TableCell>
              <TableCell>
                <Badge variant={status.variant}>{status.label}</Badge>
              </TableCell>
              <TableCell className="space-x-2">
                {item.Status === "new" && (
                  <Button size="sm" variant="outline" onClick={() => handleSelect(item.ID)}>
                    选品
                  </Button>
                )}
                {item.Status === "selected" && (
                  <Button size="sm" onClick={() => handleCreateProduct(item.ID)}>
                    创建商品
                  </Button>
                )}
              </TableCell>
            </TableRow>
          );
        })}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={7} className="text-center text-gray-400 py-8">
              暂无货源，点击「导入货源」开始
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
```

- [ ] **Step 3: 创建货源管理页 `frontend/src/app/sources/page.tsx`**

```tsx
"use client";

import { useCallback, useEffect, useState } from "react";
import { ImportDialog } from "@/components/sources/import-dialog";
import { SourceTable } from "@/components/sources/source-table";
import { Input } from "@/components/ui/input";
import { api } from "@/lib/api";

export default function SourcesPage() {
  const [items, setItems] = useState<any[]>([]);
  const [keyword, setKeyword] = useState("");

  const fetchData = useCallback(async () => {
    const params: Record<string, string> = { page: "1", pageSize: "50" };
    if (keyword) params.keyword = keyword;
    const res = await api.sourceItems.list(params);
    setItems(res.items || []);
  }, [keyword]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">货源管理</h2>
        <ImportDialog onSuccess={fetchData} />
      </div>
      <div className="mb-4">
        <Input
          placeholder="搜索货源标题..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          className="max-w-sm"
        />
      </div>
      <SourceTable items={items} onRefresh={fetchData} />
    </div>
  );
}
```

- [ ] **Step 4: 验证货源管理页**

```bash
cd frontend && npm run dev
```

打开 http://localhost:3000/sources ，确认页面渲染正常

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/sources/ frontend/src/components/sources/
git commit -m "feat: add sources management page with import and list"
```

---

## Task 10: 前端 — 商品管理页

**Files:**
- Create: `frontend/src/app/products/page.tsx`
- Create: `frontend/src/components/products/product-table.tsx`
- Create: `frontend/src/components/products/product-edit-dialog.tsx`

- [ ] **Step 1: 创建商品编辑对话框 `frontend/src/components/products/product-edit-dialog.tsx`**

```tsx
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";

interface Product {
  ID: number;
  Name: string;
  Description: string;
  CostPrice: number;
  SellPrice: number;
}

export function ProductEditDialog({
  product,
  open,
  onOpenChange,
  onSuccess,
}: {
  product: Product | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}) {
  const [name, setName] = useState(product?.Name || "");
  const [sellPrice, setSellPrice] = useState(String(product?.SellPrice || ""));
  const [loading, setLoading] = useState(false);

  const handleSave = async () => {
    if (!product) return;
    setLoading(true);
    try {
      await api.products.update(product.ID, {
        name,
        sellPrice: parseFloat(sellPrice),
      });
      onOpenChange(false);
      onSuccess();
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>编辑商品</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div>
            <Label>商品名称</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div>
            <Label>成本价</Label>
            <Input value={product?.CostPrice || 0} disabled />
          </div>
          <div>
            <Label>售价</Label>
            <Input
              type="number"
              value={sellPrice}
              onChange={(e) => setSellPrice(e.target.value)}
            />
          </div>
          <Button onClick={handleSave} disabled={loading} className="w-full">
            {loading ? "保存中..." : "保存"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
```

- [ ] **Step 2: 创建商品列表表格 `frontend/src/components/products/product-table.tsx`**

```tsx
"use client";

import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ProductEditDialog } from "./product-edit-dialog";
import { api } from "@/lib/api";

interface Product {
  ID: number;
  Name: string;
  Description: string;
  CostPrice: number;
  SellPrice: number;
  Status: string;
  CategoryID: string;
}

const statusMap: Record<string, { label: string; variant: "default" | "secondary" | "outline" }> = {
  draft: { label: "草稿", variant: "secondary" },
  ready: { label: "就绪", variant: "default" },
  published: { label: "已发布", variant: "outline" },
};

export function ProductTable({
  items,
  onRefresh,
}: {
  items: Product[];
  onRefresh: () => void;
}) {
  const [editProduct, setEditProduct] = useState<Product | null>(null);

  const handleMarkReady = async (id: number) => {
    try {
      await api.products.markReady(id);
      onRefresh();
    } catch (e: any) {
      alert(e.message);
    }
  };

  const handlePublish = async (id: number) => {
    try {
      await api.publishTasks.create({
        productId: id,
        targetPlatform: "pdd",
        categoryId: "cat-001",
        freightTemplate: "tpl-001",
      });
      alert("发品任务已创建，请到发品任务页查看");
      onRefresh();
    } catch (e: any) {
      alert(e.message);
    }
  };

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>商品名称</TableHead>
            <TableHead>成本价</TableHead>
            <TableHead>售价</TableHead>
            <TableHead>利润</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>操作</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((item) => {
            const status = statusMap[item.Status] || statusMap.draft;
            const profit = item.SellPrice > 0 ? (item.SellPrice - item.CostPrice).toFixed(2) : "-";
            return (
              <TableRow key={item.ID}>
                <TableCell className="max-w-[250px] truncate">{item.Name}</TableCell>
                <TableCell>¥{item.CostPrice}</TableCell>
                <TableCell>{item.SellPrice > 0 ? `¥${item.SellPrice}` : "-"}</TableCell>
                <TableCell>{profit !== "-" ? `¥${profit}` : "-"}</TableCell>
                <TableCell>
                  <Badge variant={status.variant}>{status.label}</Badge>
                </TableCell>
                <TableCell className="space-x-2">
                  <Button size="sm" variant="outline" onClick={() => setEditProduct(item)}>
                    编辑
                  </Button>
                  {item.Status === "draft" && (
                    <Button size="sm" variant="outline" onClick={() => handleMarkReady(item.ID)}>
                      标记就绪
                    </Button>
                  )}
                  {item.Status === "ready" && (
                    <Button size="sm" onClick={() => handlePublish(item.ID)}>
                      发布到PDD
                    </Button>
                  )}
                </TableCell>
              </TableRow>
            );
          })}
          {items.length === 0 && (
            <TableRow>
              <TableCell colSpan={6} className="text-center text-gray-400 py-8">
                暂无商品，请先从货源管理页选品
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
      <ProductEditDialog
        product={editProduct}
        open={!!editProduct}
        onOpenChange={(open) => !open && setEditProduct(null)}
        onSuccess={onRefresh}
      />
    </>
  );
}
```

- [ ] **Step 3: 创建商品管理页 `frontend/src/app/products/page.tsx`**

```tsx
"use client";

import { useCallback, useEffect, useState } from "react";
import { ProductTable } from "@/components/products/product-table";
import { api } from "@/lib/api";

export default function ProductsPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.products.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">商品管理</h2>
      <ProductTable items={items} onRefresh={fetchData} />
    </div>
  );
}
```

- [ ] **Step 4: 验证商品管理页**

打开 http://localhost:3000/products ，确认页面渲染正常

- [ ] **Step 5: Commit**

```bash
git add frontend/src/app/products/ frontend/src/components/products/
git commit -m "feat: add products management page with edit and publish"
```

---

## Task 11: 前端 — 发品任务页

**Files:**
- Create: `frontend/src/app/publish/page.tsx`
- Create: `frontend/src/components/publish/publish-table.tsx`

- [ ] **Step 1: 创建发品任务表格 `frontend/src/components/publish/publish-table.tsx`**

```tsx
"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

interface PublishTask {
  ID: number;
  ProductID: number;
  TargetPlatform: string;
  PlatformProductID: string;
  Status: string;
  ErrorMessage: string;
  CreatedAt: string;
}

const statusMap: Record<
  string,
  { label: string; variant: "default" | "secondary" | "destructive" | "outline" }
> = {
  pending: { label: "待发布", variant: "secondary" },
  publishing: { label: "发布中", variant: "outline" },
  success: { label: "成功", variant: "default" },
  failed: { label: "失败", variant: "destructive" },
};

export function PublishTable({ items }: { items: PublishTask[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>任务ID</TableHead>
          <TableHead>商品ID</TableHead>
          <TableHead>目标平台</TableHead>
          <TableHead>平台商品ID</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>错误信息</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => {
          const status = statusMap[item.Status] || statusMap.pending;
          return (
            <TableRow key={item.ID}>
              <TableCell>{item.ID}</TableCell>
              <TableCell>{item.ProductID}</TableCell>
              <TableCell>{item.TargetPlatform}</TableCell>
              <TableCell>{item.PlatformProductID || "-"}</TableCell>
              <TableCell>
                <Badge variant={status.variant}>{status.label}</Badge>
              </TableCell>
              <TableCell className="max-w-[200px] truncate text-red-500">
                {item.ErrorMessage || "-"}
              </TableCell>
            </TableRow>
          );
        })}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={6} className="text-center text-gray-400 py-8">
              暂无发品任务
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
```

- [ ] **Step 2: 创建发品任务页 `frontend/src/app/publish/page.tsx`**

```tsx
"use client";

import { useCallback, useEffect, useState } from "react";
import { PublishTable } from "@/components/publish/publish-table";
import { api } from "@/lib/api";

export default function PublishPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.publishTasks.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">发品任务</h2>
      <PublishTable items={items} />
    </div>
  );
}
```

- [ ] **Step 3: 验证发品任务页**

打开 http://localhost:3000/publish ，确认页面渲染正常

- [ ] **Step 4: Commit**

```bash
git add frontend/src/app/publish/ frontend/src/components/publish/
git commit -m "feat: add publish tasks page"
```

---

## Task 12: 端到端全链路验证

- [ ] **Step 1: 确保后端运行**

```bash
cd backend && go run main.go -f etc/config.yaml
```

- [ ] **Step 2: 确保前端运行**

```bash
cd frontend && npm run dev
```

- [ ] **Step 3: 全链路操作验证**

1. 打开 http://localhost:3000/sources
2. 点击「导入货源」→ 粘贴任意 URL → 点击导入
3. 确认货源列表显示 mock 数据
4. 点击「选品」→ 点击「创建商品」
5. 切换到「商品管理」页
6. 点击「编辑」→ 设置售价 39.9 → 保存
7. 点击「标记就绪」
8. 点击「发布到PDD」
9. 切换到「发品任务」页
10. 确认任务状态为「成功」

Expected: 全链路跑通，每个步骤正常响应

- [ ] **Step 4: 修复发现的问题（如有）**

```bash
git add -A && git commit -m "fix: resolve issues found during full e2e testing"
```

- [ ] **Step 5: 最终 Commit**

```bash
git add -A && git commit -m "feat: complete core pipeline - source import, product management, publish to platform"
```
