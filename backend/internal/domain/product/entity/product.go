package entity

import (
	"errors"
	"time"

	"github.com/yangboyi/ddd-dev/backend/infra/consts"
	sourceEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/source_item/entity"
)

type SKU struct {
	ID        int64
	SpecName  string
	SpecValue string
	Price     float64
	Stock     int
}

type Product struct {
	ID           int64
	SourceItemID int64
	Name         string
	Description  string
	Images       []string
	CostPrice    float64
	SellPrice    float64
	CategoryID   string
	Status       string
	SKUs         []SKU
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateFromSource(source *sourceEntity.SourceItem) *Product {
	now := time.Now()
	return &Product{
		SourceItemID: source.ID,
		Name:         source.Title,
		Description:  source.Description,
		Images:       source.Images,
		CostPrice:    source.Price.Min,
		SellPrice:    0,
		CategoryID:   source.Category,
		Status:       consts.ProductStatusDraft,
		SKUs:         []SKU{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (p *Product) EditInfo(name, description *string, images []string, categoryID *string) {
	if name != nil {
		p.Name = *name
	}
	if description != nil {
		p.Description = *description
	}
	if images != nil {
		p.Images = images
	}
	if categoryID != nil {
		p.CategoryID = *categoryID
	}
	p.UpdatedAt = time.Now()
}

func (p *Product) SetPrice(costPrice, sellPrice *float64) {
	if costPrice != nil {
		p.CostPrice = *costPrice
	}
	if sellPrice != nil {
		p.SellPrice = *sellPrice
	}
	p.UpdatedAt = time.Now()
}

func (p *Product) SetSKUs(skus []SKU) {
	p.SKUs = skus
	p.UpdatedAt = time.Now()
}

func (p *Product) MarkReady() error {
	if p.SellPrice <= 0 {
		return errors.New("sell price must be set before marking ready")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	p.Status = consts.ProductStatusReady
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) MarkPublished() {
	p.Status = consts.ProductStatusPublished
	p.UpdatedAt = time.Now()
}

func (p *Product) IsReady() bool {
	return p.Status == consts.ProductStatusReady
}
