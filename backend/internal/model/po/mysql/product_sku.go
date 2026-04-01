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
