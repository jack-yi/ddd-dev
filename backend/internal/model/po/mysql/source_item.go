package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type SupplierInfo struct {
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	Region string  `json:"region"`
}

func (s SupplierInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SupplierInfo) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type SourceItem struct {
	ID          int64        `gorm:"primaryKey;autoIncrement"`
	Platform    string       `gorm:"type:varchar(32);not null"`
	SourceURL   string       `gorm:"type:varchar(512);not null"`
	ExternalID  string       `gorm:"type:varchar(128)"`
	Title       string       `gorm:"type:varchar(256)"`
	Description string       `gorm:"type:text"`
	Images      StringSlice  `gorm:"type:json"`
	PriceMin    float64      `gorm:"type:decimal(10,2)"`
	PriceMax    float64      `gorm:"type:decimal(10,2)"`
	Supplier    SupplierInfo `gorm:"type:json"`
	Category    string       `gorm:"type:varchar(128)"`
	Tags        StringSlice  `gorm:"type:json"`
	SalesVolume int          `gorm:"type:int"`
	MinOrder    int          `gorm:"type:int"`
	Status      string       `gorm:"type:varchar(32);not null;default:'new'"`
	FetchedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SourceItem) TableName() string {
	return "source_items"
}
