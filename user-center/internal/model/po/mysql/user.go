package mysql

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	GoogleID  string `gorm:"type:varchar(128);uniqueIndex"`
	Email     string `gorm:"type:varchar(256);uniqueIndex;not null"`
	Name      string `gorm:"type:varchar(128)"`
	Avatar       string `gorm:"type:varchar(512)"`
	Username     string `gorm:"type:varchar(128);uniqueIndex"`
	PasswordHash string `gorm:"type:varchar(256)"`
	Status       string `gorm:"type:varchar(32);not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string { return "users" }
