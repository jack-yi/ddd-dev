# 用户中心微服务实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新增用户中心微服务，提供 Google OAuth 2.0 登录、JWT 鉴权、RBAC 角色权限，与 dropship-api 通过 gRPC 通信

**Architecture:** user-center 为独立 Go module（同一 repo），包含 REST API (:8880) 和 zRPC (:8881) 两个入口。共享 domain/model/repository 层。dropship-api 通过 etcd 发现 user-center-rpc 并调用鉴权接口。前端新增登录页、用户管理页、超管初始化页。

**Tech Stack:** Go-zero (REST + zRPC), GORM, JWT (golang-jwt/jwt/v5), Google OAuth 2.0, Protobuf, etcd | Next.js 14, shadcn/ui

**Spec:** `docs/superpowers/specs/2026-03-31-user-center-design.md`

---

## File Structure

### user-center/ (新建)

```
user-center/
├── go.mod
├── cmd/
│   ├── api/main.go                    # REST API 入口 :8880
│   └── rpc/main.go                    # zRPC 入口 :8881
├── etc/
│   ├── api.yaml                       # API 配置
│   └── rpc.yaml                       # RPC 配置
├── proto/
│   └── usercenter.proto               # Protobuf 定义
├── internal/
│   ├── config/
│   │   ├── api_config.go              # API 配置结构体
│   │   └── rpc_config.go              # RPC 配置结构体
│   ├── domain/
│   │   └── user/
│   │       ├── entity/user.go         # User/Role 实体 + 值对象
│   │       └── repository/repository.go
│   ├── model/
│   │   ├── po/mysql/user.go           # User PO
│   │   ├── po/mysql/role.go           # Role PO
│   │   ├── po/mysql/user_role.go      # UserRole PO
│   │   └── dto/user_dto.go            # DTO
│   ├── repository/
│   │   ├── user_repo.go               # User GORM 实现
│   │   └── role_repo.go               # Role GORM 实现
│   ├── application/
│   │   ├── user_app.go                # 用户业务逻辑
│   │   └── jwt.go                     # JWT 工具
│   ├── seed/
│   │   └── seed.go                    # 预设角色种子数据
│   ├── server/
│   │   ├── api/
│   │   │   ├── routes.go              # API 路由
│   │   │   ├── auth_handler.go        # OAuth handlers
│   │   │   ├── user_handler.go        # 用户管理 handlers
│   │   │   └── init_handler.go        # 超管初始化
│   │   └── rpc/
│   │       └── usercenter_server.go   # zRPC 实现
│   ├── middleware/
│   │   └── auth.go                    # JWT 中间件
│   └── wire.go                        # 依赖注入
```

### backend/ (修改)

```
backend/
├── infra/config/config.go             # +UserCenterRpc 配置
├── etc/config.yaml                    # +RPC client 配置
├── internal/
│   ├── middleware/
│   │   └── auth.go                    # 新增: 鉴权中间件
│   ├── server/routes.go               # 修改: 应用中间件
│   └── wire.go                        # 修改: 注入RPC client
└── main.go                            # 修改: 初始化RPC client
```

### frontend/ (修改+新增)

```
frontend/src/
├── lib/
│   ├── api.ts                         # 修改: 增加auth header + user-center API
│   └── auth.ts                        # 新增: token 管理
├── app/
│   ├── layout.tsx                     # 修改: auth guard
│   ├── login/page.tsx                 # 新增: 登录页
│   ├── init/page.tsx                  # 新增: 超管初始化
│   └── users/page.tsx                 # 新增: 用户管理
├── components/
│   ├── layout/sidebar.tsx             # 修改: 角色可见性
│   └── users/user-table.tsx           # 新增: 用户表格
```

---

## Task 1: user-center 项目初始化 + Protobuf

**Files:**
- Create: `user-center/go.mod`
- Create: `user-center/proto/usercenter.proto`
- Create: `user-center/etc/api.yaml`
- Create: `user-center/etc/rpc.yaml`

- [ ] **Step 1: 创建 go.mod**

```bash
mkdir -p /Users/yangboyi/github/ddd-dev/user-center
cd /Users/yangboyi/github/ddd-dev/user-center
go mod init github.com/yangboyi/ddd-dev/user-center
```

- [ ] **Step 2: 创建 Protobuf `user-center/proto/usercenter.proto`**

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

- [ ] **Step 3: 生成 Go 代码**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center
# 安装 protoc-gen-go 和 protoc-gen-go-grpc（如未安装）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

mkdir -p proto/pb
protoc --go_out=proto/pb --go_opt=paths=source_relative \
       --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
       proto/usercenter.proto
```

- [ ] **Step 4: 创建配置文件 `user-center/etc/api.yaml`**

```yaml
Name: user-center-api
Host: 0.0.0.0
Port: 8880

MySQL:
  DataSource: "root:root123@tcp(127.0.0.1:3306)/dropship?charset=utf8mb4&parseTime=True&loc=Local"

JWT:
  Secret: "your-jwt-secret-key-change-in-production"
  Expire: 86400

Google:
  ClientID: "your-google-client-id"
  ClientSecret: "your-google-client-secret"
  RedirectURL: "http://localhost:8880/api/auth/google/callback"

UserCenterRpc:
  Etcd:
    Hosts:
      - "127.0.0.1:2379"
    Key: user-center.rpc
```

- [ ] **Step 5: 创建配置文件 `user-center/etc/rpc.yaml`**

```yaml
Name: user-center.rpc
ListenOn: 0.0.0.0:8881

Etcd:
  Hosts:
    - "127.0.0.1:2379"
  Key: user-center.rpc

MySQL:
  DataSource: "root:root123@tcp(127.0.0.1:3306)/dropship?charset=utf8mb4&parseTime=True&loc=Local"

JWT:
  Secret: "your-jwt-secret-key-change-in-production"
  Expire: 86400
```

- [ ] **Step 6: 安装依赖**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center
go get github.com/zeromicro/go-zero@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/mysql@latest
go get github.com/golang-jwt/jwt/v5@latest
go get golang.org/x/oauth2@latest
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
go mod tidy
```

- [ ] **Step 7: Commit**

```bash
git add user-center/
git commit -m "feat: init user-center module with protobuf definition"
```

---

## Task 2: 模型层 + 领域层

