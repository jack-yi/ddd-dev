package platform

import (
	"context"

	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type SourcePlatformGateway interface {
	FetchProduct(ctx context.Context, sourceURL string) (*anticorruption.SourceProduct, error)
}

type TargetPlatformGateway interface {
	PublishProduct(ctx context.Context, name, description string, images []string,
		sellPrice float64, config anticorruption.PublishConfig) (*anticorruption.PublishResult, error)
}
