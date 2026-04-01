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
		Where("user_roles.user_id = ?", userID).Find(&roles).Error
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
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).FirstOrCreate(ur).Error
}

func (r *RoleRepoImpl) RemoveUserRoles(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&po.UserRole{}).Error
}

func (r *RoleRepoImpl) SaveIfNotExist(ctx context.Context, role *entity.Role) error {
	record := &po.Role{Name: role.Name, Description: role.Description, IsDefault: role.IsDefault}
	return r.db.WithContext(ctx).Where("name = ?", role.Name).FirstOrCreate(record).Error
}

func (r *RoleRepoImpl) HasSuperAdmin(ctx context.Context) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&po.UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", entity.RoleSuperAdmin).Count(&count).Error
	return count > 0, err
}

func toRoleEntity(p *po.Role) *entity.Role {
	return &entity.Role{ID: p.ID, Name: p.Name, Description: p.Description, IsDefault: p.IsDefault, CreatedAt: p.CreatedAt}
}