**Files:**
- Create: `user-center/internal/model/po/mysql/user.go`
- Create: `user-center/internal/model/po/mysql/role.go`
- Create: `user-center/internal/model/po/mysql/user_role.go`
- Create: `user-center/internal/model/dto/user_dto.go`
- Create: `user-center/internal/domain/user/entity/user.go`
- Create: `user-center/internal/domain/user/repository/repository.go`

- [ ] **Step 1: 创建 User PO `user-center/internal/model/po/mysql/user.go`**

```go
package mysql

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	GoogleID  string `gorm:"type:varchar(128);uniqueIndex"`
	Email     string `gorm:"type:varchar(256);uniqueIndex;not null"`
	Name      string `gorm:"type:varchar(128)"`
	Avatar    string `gorm:"type:varchar(512)"`
	Status    string `gorm:"type:varchar(32);not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
```

- [ ] **Step 2: 创建 Role PO `user-center/internal/model/po/mysql/role.go`**

```go
package mysql

import "time"

type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(256)"`
	IsDefault   bool   `gorm:"type:bool;default:false"`
	CreatedAt   time.Time
}

func (Role) TableName() string {
	return "roles"
}
```

- [ ] **Step 3: 创建 UserRole PO `user-center/internal/model/po/mysql/user_role.go`**

```go
package mysql

import "time"

type UserRole struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"type:bigint;not null;index:idx_user_role,unique"`
	RoleID    int64 `gorm:"type:bigint;not null;index:idx_user_role,unique"`
	CreatedAt time.Time
}

func (UserRole) TableName() string {
	return "user_roles"
}
```

- [ ] **Step 4: 创建 DTO `user-center/internal/model/dto/user_dto.go`**

```go
package dto

type UserResp struct {
	ID     int64    `json:"id"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Status string   `json:"status"`
	Roles  []string `json:"roles"`
}

type LoginResp struct {
	Token string   `json:"token"`
	User  UserResp `json:"user"`
}

type UpdateUserStatusReq struct {
	Status string `json:"status"`
}

type AssignRoleReq struct {
	RoleName string `json:"roleName"`
}

type InitSuperAdminReq struct {
	Email string `json:"email"`
}

type UserFilter struct {
	Status   *string `json:"status,omitempty"`
	Keyword  *string `json:"keyword,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

type CheckInitResp struct {
	NeedInit bool `json:"needInit"`
}
```

- [ ] **Step 5: 创建 User 实体 `user-center/internal/domain/user/entity/user.go`**

```go
package entity

import (
	"errors"
	"time"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"

	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleOperator   = "operator"
	RoleViewer     = "viewer"
)

type User struct {
	ID        int64
	GoogleID  string
	Email     string
	Name      string
	Avatar    string
	Status    string
	Roles     []Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Role struct {
	ID          int64
	Name        string
	Description string
	IsDefault   bool
	CreatedAt   time.Time
}

func NewUser(googleID, email, name, avatar string) *User {
	now := time.Now()
	return &User{
		GoogleID:  googleID,
		Email:     email,
		Name:      name,
		Avatar:    avatar,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Disable() error {
	if u.HasRole(RoleSuperAdmin) {
		return errors.New("cannot disable super admin")
	}
	u.Status = StatusDisabled
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Enable() {
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
}

func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) RoleNames() []string {
	names := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		names[i] = r.Name
	}
	return names
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}
```

- [ ] **Step 6: 创建 Repository 接口 `user-center/internal/domain/user/repository/repository.go`**

```go
package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindByUserID(ctx context.Context, userID int64) ([]entity.Role, error)
	FindDefault(ctx context.Context) (*entity.Role, error)
	AssignRole(ctx context.Context, userID, roleID int64) error
	RemoveUserRoles(ctx context.Context, userID int64) error
	SaveIfNotExist(ctx context.Context, role *entity.Role) error
	HasSuperAdmin(ctx context.Context) (bool, error)
}
```

- [ ] **Step 7: 验证编译**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go build ./...
```

- [ ] **Step 8: Commit**

```bash
git add user-center/internal/model/ user-center/internal/domain/
git commit -m "feat(user-center): add model and domain layers"
```

---

## Task 3: Repository 实现 + 种子数据

**Files:**
- Create: `user-center/internal/repository/user_repo.go`
- Create: `user-center/internal/repository/role_repo.go`
- Create: `user-center/internal/seed/seed.go`

- [ ] **Step 1: 创建 User Repository `user-center/internal/repository/user_repo.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) Save(ctx context.Context, user *entity.User) error {
	record := &po.User{
		GoogleID: user.GoogleID,
		Email:    user.Email,
		Name:     user.Name,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	user.ID = record.ID
	return nil
}

func (r *UserRepoImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var record po.User
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return toUserEntity(&record), nil
}

func (r *UserRepoImpl) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	var record po.User
	if err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&record).Error; err != nil {
		return nil, err
	}
	return toUserEntity(&record), nil
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var record po.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&record).Error; err != nil {
		return nil, err
	}
	return toUserEntity(&record), nil
}

func (r *UserRepoImpl) Update(ctx context.Context, user *entity.User) error {
	record := &po.User{
		ID:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
		Status: user.Status,
	}
	if err := r.db.WithContext(ctx).Model(record).Updates(map[string]interface{}{
		"name":   record.Name,
		"avatar": record.Avatar,
		"status": record.Status,
	}).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func toUserEntity(p *po.User) *entity.User {
	return &entity.User{
		ID:        p.ID,
		GoogleID:  p.GoogleID,
		Email:     p.Email,
		Name:      p.Name,
		Avatar:    p.Avatar,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
```

- [ ] **Step 2: 创建 Role Repository `user-center/internal/repository/role_repo.go`**

```go
package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"gorm.io/gorm"
)

type RoleRepoImpl struct {
	db *gorm.DB
}

func NewRoleRepoImpl(db *gorm.DB) *RoleRepoImpl {
	return &RoleRepoImpl{db: db}
}

func (r *RoleRepoImpl) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var record po.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&record).Error; err != nil {
		return nil, err
	}
	return toRoleEntity(&record), nil
}

func (r *RoleRepoImpl) FindByUserID(ctx context.Context, userID int64) ([]entity.Role, error) {
	var roles []po.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("find roles by user id: %w", err)
	}
	result := make([]entity.Role, len(roles))
	for i, role := range roles {
		result[i] = *toRoleEntity(&role)
	}
	return result, nil
}

