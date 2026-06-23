package service

import (
	"context"
	"errors"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
	"github.com/OnurCeliiik/ecommerce/services/inventory/model"
	"github.com/OnurCeliiik/ecommerce/services/inventory/repository"
	"github.com/google/uuid"
)

type InventoryRepository interface {
	FindByProductID(ctx context.Context, productID uuid.UUID) (*model.InventoryItem, error)
	Upsert(ctx context.Context, item *model.InventoryItem) error
}

type inventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *inventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) GetInventory(ctx context.Context, productID uuid.UUID) (*dto.InventoryResponse, error) {
	item, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrInventoryNotFound) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}

	return toInventoryResponse(item), nil
}

func (s *inventoryService) UpdateInventory(ctx context.Context, productID uuid.UUID, req dto.UpdateInventoryRequest) (*dto.InventoryResponse, error) {
	if req.Quantity < 0 {
		return nil, ErrInvalidInput
	}

	now := time.Now()
	item := &model.InventoryItem{
		ProductID: productID,
		Quantity:  req.Quantity,
		UpdatedAt: now,
	}
	if err := s.repo.Upsert(ctx, item); err != nil {
		return nil, err
	}

	return toInventoryResponse(item), nil
}

func toInventoryResponse(item *model.InventoryItem) *dto.InventoryResponse {
	return &dto.InventoryResponse{
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UpdatedAt: item.UpdatedAt,
	}
}
