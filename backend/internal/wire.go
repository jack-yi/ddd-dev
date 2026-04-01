package internal

import (
	"github.com/yangboyi/ddd-dev/backend/internal/application"
	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	"github.com/yangboyi/ddd-dev/backend/internal/gateway"
	"github.com/yangboyi/ddd-dev/backend/internal/middleware"
	"github.com/yangboyi/ddd-dev/backend/internal/queries"
	repo "github.com/yangboyi/ddd-dev/backend/internal/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/server"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type Handlers struct {
	SourceItem     *server.SourceItemHandler
	Product        *server.ProductHandler
	Publish        *server.PublishHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitHandlers(db *gorm.DB, userCenterRpc zrpc.Client) *Handlers {
	sourceItemRepo := repo.NewSourceItemRepoImpl(db)
	productRepo := repo.NewProductRepoImpl(db)
	publishTaskRepo := repo.NewPublishTaskRepoImpl(db)

	sourceGateway := gateway.NewMockSourceGateway()
	targetGateway := gateway.NewMockTargetGateway()

	sourceItemQuery := queries.NewSourceItemQuery(db)

	publishDomainService := domainservice.NewPublishDomainService(productRepo, publishTaskRepo, targetGateway)

	sourceItemApp := application.NewSourceItemApp(sourceItemRepo, sourceGateway, sourceItemQuery)
	productApp := application.NewProductApp(productRepo, sourceItemRepo, db)
	publishApp := application.NewPublishApp(productRepo, publishDomainService, db)

	authMw := middleware.NewAuthMiddleware(userCenterRpc)

	return &Handlers{
		SourceItem:     server.NewSourceItemHandler(sourceItemApp),
		Product:        server.NewProductHandler(productApp),
		Publish:        server.NewPublishHandler(publishApp),
		AuthMiddleware: authMw,
	}
}