func (r *RoleRepoImpl) FindDefault(ctx context.Context) (*entity.Role, error) {
	var record po.Role
	if err := r.db.WithContext(ctx).Where("is_default = ?", true).First(&record).Error; err != nil {
		return nil, err
	}
	return toRoleEntity(&record), nil
}

func (r *RoleRepoImpl) AssignRole(ctx context.Context, userID, roleID int64) error {
	ur := &po.UserRole{UserID: userID, RoleID: roleID}
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).
		FirstOrCreate(ur).Error
}

func (r *RoleRepoImpl) RemoveUserRoles(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&po.UserRole{}).Error
}

func (r *RoleRepoImpl) SaveIfNotExist(ctx context.Context, role *entity.Role) error {
	record := &po.Role{
		Name:        role.Name,
		Description: role.Description,
		IsDefault:   role.IsDefault,
	}
	return r.db.WithContext(ctx).Where("name = ?", role.Name).FirstOrCreate(record).Error
}

func (r *RoleRepoImpl) HasSuperAdmin(ctx context.Context) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&po.UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", entity.RoleSuperAdmin).
		Count(&count).Error
	return count > 0, err
}

func toRoleEntity(p *po.Role) *entity.Role {
	return &entity.Role{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		IsDefault:   p.IsDefault,
		CreatedAt:   p.CreatedAt,
	}
}
```

- [ ] **Step 3: 创建种子数据 `user-center/internal/seed/seed.go`**

```go
package seed

import (
	"context"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/repository"
)

func SeedRoles(ctx context.Context, roleRepo repository.RoleRepository) {
	roles := []entity.Role{
		{Name: entity.RoleSuperAdmin, Description: "超级管理员", IsDefault: false},
		{Name: entity.RoleAdmin, Description: "管理员", IsDefault: false},
		{Name: entity.RoleOperator, Description: "运营人员", IsDefault: false},
		{Name: entity.RoleViewer, Description: "只读用户", IsDefault: true},
	}
	for _, role := range roles {
		if err := roleRepo.SaveIfNotExist(ctx, &role); err != nil {
			log.Printf("warn: seed role %s failed: %v", role.Name, err)
		}
	}
	log.Println("roles seeded successfully")
}
```

- [ ] **Step 4: 验证编译**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add user-center/internal/repository/ user-center/internal/seed/
git commit -m "feat(user-center): add repository implementations and role seed"
```

---

## Task 4: Application 层 (JWT + 用户业务逻辑)

**Files:**
- Create: `user-center/internal/application/jwt.go`
- Create: `user-center/internal/application/user_app.go`
- Create: `user-center/internal/config/api_config.go`
- Create: `user-center/internal/config/rpc_config.go`

- [ ] **Step 1: 创建 JWT 工具 `user-center/internal/application/jwt.go`**

```go
package application

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret string
	Expire int64 // seconds
}

type Claims struct {
	UserID int64    `json:"userId"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(cfg JWTConfig, userID int64, email, name string, roles []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(secret string, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
```

- [ ] **Step 2: 创建用户应用服务 `user-center/internal/application/user_app.go`**

```go
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/repository"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"gorm.io/gorm"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Avatar   string
}

type UserApp struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	jwtCfg   JWTConfig
	db       *gorm.DB
}

func NewUserApp(ur repository.UserRepository, rr repository.RoleRepository, jwtCfg JWTConfig, db *gorm.DB) *UserApp {
	return &UserApp{userRepo: ur, roleRepo: rr, jwtCfg: jwtCfg, db: db}
}

func (a *UserApp) LoginOrRegister(ctx context.Context, info GoogleUserInfo) (*entity.User, string, error) {
	user, err := a.userRepo.FindByGoogleID(ctx, info.GoogleID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", fmt.Errorf("find user: %w", err)
		}
		// 新用户注册
		user = entity.NewUser(info.GoogleID, info.Email, info.Name, info.Avatar)
		if err := a.userRepo.Save(ctx, user); err != nil {
			return nil, "", fmt.Errorf("save user: %w", err)
		}
		// 分配默认角色
		defaultRole, err := a.roleRepo.FindDefault(ctx)
		if err == nil {
			_ = a.roleRepo.AssignRole(ctx, user.ID, defaultRole.ID)
		}
	}

	if !user.IsActive() {
		return nil, "", fmt.Errorf("user is disabled")
	}

	// 加载角色
	roles, err := a.roleRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("find roles: %w", err)
	}
	user.Roles = roles

	token, err := GenerateToken(a.jwtCfg, user.ID, user.Email, user.Name, user.RoleNames())
	if err != nil {
		return nil, "", fmt.Errorf("generate token: %w", err)
	}

	return user, token, nil
}

func (a *UserApp) VerifyToken(ctx context.Context, tokenStr string) (*Claims, error) {
	return ParseToken(a.jwtCfg.Secret, tokenStr)
}

func (a *UserApp) GetUserInfo(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	roles, _ := a.roleRepo.FindByUserID(ctx, userID)
	user.Roles = roles
	return user, nil
}

func (a *UserApp) ListUsers(ctx context.Context, page, pageSize int) ([]po.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var total int64
	var users []po.User
	query := a.db.WithContext(ctx).Model(&po.User{})
	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	return users, total, nil
}

func (a *UserApp) UpdateStatus(ctx context.Context, userID int64, status string) error {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	switch status {
	case "active":
		user.Enable()
	case "disabled":
		if err := user.Disable(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid status: %s", status)
	}
	return a.userRepo.Update(ctx, user)
}

func (a *UserApp) AssignRole(ctx context.Context, userID int64, roleName string) error {
	role, err := a.roleRepo.FindByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("role not found: %s", roleName)
	}
	if err := a.roleRepo.RemoveUserRoles(ctx, userID); err != nil {
		return err
	}
	return a.roleRepo.AssignRole(ctx, userID, role.ID)
}

func (a *UserApp) InitSuperAdmin(ctx context.Context, userID int64) error {
	has, err := a.roleRepo.HasSuperAdmin(ctx)
	if err != nil {
		return err
	}
	if has {
		return fmt.Errorf("super admin already exists")
	}
	return a.AssignRole(ctx, userID, entity.RoleSuperAdmin)
}

