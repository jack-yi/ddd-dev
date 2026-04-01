package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
)

type SourceItemRepository interface {
	Save(ctx context.Context, item *entity.SourceItem) error
	FindByID(ctx context.Context, id int64) (*entity.SourceItem, error)
	Update(ctx context.Context, item *entity.SourceItem) error
}
