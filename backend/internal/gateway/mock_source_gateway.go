package gateway

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type MockSourceGateway struct{}

func NewMockSourceGateway() *MockSourceGateway {
	return &MockSourceGateway{}
}

func (g *MockSourceGateway) FetchProduct(ctx context.Context, sourceURL string) (*anticorruption.SourceProduct, error) {
	return &anticorruption.SourceProduct{
		ExternalID:  "mock-ext-001",
		Title:       "Mock商品 - 高品质T恤",
		Description: "优质纯棉T恤，多色可选，批发价格优惠",
		Images:      []string{"https://via.placeholder.com/800x800?text=Mock+Image+1", "https://via.placeholder.com/800x800?text=Mock+Image+2"},
		PriceMin:    15.00,
		PriceMax:    25.00,
		Supplier: anticorruption.Supplier{
			Name:   "广州优品服饰有限公司",
			Rating: 4.8,
			Region: "广东广州",
		},
		Category:    "服装/T恤",
		SalesVolume: 10000,
		MinOrder:    2,
	}, nil
}
