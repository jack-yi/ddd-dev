package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/repository"
	po "github.com/yangboyi/ddd-dev/user-center/internal/model/po/mysql"
	"golang.org/x/crypto/bcrypt"
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
		user = entity.NewUser(info.GoogleID, info.Email, info.Name, info.Avatar)
		if err := a.userRepo.Save(ctx, user); err != nil {
			return nil, "", fmt.Errorf("save user: %w", err)
		}
		defaultRole, err := a.roleRepo.FindDefault(ctx)
		if err == nil {
			_ = a.roleRepo.AssignRole(ctx, user.ID, defaultRole.ID)
		}
	}
	if !user.IsActive() {
		return nil, "", fmt.Errorf("user is disabled")
	}
	roles, _ := a.roleRepo.FindByUserID(ctx, user.ID)
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

func (a *UserApp) LoginByPassword(ctx context.Context, username, password string) (*entity.User, string, error) {
	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid username or password")
	}
	if !checkPassword(user.PasswordHash, password) {
		return nil, "", fmt.Errorf("invalid username or password")
	}
	if !user.IsActive() {
		return nil, "", fmt.Errorf("user is disabled")
	}
	roles, _ := a.roleRepo.FindByUserID(ctx, user.ID)
	user.Roles = roles
	token, err := GenerateToken(a.jwtCfg, user.ID, user.Email, user.Name, user.RoleNames())
	if err != nil {
		return nil, "", fmt.Errorf("generate token: %w", err)
	}
	return user, token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