func (a *UserApp) CheckRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	roles, err := a.roleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r.Name == roleName || r.Name == entity.RoleSuperAdmin {
			return true, nil
		}
	}
	return false, nil
}

func (a *UserApp) NeedInit(ctx context.Context) (bool, error) {
	has, err := a.roleRepo.HasSuperAdmin(ctx)
	return !has, err
}
```

- [ ] **Step 3: 创建配置结构体 `user-center/internal/config/api_config.go`**

```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ApiConfig struct {
	rest.RestConf
	MySQL         MySQLConfig
	JWT           JWTConfig
	Google        GoogleConfig
	UserCenterRpc zrpc.RpcClientConf
}

type MySQLConfig struct {
	DataSource string
}

type JWTConfig struct {
	Secret string
	Expire int64
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
```

- [ ] **Step 4: 创建 RPC 配置 `user-center/internal/config/rpc_config.go`**

```go
package config

import "github.com/zeromicro/go-zero/zrpc"

type RpcConfig struct {
	zrpc.RpcServerConf
	MySQL MySQLConfig
	JWT   JWTConfig
}
```

- [ ] **Step 5: 验证编译 + Commit**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go build ./...
git add user-center/internal/application/ user-center/internal/config/
git commit -m "feat(user-center): add application layer - JWT, user service, configs"
```

---

## Task 5: RPC 服务

**Files:**
- Create: `user-center/internal/server/rpc/usercenter_server.go`
- Create: `user-center/cmd/rpc/main.go`

- [ ] **Step 1: 创建 RPC 服务实现 `user-center/internal/server/rpc/usercenter_server.go`**

```go
package rpc

import (
	"context"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
)

type UserCenterServer struct {
	pb.UnimplementedUserCenterServer
	userApp *application.UserApp
}

func NewUserCenterServer(app *application.UserApp) *UserCenterServer {
	return &UserCenterServer{userApp: app}
}

func (s *UserCenterServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenReq) (*pb.VerifyTokenResp, error) {
	claims, err := s.userApp.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResp{
		UserId: claims.UserID,
		Email:  claims.Email,
		Name:   claims.Name,
		Roles:  claims.Roles,
	}, nil
}

func (s *UserCenterServer) CheckRole(ctx context.Context, req *pb.CheckRoleReq) (*pb.CheckRoleResp, error) {
	has, err := s.userApp.CheckRole(ctx, req.UserId, req.Role)
	if err != nil {
		return nil, err
	}
	return &pb.CheckRoleResp{HasRole: has}, nil
}

func (s *UserCenterServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoReq) (*pb.UserInfo, error) {
	user, err := s.userApp.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Avatar: user.Avatar,
		Status: user.Status,
		Roles:  user.RoleNames(),
	}, nil
}
```

- [ ] **Step 2: 创建 RPC 入口 `user-center/cmd/rpc/main.go`**

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	"github.com/yangboyi/ddd-dev/user-center/internal/repository"
	rpcServer "github.com/yangboyi/ddd-dev/user-center/internal/server/rpc"
	"github.com/yangboyi/ddd-dev/user-center/internal/seed"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/rpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.RpcConfig
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}

	if err := db.AutoMigrate(&po.User{}, &po.Role{}, &po.UserRole{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	userRepo := repository.NewUserRepoImpl(db)
	roleRepo := repository.NewRoleRepoImpl(db)

	jwtCfg := application.JWTConfig{Secret: c.JWT.Secret, Expire: c.JWT.Expire}
	userApp := application.NewUserApp(userRepo, roleRepo, jwtCfg, db)

	// Seed roles
	seed.SeedRoles(context.Background(), roleRepo)

	srv := zrpc.MustNewServer(c.RpcServerConf, func(s *grpc.Server) {
		pb.RegisterUserCenterServer(s, rpcServer.NewUserCenterServer(userApp))
	})
	defer srv.Stop()

	fmt.Printf("Starting user-center-rpc at %s...\n", c.ListenOn)
	srv.Start()
}
```

- [ ] **Step 3: 验证编译 + Commit**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go build ./...
git add user-center/cmd/rpc/ user-center/internal/server/rpc/
git commit -m "feat(user-center): add zRPC server with VerifyToken, CheckRole, GetUserInfo"
```

---

## Task 6: API 服务 (OAuth + 用户管理 + 中间件)

**Files:**
- Create: `user-center/internal/middleware/auth.go`
- Create: `user-center/internal/server/api/auth_handler.go`
- Create: `user-center/internal/server/api/user_handler.go`
- Create: `user-center/internal/server/api/init_handler.go`
- Create: `user-center/internal/server/api/routes.go`
- Create: `user-center/cmd/api/main.go`

- [ ] **Step 1: 创建 JWT 中间件 `user-center/internal/middleware/auth.go`**

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
)

type contextKey string

const (
	CtxUserID contextKey = "userId"
	CtxEmail  contextKey = "email"
	CtxRoles  contextKey = "roles"
)

func AuthMiddleware(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := application.ParseToken(jwtSecret, tokenStr)
			if err != nil {
				http.Error(w, `{"code":401,"message":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			ctx = context.WithValue(ctx, CtxRoles, claims.Roles)
			next(w, r.WithContext(ctx))
		}
	}
}

func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(CtxUserID).(int64)
	return v
}

func GetRoles(ctx context.Context) []string {
	v, _ := ctx.Value(CtxRoles).([]string)
	return v
}

func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role || r == "super_admin" {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: 创建 OAuth Handler `user-center/internal/server/api/auth_handler.go`**

```go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	userApp     *application.UserApp
	oauthConfig *oauth2.Config
}

func NewAuthHandler(userApp *application.UserApp, googleCfg config.GoogleConfig) *AuthHandler {
	return &AuthHandler{
		userApp: userApp,
		oauthConfig: &oauth2.Config{
			ClientID:     googleCfg.ClientID,
			ClientSecret: googleCfg.ClientSecret,
			RedirectURL:  googleCfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeJSON(w, 400, "missing code")
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		writeJSON(w, 500, fmt.Sprintf("exchange token: %v", err))
		return
	}

	client := h.oauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		writeJSON(w, 500, fmt.Sprintf("get user info: %v", err))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		writeJSON(w, 500, "parse google user info failed")
		return
	}

	user, jwtToken, err := h.userApp.LoginOrRegister(r.Context(), application.GoogleUserInfo{
		GoogleID: googleUser.ID,
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Avatar:   googleUser.Picture,
	})
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}

	// 重定向前端，带上 token
	frontendURL := fmt.Sprintf("http://localhost:3000/login/callback?token=%s", jwtToken)
	_ = user
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	user, err := h.userApp.GetUserInfo(r.Context(), userID)
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, dto.UserResp{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Avatar: user.Avatar,
		Status: user.Status,
		Roles:  user.RoleNames(),
	})
}

func writeJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"code": code, "message": msg})
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "message": "ok", "data": data})
}
```

- [ ] **Step 3: 创建用户管理 Handler `user-center/internal/server/api/user_handler.go`**

```go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
)

type UserHandler struct {
	userApp *application.UserApp
}

func NewUserHandler(userApp *application.UserApp) *UserHandler {
	return &UserHandler{userApp: userApp}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin") {
		writeJSON(w, 403, "forbidden")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	users, total, err := h.userApp.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, map[string]interface{}{"items": users, "total": total})
}

func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin") {
		writeJSON(w, 403, "forbidden")
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req dto.UpdateUserStatusReq
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.userApp.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "super_admin") {
		writeJSON(w, 403, "only super admin can assign roles")
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req dto.AssignRoleReq
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.userApp.AssignRole(r.Context(), id, req.RoleName); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}
```

- [ ] **Step 4: 创建超管初始化 Handler `user-center/internal/server/api/init_handler.go`**

```go
package api

import (
	"net/http"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
)

type InitHandler struct {
	userApp *application.UserApp
}

func NewInitHandler(userApp *application.UserApp) *InitHandler {
	return &InitHandler{userApp: userApp}
}

func (h *InitHandler) Check(w http.ResponseWriter, r *http.Request) {
	needInit, err := h.userApp.NeedInit(r.Context())
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, dto.CheckInitResp{NeedInit: needInit})
}

func (h *InitHandler) InitSuperAdmin(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if err := h.userApp.InitSuperAdmin(r.Context(), userID); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}
```

- [ ] **Step 5: 创建路由 `user-center/internal/server/api/routes.go`**

```go
package api

import (
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterRoutes(engine *rest.Server, auth *AuthHandler, user *UserHandler, init *InitHandler, jwtSecret string) {
	authMw := middleware.AuthMiddleware(jwtSecret)

	// 公开路由
	engine.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/auth/google/login", Handler: auth.GoogleLogin},
		{Method: http.MethodGet, Path: "/api/auth/google/callback", Handler: auth.GoogleCallback},
		{Method: http.MethodGet, Path: "/api/init/check", Handler: init.Check},
	})

	// 需要鉴权的路由
	engine.AddRoutes(rest.WithMiddleware(authMw,
		rest.Route{Method: http.MethodGet, Path: "/api/auth/me", Handler: auth.Me},
		rest.Route{Method: http.MethodPost, Path: "/api/init/super-admin", Handler: init.InitSuperAdmin},
		rest.Route{Method: http.MethodGet, Path: "/api/users", Handler: user.List},
		rest.Route{Method: http.MethodPut, Path: "/api/users/status", Handler: user.UpdateStatus},
		rest.Route{Method: http.MethodPut, Path: "/api/users/role", Handler: user.AssignRole},
	)...)
}
```

- [ ] **Step 6: 创建 API 入口 `user-center/cmd/api/main.go`**

```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	"github.com/yangboyi/ddd-dev/user-center/internal/repository"
	apiServer "github.com/yangboyi/ddd-dev/user-center/internal/server/api"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.ApiConfig
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.MySQL.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}

	if err := db.AutoMigrate(&po.User{}, &po.Role{}, &po.UserRole{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}

	userRepo := repository.NewUserRepoImpl(db)
	roleRepo := repository.NewRoleRepoImpl(db)
	jwtCfg := application.JWTConfig{Secret: c.JWT.Secret, Expire: c.JWT.Expire}
	userApp := application.NewUserApp(userRepo, roleRepo, jwtCfg, db)

	authHandler := apiServer.NewAuthHandler(userApp, c.Google)
	userHandler := apiServer.NewUserHandler(userApp)
	initHandler := apiServer.NewInitHandler(userApp)

	srv := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer srv.Stop()

	apiServer.RegisterRoutes(srv, authHandler, userHandler, initHandler, c.JWT.Secret)

	fmt.Printf("Starting user-center-api at %s:%d...\n", c.Host, c.Port)
	srv.Start()
}
```

- [ ] **Step 7: 验证编译 + Commit**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go build ./...
git add user-center/internal/middleware/ user-center/internal/server/api/ user-center/cmd/api/
git commit -m "feat(user-center): add REST API with OAuth, user management, auth middleware"
```

---

## Task 7: dropship-api 集成鉴权

**Files:**
- Modify: `backend/infra/config/config.go`
- Modify: `backend/etc/config.yaml`
- Create: `backend/internal/middleware/auth.go`
- Modify: `backend/internal/server/routes.go`
- Modify: `backend/internal/wire.go`
- Modify: `backend/main.go`

- [ ] **Step 1: 更新配置 `backend/infra/config/config.go`**

在 Config 结构体中新增 UserCenterRpc：

```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	MySQL         MySQLConfig
	Etcd          EtcdConfig  `json:",optional"`
	UserCenterRpc zrpc.RpcClientConf `json:",optional"`
}

type MySQLConfig struct {
	DataSource string
}

type EtcdConfig struct {
	Hosts []string `json:",optional"`
	Key   string   `json:",optional"`
}
```

- [ ] **Step 2: 更新 config.yaml**

在 `backend/etc/config.yaml` 末尾追加：

```yaml
UserCenterRpc:
  Etcd:
    Hosts:
      - "127.0.0.1:2379"
    Key: user-center.rpc
```

- [ ] **Step 3: 创建鉴权中间件 `backend/internal/middleware/auth.go`**

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
	"github.com/zeromicro/go-zero/zrpc"
)

type contextKey string

const (
	CtxUserID contextKey = "userId"
	CtxRoles  contextKey = "roles"
)

type AuthMiddleware struct {
	userCenterClient pb.UserCenterClient
}

