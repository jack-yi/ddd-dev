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
