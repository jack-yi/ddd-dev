package application

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/platform"
	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
	"github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
)

type SourceItemApp struct {
	repo          repository.SourceItemRepository
	sourceGateway platform.SourcePlatformGateway
	query         *queries.SourceItemQuery
}

func NewSourceItemApp(
	repo repository.SourceItemRepository,
	gw platform.SourcePlatformGateway,
	query *queries.SourceItemQuery,
) *SourceItemApp {
	return &SourceItemApp{repo: repo, sourceGateway: gw, query: query}
}

func (a *SourceItemApp) Import(ctx context.Context, req *dto.ImportSourceItemReq) (*entity.SourceItem, error) {
	product, err := a.sourceGateway.FetchProduct(ctx, req.SourceURL)
	if err != nil {
		return nil, fmt.Errorf("fetch product from source: %w", err)
	}

	item := entity.NewSourceItem(
		req.Platform, req.SourceURL, product.ExternalID,
		product.Title, product.Description, product.Images,
		entity.Price{Min: product.PriceMin, Max: product.PriceMax},
		entity.Supplier{Name: product.Supplier.Name, Rating: product.Supplier.Rating, Region: product.Supplier.Region},
		product.Category, product.SalesVolume, product.MinOrder,
	)

	if err := a.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("save source item: %w", err)
	}
	return item, nil
}

func (a *SourceItemApp) GetByID(ctx context.Context, id int64) (*entity.SourceItem, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *SourceItemApp) List(ctx context.Context, filter *dto.SourceItemFilter) (*queries.ListResult, error) {
	return a.query.List(ctx, filter)
}

func (a *SourceItemApp) UpdateStatus(ctx context.Context, id int64, status string) error {
	item, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find source item: %w", err)
	}

	switch status {
	case "selected":
		if err := item.Select(); err != nil {
			return err
		}
	case "ignored":
		if err := item.Ignore(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	return a.repo.Update(ctx, item)
}

func (a *SourceItemApp) AddTag(ctx context.Context, id int64, tag string) error {
	item, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find source item: %w", err)
	}
	item.AddTag(tag)
	return a.repo.Update(ctx, item)
}
