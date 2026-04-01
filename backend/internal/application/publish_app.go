package application

import (
	"context"
	"fmt"

	domainservice "github.com/yangboyi/ddd-dev/backend/internal/domain/domain_service"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	publishEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	po "github.com/yangboyi/ddd-dev/backend/internal/model/po/mysql"
	"gorm.io/gorm"
)

type PublishApp struct {
	productRepo    productRepo.ProductRepository
	publishService *domainservice.PublishDomainService
	db             *gorm.DB
}

func NewPublishApp(
	pr productRepo.ProductRepository,
	ps *domainservice.PublishDomainService,
	db *gorm.DB,
) *PublishApp {
	return &PublishApp{productRepo: pr, publishService: ps, db: db}
}

func (a *PublishApp) CreateTask(ctx context.Context, req *dto.CreatePublishTaskReq) (*publishEntity.PublishTask, error) {
	product, err := a.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}

	config := publishEntity.PublishConfig{
		CategoryID:      req.CategoryID,
		FreightTemplate: req.FreightTemplate,
	}

	task, err := a.publishService.PublishProduct(ctx, product, req.TargetPlatform, config)
	if err != nil {
		return nil, fmt.Errorf("publish product: %w", err)
	}
	return task, nil
}

func (a *PublishApp) List(ctx context.Context, filter *dto.PublishTaskFilter) ([]po.PublishTask, int64, error) {
	query := a.db.WithContext(ctx).Model(&po.PublishTask{})

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.TargetPlatform != nil {
		query = query.Where("target_platform = ?", *filter.TargetPlatform)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count publish tasks: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	var items []po.PublishTask
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list publish tasks: %w", err)
	}

	return items, total, nil
}