func NewAuthMiddleware(conn zrpc.Client) *AuthMiddleware {
	return &AuthMiddleware{
		userCenterClient: pb.NewUserCenterClient(conn.Conn()),
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		resp, err := m.userCenterClient.VerifyToken(r.Context(), &pb.VerifyTokenReq{Token: tokenStr})
		if err != nil {
			http.Error(w, `{"code":401,"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxUserID, resp.UserId)
		ctx = context.WithValue(ctx, CtxRoles, resp.Roles)
		next(w, r.WithContext(ctx))
	}
}
```

- [ ] **Step 4: 更新 wire.go — 注入 RPC client**

修改 `backend/internal/wire.go`，增加 `authMiddleware` 参数：

```go
package internal

import (
	"github.com/yangboyi/ddd-dev/backend/internal/application"
	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	"github.com/yangboyi/ddd-dev/backend/internal/gateway"
	"github.com/yangboyi/ddd-dev/backend/internal/middleware"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
	repo "github.com/yangboyi/ddd-dev/backend/internal/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type Handlers struct {
	SourceItem     *server.SourceItemHandler
	Product        *server.ProductHandler
	Publish        *server.PublishHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitHandlers(db *gorm.DB, userCenterRpc zrpc.Client) *Handlers {
	sourceItemRepo := repo.NewSourceItemRepoImpl(db)
	productRepo := repo.NewProductRepoImpl(db)
	publishTaskRepo := repo.NewPublishTaskRepoImpl(db)

	sourceGateway := gateway.NewMockSourceGateway()
	targetGateway := gateway.NewMockTargetGateway()

	sourceItemQuery := queries.NewSourceItemQuery(db)

	publishDomainService := domainservice.NewPublishDomainService(productRepo, publishTaskRepo, targetGateway)

	sourceItemApp := application.NewSourceItemApp(sourceItemRepo, sourceGateway, sourceItemQuery)
	productApp := application.NewProductApp(productRepo, sourceItemRepo, db)
	publishApp := application.NewPublishApp(productRepo, publishDomainService, db)

	authMw := middleware.NewAuthMiddleware(userCenterRpc)

	return &Handlers{
		SourceItem:     server.NewSourceItemHandler(sourceItemApp),
		Product:        server.NewProductHandler(productApp),
		Publish:        server.NewPublishHandler(publishApp),
		AuthMiddleware: authMw,
	}
}
```

- [ ] **Step 5: 更新路由 — 应用鉴权中间件**

修改 `backend/internal/server/routes.go`：

```go
package server

import (
	"net/http"

	"github.com/yangboyi/ddd-dev/backend/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(engine *rest.Server, si *SourceItemHandler, p *ProductHandler, pub *PublishHandler, authMw *middleware.AuthMiddleware) {
	engine.AddRoutes(rest.WithMiddleware(authMw.Handle,
		rest.Route{Method: http.MethodPost, Path: "/api/source-items/import", Handler: si.Import},
		rest.Route{Method: http.MethodGet, Path: "/api/source-items", Handler: si.List},
		rest.Route{Method: http.MethodGet, Path: "/api/source-items/detail", Handler: si.GetByID},
		rest.Route{Method: http.MethodPut, Path: "/api/source-items/status", Handler: si.UpdateStatus},
		rest.Route{Method: http.MethodPost, Path: "/api/source-items/tag", Handler: si.AddTag},

		rest.Route{Method: http.MethodPost, Path: "/api/products/create-from-source", Handler: p.CreateFromSource},
		rest.Route{Method: http.MethodGet, Path: "/api/products", Handler: p.List},
		rest.Route{Method: http.MethodGet, Path: "/api/products/detail", Handler: p.GetByID},
		rest.Route{Method: http.MethodPut, Path: "/api/products", Handler: p.Update},
		rest.Route{Method: http.MethodPut, Path: "/api/products/ready", Handler: p.MarkReady},

		rest.Route{Method: http.MethodPost, Path: "/api/publish-tasks", Handler: pub.CreateTask},
		rest.Route{Method: http.MethodGet, Path: "/api/publish-tasks", Handler: pub.List},
	)...)
}
```

- [ ] **Step 6: 更新 main.go — 初始化 RPC client**

修改 `backend/main.go` 中 handlers 初始化：

```go
// 在 srv := rest.MustNewServer(...) 之前添加:
userCenterRpc := zrpc.MustNewClient(c.UserCenterRpc)

handlers := internal.InitHandlers(db, userCenterRpc)
server.RegisterRoutes(srv, handlers.SourceItem, handlers.Product, handlers.Publish, handlers.AuthMiddleware)
```

需要增加 import:
```go
"github.com/zeromicro/go-zero/zrpc"
```

- [ ] **Step 7: backend 需要引用 user-center proto**

在 `backend/go.mod` 中添加 replace 指令指向本地 user-center 模块：

```bash
cd /Users/yangboyi/github/ddd-dev/backend
go mod edit -require github.com/yangboyi/ddd-dev/user-center@v0.0.0
go mod edit -replace github.com/yangboyi/ddd-dev/user-center=../user-center
go mod tidy
```

- [ ] **Step 8: 验证编译 + Commit**

```bash
cd /Users/yangboyi/github/ddd-dev/backend && go build ./...
git add backend/
git commit -m "feat(dropship-api): add auth middleware via user-center RPC"
```

---

## Task 8: 前端 — 登录 + 鉴权

**Files:**
- Create: `frontend/src/lib/auth.ts`
- Modify: `frontend/src/lib/api.ts`
- Create: `frontend/src/app/login/page.tsx`
- Create: `frontend/src/app/login/callback/page.tsx`
- Modify: `frontend/src/app/layout.tsx`
- Modify: `frontend/src/components/layout/sidebar.tsx`

- [ ] **Step 1: 创建 auth 工具 `frontend/src/lib/auth.ts`**

```typescript
export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
}

export function setToken(token: string) {
  localStorage.setItem("token", token);
}

export function clearToken() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}

export function isLoggedIn(): boolean {
  return !!getToken();
}

export function getUser(): any | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem("user");
  return raw ? JSON.parse(raw) : null;
}

export function setUser(user: any) {
  localStorage.setItem("user", JSON.stringify(user));
}
```

- [ ] **Step 2: 修改 api.ts — 增加 auth header + user-center API**

替换 `frontend/src/lib/api.ts`：

```typescript
import { getToken, clearToken } from "./auth";

const API_BASE = "http://localhost:8888/api";
const USER_CENTER_API = "http://localhost:8880/api";

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

async function request<T>(base: string, path: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${base}${path}`, { headers, ...options });
  if (res.status === 401) {
    clearToken();
    window.location.href = "/login";
    throw new Error("unauthorized");
  }
  const json: ApiResponse<T> = await res.json();
  if (json.code !== 0) {
    throw new Error(json.message);
  }
  return json.data;
}

export const api = {
  // 货源
  sourceItems: {
    import: (data: { platform: string; sourceUrl: string }) =>
      request(API_BASE, "/source-items/import", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/source-items?${new URLSearchParams(params)}`),
    updateStatus: (id: number, status: string) =>
      request(API_BASE, `/source-items/status?id=${id}`, { method: "PUT", body: JSON.stringify({ status }) }),
    addTag: (id: number, tag: string) =>
      request(API_BASE, `/source-items/tag?id=${id}`, { method: "POST", body: JSON.stringify({ tag }) }),
  },
  // 商品
  products: {
    createFromSource: (sourceItemId: number) =>
      request(API_BASE, "/products/create-from-source", { method: "POST", body: JSON.stringify({ sourceItemId }) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/products?${new URLSearchParams(params)}`),
    get: (id: number) => request(API_BASE, `/products/detail?id=${id}`),
    update: (id: number, data: any) =>
      request(API_BASE, `/products?id=${id}`, { method: "PUT", body: JSON.stringify(data) }),
    markReady: (id: number) =>
      request(API_BASE, `/products/ready?id=${id}`, { method: "PUT" }),
  },
  // 发品
  publishTasks: {
    create: (data: { productId: number; targetPlatform: string; categoryId: string; freightTemplate: string }) =>
      request(API_BASE, "/publish-tasks", { method: "POST", body: JSON.stringify(data) }),
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(API_BASE, `/publish-tasks?${new URLSearchParams(params)}`),
  },
  // 用户中心
  auth: {
    me: () => request<any>(USER_CENTER_API, "/auth/me"),
    checkInit: () => request<{ needInit: boolean }>(USER_CENTER_API, "/init/check"),
    initSuperAdmin: () => request(USER_CENTER_API, "/init/super-admin", { method: "POST" }),
  },
  users: {
    list: (params: Record<string, string>) =>
      request<{ items: any[]; total: number }>(USER_CENTER_API, `/users?${new URLSearchParams(params)}`),
    updateStatus: (id: number, status: string) =>
      request(USER_CENTER_API, `/users/status?id=${id}`, { method: "PUT", body: JSON.stringify({ status }) }),
    assignRole: (id: number, roleName: string) =>
      request(USER_CENTER_API, `/users/role?id=${id}`, { method: "PUT", body: JSON.stringify({ roleName }) }),
  },
};
```

- [ ] **Step 3: 创建登录页 `frontend/src/app/login/page.tsx`**

```tsx
"use client";

import { Button } from "@/components/ui/button";

export default function LoginPage() {
  const handleGoogleLogin = () => {
    window.location.href = "http://localhost:8880/api/auth/google/login";
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-sm w-full space-y-6 p-8 bg-white rounded-xl shadow-md">
        <div className="text-center">
          <h1 className="text-2xl font-bold">代发工具</h1>
          <p className="text-gray-500 mt-2">电商一键代发运营工具</p>
        </div>
        <Button onClick={handleGoogleLogin} className="w-full" size="lg">
          Sign in with Google
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: 创建 OAuth 回调页 `frontend/src/app/login/callback/page.tsx`**

```tsx
"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { setToken, setUser } from "@/lib/auth";
import { api } from "@/lib/api";

export default function CallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setToken(token);
      api.auth.me().then((user) => {
        setUser(user);
        router.replace("/sources");
      });
    }
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <p>登录中...</p>
    </div>
  );
}
```

- [ ] **Step 5: 更新 layout.tsx — auth guard**

```tsx
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthGuard } from "@/components/layout/auth-guard";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "代发工具",
  description: "电商一键代发运营工具",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        <AuthGuard>{children}</AuthGuard>
      </body>
    </html>
  );
}
```

- [ ] **Step 6: 创建 AuthGuard `frontend/src/components/layout/auth-guard.tsx`**

```tsx
"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { isLoggedIn } from "@/lib/auth";
import { Sidebar } from "./sidebar";

const publicPaths = ["/login", "/login/callback"];

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [checked, setChecked] = useState(false);
  const isPublic = publicPaths.some((p) => pathname.startsWith(p));

  useEffect(() => {
    if (!isPublic && !isLoggedIn()) {
      router.replace("/login");
    } else {
      setChecked(true);
    }
  }, [pathname, isPublic, router]);

  if (!checked) return null;

  if (isPublic) {
    return <>{children}</>;
  }

  return (
    <div className="flex">
      <Sidebar />
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
```

- [ ] **Step 7: 更新 sidebar — 增加用户管理 + 角色控制 + 登出**

```tsx
"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { cn } from "@/lib/utils";
import { clearToken, getUser } from "@/lib/auth";
import { Button } from "@/components/ui/button";

const navItems = [
  { href: "/sources", label: "货源管理", icon: "📦", roles: null },
  { href: "/products", label: "商品管理", icon: "🏷️", roles: null },
  { href: "/publish", label: "发品任务", icon: "🚀", roles: null },
  { href: "/users", label: "用户管理", icon: "👥", roles: ["super_admin", "admin"] },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const user = getUser();
  const userRoles: string[] = user?.roles || [];

  const handleLogout = () => {
    clearToken();
    router.replace("/login");
  };

  const visibleItems = navItems.filter(
    (item) => !item.roles || item.roles.some((r) => userRoles.includes(r))
  );

  return (
    <aside className="w-56 border-r bg-gray-50 p-4 min-h-screen flex flex-col">
      <h1 className="text-lg font-bold mb-6 px-2">代发工具</h1>
      <nav className="space-y-1 flex-1">
        {visibleItems.map((item) => (
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
      {user && (
        <div className="border-t pt-4 mt-4">
          <div className="px-2 text-sm text-gray-600 mb-2 truncate">{user.email}</div>
          <Button variant="outline" size="sm" className="w-full" onClick={handleLogout}>
            退出登录
          </Button>
        </div>
      )}
    </aside>
  );
}
```

- [ ] **Step 8: 更新首页重定向**

`frontend/src/app/page.tsx` 保持不变（已重定向到 /sources）。

- [ ] **Step 9: Commit**

```bash
git add frontend/
git commit -m "feat(frontend): add login page, auth guard, user-center API integration"
```

---

## Task 9: 前端 — 超管初始化 + 用户管理页

**Files:**
- Create: `frontend/src/app/init/page.tsx`
- Create: `frontend/src/app/users/page.tsx`
- Create: `frontend/src/components/users/user-table.tsx`

- [ ] **Step 1: 创建超管初始化页 `frontend/src/app/init/page.tsx`**

```tsx
"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUser, setUser } from "@/lib/auth";

export default function InitPage() {
  const router = useRouter();
  const [needInit, setNeedInit] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(false);
  const user = getUser();

  useEffect(() => {
    api.auth.checkInit().then((res) => setNeedInit(res.needInit));
  }, []);

  const handleInit = async () => {
    setLoading(true);
    try {
      await api.auth.initSuperAdmin();
      const updatedUser = await api.auth.me();
      setUser(updatedUser);
      alert("超级管理员初始化成功！");
      router.replace("/sources");
    } catch (e: any) {
      alert(e.message);
    } finally {
      setLoading(false);
    }
  };

  if (needInit === null) return <div className="p-8">检查中...</div>;
  if (!needInit) {
    router.replace("/sources");
    return null;
  }

  return (
    <div className="max-w-md mx-auto mt-20 p-8 bg-white rounded-xl shadow-md space-y-6">
      <h2 className="text-xl font-bold">初始化超级管理员</h2>
      <p className="text-gray-500">
        系统首次使用，需要初始化超级管理员。当前登录用户
        <strong> {user?.email} </strong>
        将被设为超级管理员。
      </p>
      <Button onClick={handleInit} disabled={loading} className="w-full">
        {loading ? "初始化中..." : "确认初始化"}
      </Button>
    </div>
  );
}
```

- [ ] **Step 2: 创建用户表格 `frontend/src/components/users/user-table.tsx`**

```tsx
"use client";

import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { getUser } from "@/lib/auth";

interface User {
  ID: number;
  Email: string;
  Name: string;
  Status: string;
}

const roleOptions = ["super_admin", "admin", "operator", "viewer"];

export function UserTable({ items, onRefresh }: { items: User[]; onRefresh: () => void }) {
  const currentUser = getUser();
  const isSuperAdmin = currentUser?.roles?.includes("super_admin");

  const handleToggleStatus = async (id: number, currentStatus: string) => {
    const newStatus = currentStatus === "active" ? "disabled" : "active";
    await api.users.updateStatus(id, newStatus);
    onRefresh();
  };

  const handleAssignRole = async (id: number, roleName: string) => {
    await api.users.assignRole(id, roleName);
    onRefresh();
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>邮箱</TableHead>
          <TableHead>名称</TableHead>
          <TableHead>状态</TableHead>
          <TableHead>操作</TableHead>
          {isSuperAdmin && <TableHead>角色分配</TableHead>}
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => (
          <TableRow key={item.ID}>
            <TableCell>{item.Email}</TableCell>
            <TableCell>{item.Name}</TableCell>
            <TableCell>
              <Badge variant={item.Status === "active" ? "default" : "destructive"}>
                {item.Status === "active" ? "正常" : "禁用"}
              </Badge>
            </TableCell>
            <TableCell>
              <Button size="sm" variant="outline" onClick={() => handleToggleStatus(item.ID, item.Status)}>
                {item.Status === "active" ? "禁用" : "启用"}
              </Button>
            </TableCell>
            {isSuperAdmin && (
              <TableCell>
                <select
                  className="border rounded px-2 py-1 text-sm"
                  onChange={(e) => handleAssignRole(item.ID, e.target.value)}
                  defaultValue=""
                >
                  <option value="" disabled>选择角色</option>
                  {roleOptions.map((r) => (
                    <option key={r} value={r}>{r}</option>
                  ))}
                </select>
              </TableCell>
            )}
          </TableRow>
        ))}
        {items.length === 0 && (
          <TableRow>
            <TableCell colSpan={5} className="text-center text-gray-400 py-8">暂无用户</TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
```

- [ ] **Step 3: 创建用户管理页 `frontend/src/app/users/page.tsx`**

```tsx
"use client";

import { useCallback, useEffect, useState } from "react";
import { UserTable } from "@/components/users/user-table";
import { api } from "@/lib/api";

export default function UsersPage() {
  const [items, setItems] = useState<any[]>([]);

  const fetchData = useCallback(async () => {
    const res = await api.users.list({ page: "1", pageSize: "50" });
    setItems(res.items || []);
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">用户管理</h2>
      <UserTable items={items} onRefresh={fetchData} />
    </div>
  );
}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/
git commit -m "feat(frontend): add super admin init page and user management page"
```

---

## Task 10: 端到端验证

- [ ] **Step 1: 启动 user-center-rpc**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go run cmd/rpc/main.go -f etc/rpc.yaml
```

Expected: `Starting user-center-rpc at 0.0.0.0:8881...` + roles seeded

- [ ] **Step 2: 启动 user-center-api**

```bash
cd /Users/yangboyi/github/ddd-dev/user-center && go run cmd/api/main.go -f etc/api.yaml
```

Expected: `Starting user-center-api at 0.0.0.0:8880...`

- [ ] **Step 3: 启动 dropship-api**

```bash
cd /Users/yangboyi/github/ddd-dev/backend && go run main.go -f etc/config.yaml
```

Expected: 正常启动，通过 etcd 发现 user-center-rpc

- [ ] **Step 4: 验证鉴权拦截**

```bash
curl -s "http://localhost:8888/api/source-items?page=1&pageSize=10"
```

Expected: 返回 401 unauthorized（无 token）

- [ ] **Step 5: 验证 Google OAuth 流程**

浏览器打开 http://localhost:3000 → 自动跳转 /login → 点击 Sign in with Google → 完成授权 → 回调带 token → 自动跳转 /sources

- [ ] **Step 6: 验证超管初始化**

首次登录后访问 http://localhost:3000/init → 点击初始化 → 当前用户成为 super_admin

- [ ] **Step 7: 修复发现的问题（如有）**

```bash
git add -A && git commit -m "fix: resolve issues found during e2e testing"
```
