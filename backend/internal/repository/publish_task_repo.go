package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type PublishTaskRepoImpl struct {
	db *gorm.DB
}

func NewPublishTaskRepoImpl(db *gorm.DB) *PublishTaskRepoImpl {
	return &PublishTaskRepoImpl{db: db}
}

func (r *PublishTaskRepoImpl) Save(ctx context.Context, task *entity.PublishTask) error {
	record := toPublishTaskPO(task)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create publish task: %w", err)
	}
	task.ID = record.ID
	return nil
}

func (r *PublishTaskRepoImpl) FindByID(ctx context.Context, id int64) (*entity.PublishTask, error) {
	var record po.PublishTask
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find publish task by id: %w", err)
	}
	return toPublishTaskEntity(&record), nil
}

func (r *PublishTaskRepoImpl) Update(ctx context.Context, task *entity.PublishTask) error {
	record := toPublishTaskPO(task)
	record.ID = task.ID
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update publish task: %w", err)
	}
	return nil
}

func toPublishTaskPO(e *entity.PublishTask) *po.PublishTask {
	return &po.PublishTask{
		ID:                e.ID,
		ProductID:         e.ProductID,
		TargetPlatform:    e.TargetPlatform,
		PlatformProductID: e.PlatformProductID,
		PublishConfig: po.PublishConfig{
			CategoryID:      e.PublishConfig.CategoryID,
			FreightTemplate: e.PublishConfig.FreightTemplate,
		},
		Status:       e.Status,
		ErrorMessage: e.ErrorMessage,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toPublishTaskEntity(p *po.PublishTask) *entity.PublishTask {
	return &entity.PublishTask{
		ID:                p.ID,
		ProductID:         p.ProductID,
		TargetPlatform:    p.TargetPlatform,
		PlatformProductID: p.PlatformProductID,
		PublishConfig: entity.PublishConfig{
			CategoryID:      p.PublishConfig.CategoryID,
			FreightTemplate: p.PublishConfig.FreightTemplate,
		},
		Status:       p.Status,
		ErrorMessage: p.ErrorMessage,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
