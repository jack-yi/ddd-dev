package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
)

type PublishConfig struct {
	CategoryID      string
	FreightTemplate string
}

type PublishTask struct {
	ID                int64
	ProductID         int64
	TargetPlatform    string
	PlatformProductID string
	PublishConfig     PublishConfig
	Status            string
	ErrorMessage      string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewPublishTask(productID int64, targetPlatform string, config PublishConfig) *PublishTask {
	now := time.Now()
	return &PublishTask{
		ProductID:      productID,
		TargetPlatform: targetPlatform,
		PublishConfig:  config,
		Status:         consts.PublishTaskStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (t *PublishTask) MarkPublishing() error {
	if t.Status != consts.PublishTaskStatusPending {
		return errors.New("task must be pending to start publishing")
	}
	t.Status = consts.PublishTaskStatusPublishing
	t.UpdatedAt = time.Now()
	return nil
}

func (t *PublishTask) MarkSuccess(platformProductID string) {
	t.Status = consts.PublishTaskStatusSuccess
	t.PlatformProductID = platformProductID
	t.UpdatedAt = time.Now()
}

func (t *PublishTask) MarkFailed(errMsg string) {
	t.Status = consts.PublishTaskStatusFailed
	t.ErrorMessage = errMsg
	t.UpdatedAt = time.Now()
}
