package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type SourceItemRepoImpl struct {
	db *gorm.DB
}

func NewSourceItemRepoImpl(db *gorm.DB) *SourceItemRepoImpl {
	return &SourceItemRepoImpl{db: db}
}

func (r *SourceItemRepoImpl) Save(ctx context.Context, item *entity.SourceItem) error {
	record := toSourceItemPO(item)
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create source item: %w", err)
	}
	item.ID = record.ID
	return nil
}

func (r *SourceItemRepoImpl) FindByID(ctx context.Context, id int64) (*entity.SourceItem, error) {
	var record po.SourceItem
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find source item by id: %w", err)
	}
	return toSourceItemEntity(&record), nil
}

func (r *SourceItemRepoImpl) Update(ctx context.Context, item *entity.SourceItem) error {
	record := toSourceItemPO(item)
	record.ID = item.ID
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update source item: %w", err)
	}
	return nil
}

func toSourceItemPO(e *entity.SourceItem) *po.SourceItem {
	return &po.SourceItem{
		ID:          e.ID,
		Platform:    e.Platform,
		SourceURL:   e.SourceURL,
		ExternalID:  e.ExternalID,
		Title:       e.Title,
		Description: e.Description,
		Images:      po.StringSlice(e.Images),
		PriceMin:    e.Price.Min,
		PriceMax:    e.Price.Max,
		Supplier: po.SupplierInfo{
			Name:   e.Supplier.Name,
			Rating: e.Supplier.Rating,
			Region: e.Supplier.Region,
		},
		Category:    e.Category,
		Tags:        po.StringSlice(e.Tags),
		SalesVolume: e.SalesVolume,
		MinOrder:    e.MinOrder,
		Status:      e.Status,
		FetchedAt:   e.FetchedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toSourceItemEntity(p *po.SourceItem) *entity.SourceItem {
	return &entity.SourceItem{
		ID:          p.ID,
		Platform:    p.Platform,
		SourceURL:   p.SourceURL,
		ExternalID:  p.ExternalID,
		Title:       p.Title,
		Description: p.Description,
		Images:      []string(p.Images),
		Price:       entity.Price{Min: p.PriceMin, Max: p.PriceMax},
		Supplier: entity.Supplier{
			Name:   p.Supplier.Name,
			Rating: p.Supplier.Rating,
			Region: p.Supplier.Region,
		},
		Category:    p.Category,
		Tags:        []string(p.Tags),
		SalesVolume: p.SalesVolume,
		MinOrder:    p.MinOrder,
		Status:      p.Status,
		FetchedAt:   p.FetchedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
