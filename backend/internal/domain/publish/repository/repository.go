package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
)

type PublishTaskRepository interface {
	Save(ctx context.Context, task *entity.PublishTask) error
	FindByID(ctx context.Context, id int64) (*entity.PublishTask, error)
	Update(ctx context.Context, task *entity.PublishTask) error
}
