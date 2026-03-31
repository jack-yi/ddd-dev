package internal

import (
	"github.com/yangboyi/ddd-dev/backend/internal/application"
	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	"github.com/yangboyi/ddd-dev/backend/internal/gateway"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
	repo "github.com/yangboyi/ddd-dev/backend/internal/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"gorm.io/gorm"
)

type Handlers struct {
	SourceItem *server.SourceItemHandler
	Product    *server.ProductHandler
	Publish    *server.PublishHandler
}

func InitHandlers(db *gorm.DB) *Handlers {
	// Repositories
	sourceItemRepo := repo.NewSourceItemRepoImpl(db)
	productRepo := repo.NewProductRepoImpl(db)
	publishTaskRepo := repo.NewPublishTaskRepoImpl(db)

	// Gateways (Mock)
	sourceGateway := gateway.NewMockSourceGateway()
	targetGateway := gateway.NewMockTargetGateway()

	// Queries
	sourceItemQuery := queries.NewSourceItemQuery(db)

	// Domain Services
	publishDomainService := domainservice.NewPublishDomainService(productRepo, publishTaskRepo, targetGateway)

	// Application Services
	sourceItemApp := application.NewSourceItemApp(sourceItemRepo, sourceGateway, sourceItemQuery)
	productApp := application.NewProductApp(productRepo, sourceItemRepo, db)
	publishApp := application.NewPublishApp(productRepo, publishDomainService, db)

	// Handlers
	return &Handlers{
		SourceItem: server.NewSourceItemHandler(sourceItemApp),
		Product:    server.NewProductHandler(productApp),
		Publish:    server.NewPublishHandler(publishApp),
	}
}
