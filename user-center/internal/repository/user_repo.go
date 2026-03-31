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
	record := &po.User{GoogleID: user.GoogleID, Email: user.Email, Name: user.Name, Avatar: user.Avatar, Status: user.Status}
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
	return r.db.WithContext(ctx).Model(&po.User{ID: user.ID}).Updates(map[string]interface{}{
		"name": user.Name, "avatar": user.Avatar, "status": user.Status,
	}).Error
}

func toUserEntity(p *po.User) *entity.User {
	return &entity.User{
		ID: p.ID, GoogleID: p.GoogleID, Email: p.Email, Name: p.Name,
		Avatar: p.Avatar, Status: p.Status, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
