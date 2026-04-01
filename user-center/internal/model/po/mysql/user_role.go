package mysql

import "time"

type UserRole struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"type:bigint;not null;uniqueIndex:idx_user_role"`
	RoleID    int64 `gorm:"type:bigint;not null;uniqueIndex:idx_user_role"`
	CreatedAt time.Time
}

func (UserRole) TableName() string { return "user_roles" }
