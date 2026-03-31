package repository

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id int64) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
}
