package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PublishConfig struct {
	CategoryID      string `json:"categoryId"`
	FreightTemplate string `json:"freightTemplate"`
}

func (p PublishConfig) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PublishConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

type PublishTask struct {
	ID                int64         `gorm:"primaryKey;autoIncrement"`
	ProductID         int64         `gorm:"type:bigint;not null"`
	TargetPlatform    string        `gorm:"type:varchar(32);not null"`
	PlatformProductID string        `gorm:"type:varchar(128)"`
	PublishConfig     PublishConfig `gorm:"type:json"`
	Status            string        `gorm:"type:varchar(32);not null;default:'pending'"`
	ErrorMessage      string        `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (PublishTask) TableName() string {
	return "publish_tasks"
}
