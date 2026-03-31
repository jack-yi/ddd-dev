package mysql

import "time"

type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(256)"`
	IsDefault   bool   `gorm:"type:bool;default:false"`
	CreatedAt   time.Time
}

func (Role) TableName() string { return "roles" }
