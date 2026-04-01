package domainservice

import (
	"context"
	"fmt"

	"github.com/yangboyi/ddd-dev/backend/internal/domain/platform"
	productEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/product/entity"
	productRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/product/repository"
	publishEntity "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/entity"
	publishRepo "github.com/yangboyi/ddd-dev/backend/internal/domain/publish/repository"
	"github.com/yangboyi/ddd-dev/backend/internal/model/anticorruption"
)

type PublishDomainService struct {
	productRepo   productRepo.ProductRepository
	publishRepo   publishRepo.PublishTaskRepository
	targetGateway platform.TargetPlatformGateway
}

func NewPublishDomainService(
	pr productRepo.ProductRepository,
	ptr publishRepo.PublishTaskRepository,
	tg platform.TargetPlatformGateway,
) *PublishDomainService {
	return &PublishDomainService{
		productRepo:   pr,
		publishRepo:   ptr,
		targetGateway: tg,
	}
}

func (s *PublishDomainService) PublishProduct(ctx context.Context, product *productEntity.Product,
	targetPlatform string, config publishEntity.PublishConfig) (*publishEntity.PublishTask, error) {

	if !product.IsReady() {
		return nil, fmt.Errorf("product %d is not ready for publishing", product.ID)
	}

	task := publishEntity.NewPublishTask(product.ID, targetPlatform, config)
	if err := s.publishRepo.Save(ctx, task); err != nil {
		return nil, fmt.Errorf("save publish task: %w", err)
	}

	if err := task.MarkPublishing(); err != nil {
		return nil, fmt.Errorf("mark publishing: %w", err)
	}

	acConfig := anticorruption.PublishConfig{
		CategoryID:      config.CategoryID,
		FreightTemplate: config.FreightTemplate,
	}
	result, err := s.targetGateway.PublishProduct(ctx, product.Name, product.Description,
		product.Images, product.SellPrice, acConfig)
	if err != nil {
		task.MarkFailed(err.Error())
		_ = s.publishRepo.Update(ctx, task)
		return task, nil
	}

	if result.Success {
		task.MarkSuccess(result.PlatformProductID)
		product.MarkPublished()
		_ = s.productRepo.Update(ctx, product)
	} else {
		task.MarkFailed(result.ErrorMessage)
	}

	_ = s.publishRepo.Update(ctx, task)
	return task, nil
}
