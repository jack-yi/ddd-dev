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
