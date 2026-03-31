package gateway

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type MockTargetGateway struct{}

func NewMockTargetGateway() *MockTargetGateway {
	return &MockTargetGateway{}
}

func (g *MockTargetGateway) PublishProduct(ctx context.Context, name, description string,
	images []string, sellPrice float64, config anticorruption.PublishConfig) (*anticorruption.PublishResult, error) {
	return &anticorruption.PublishResult{
		PlatformProductID: fmt.Sprintf("pdd-mock-%d", rand.Intn(1000000)),
		Success:           true,
		ErrorMessage:      "",
	}, nil
}
