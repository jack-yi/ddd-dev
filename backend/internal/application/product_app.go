package application

import (
	"context"
	"fmt"

	productEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	sourceRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type ProductApp struct {
	productRepo productRepo.ProductRepository
	sourceRepo  sourceRepo.SourceItemRepository
	db          *gorm.DB
}

func NewProductApp(pr productRepo.ProductRepository, sr sourceRepo.SourceItemRepository, db *gorm.DB) *ProductApp {
	return &ProductApp{productRepo: pr, sourceRepo: sr, db: db}
}

func (a *ProductApp) CreateFromSource(ctx context.Context, sourceItemID int64) (*productEntity.Product, error) {
	source, err := a.sourceRepo.FindByID(ctx, sourceItemID)
	if err != nil {
		return nil, fmt.Errorf("find source item: %w", err)
	}

	product := productEntity.CreateFromSource(source)
	if err := a.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}
	return product, nil
}

func (a *ProductApp) GetByID(ctx context.Context, id int64) (*productEntity.Product, error) {
	return a.productRepo.FindByID(ctx, id)
}

func (a *ProductApp) List(ctx context.Context, filter *dto.ProductFilter) ([]po.Product, int64, error) {
	query := a.db.WithContext(ctx).Model(&po.Product{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *filter.Keyword))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.Product
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}

	return items, total, nil
}

func (a *ProductApp) Update(ctx context.Context, id int64, req *dto.UpdateProductReq) (*productEntity.Product, error) {
	product, err := a.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	product.EditInfo(req.Name, req.Description, req.Images, req.CategoryID)

	if req.CostPrice != nil || req.SellPrice != nil {
		product.SetPrice(req.CostPrice, req.SellPrice)
	}

	if req.SKUs != nil {
		skus := make([]productEntity.SKU, len(req.SKUs))
		for i, s := range req.SKUs {
			skus[i] = productEntity.SKU{
				ID:        s.ID,
				SpecName:  s.SpecName,
				SpecValue: s.SpecValue,
				Price:     s.Price,
				Stock:     s.Stock,
			}
		}
		product.SetSKUs(skus)
	}

	if err := a.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	return product, nil
}

func (a *ProductApp) MarkReady(ctx context.Context, id int64) error {
	product, err := a.productRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}
	if err := product.MarkReady(); err != nil {
		return err
	}
	return a.productRepo.Update(ctx, product)
}
