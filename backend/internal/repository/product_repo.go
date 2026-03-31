package repository

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type ProductRepoImpl struct {
	db *gorm.DB
}

func NewProductRepoImpl(db *gorm.DB) *ProductRepoImpl {
	return &ProductRepoImpl{db: db}
}

func (r *ProductRepoImpl) Save(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := toProductPO(product)
		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("create product: %w", err)
		}
		product.ID = record.ID

		for i := range product.SKUs {
			sku := &po.ProductSKU{
				ProductID: product.ID,
				SpecName:  product.SKUs[i].SpecName,
				SpecValue: product.SKUs[i].SpecValue,
				Price:     product.SKUs[i].Price,
				Stock:     product.SKUs[i].Stock,
			}
			if err := tx.Create(sku).Error; err != nil {
				return fmt.Errorf("create product sku: %w", err)
			}
			product.SKUs[i].ID = sku.ID
		}
		return nil
	})
}

func (r *ProductRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	var record po.Product
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}

	var skuRecords []po.ProductSKU
	if err := r.db.WithContext(ctx).Where("product_id = ?", id).Find(&skuRecords).Error; err != nil {
		return nil, fmt.Errorf("find product skus: %w", err)
	}

	product := toProductEntity(&record)
	for _, s := range skuRecords {
		product.SKUs = append(product.SKUs, entity.SKU{
			ID:        s.ID,
			SpecName:  s.SpecName,
			SpecValue: s.SpecValue,
			Price:     s.Price,
			Stock:     s.Stock,
		})
	}
	return product, nil
}

func (r *ProductRepoImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		record := toProductPO(product)
		record.ID = product.ID
		if err := tx.Save(record).Error; err != nil {
			return fmt.Errorf("update product: %w", err)
		}

		if err := tx.Where("product_id = ?", product.ID).Delete(&po.ProductSKU{}).Error; err != nil {
			return fmt.Errorf("delete old skus: %w", err)
		}
		for i := range product.SKUs {
			sku := &po.ProductSKU{
				ProductID: product.ID,
				SpecName:  product.SKUs[i].SpecName,
				SpecValue: product.SKUs[i].SpecValue,
				Price:     product.SKUs[i].Price,
				Stock:     product.SKUs[i].Stock,
			}
			if err := tx.Create(sku).Error; err != nil {
				return fmt.Errorf("create product sku: %w", err)
			}
			product.SKUs[i].ID = sku.ID
		}
		return nil
	})
}

func toProductPO(e *entity.Product) *po.Product {
	return &po.Product{
		ID:           e.ID,
		SourceItemID: e.SourceItemID,
		Name:         e.Name,
		Description:  e.Description,
		Images:       po.StringSlice(e.Images),
		CostPrice:    e.CostPrice,
		SellPrice:    e.SellPrice,
		CategoryID:   e.CategoryID,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func toProductEntity(p *po.Product) *entity.Product {
	return &entity.Product{
		ID:           p.ID,
		SourceItemID: p.SourceItemID,
		Name:         p.Name,
		Description:  p.Description,
		Images:       []string(p.Images),
		CostPrice:    p.CostPrice,
		SellPrice:    p.SellPrice,
		CategoryID:   p.CategoryID,
		Status:       p.Status,
		SKUs:         []entity.SKU{},
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
