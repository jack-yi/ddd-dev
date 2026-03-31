package queries

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type SourceItemQuery struct {
	db *gorm.DB
}

func NewSourceItemQuery(db *gorm.DB) *SourceItemQuery {
	return &SourceItemQuery{db: db}
}

type ListResult struct {
	Items []po.SourceItem
	Total int64
}

func (q *SourceItemQuery) List(ctx context.Context, filter *dto.SourceItemFilter) (*ListResult, error) {
	query := q.db.WithContext(ctx).Model(&po.SourceItem{})

	if filter.Platform != nil {
		query = query.Where("platform = ?", *filter.Platform)
	}
	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.PriceMin != nil {
		query = query.Where("price_min >= ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		query = query.Where("price_max <= ?", *filter.PriceMax)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", *filter.Keyword))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count source items: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.SourceItem
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list source items: %w", err)
	}

	return &ListResult{Items: items, Total: total}, nil
}
