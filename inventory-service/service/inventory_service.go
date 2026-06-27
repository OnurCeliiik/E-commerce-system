package service

import (
	"context"
	"errors"
	"log"
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

type InventoryEventPublisher interface {
	PublishInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error
	PublishInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error
}

type inventoryService struct {
	repo      InventoryRepository
	publisher InventoryEventPublisher
}

func NewInventoryService(repo InventoryRepository, publisher InventoryEventPublisher) *inventoryService {
	return &inventoryService{repo: repo, publisher: publisher}
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

func (s *inventoryService) ProcessOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error {
	if err := s.reserveStock(ctx, event); err != nil {
		failEvent := dto.InventoryReservationFailedEvent{
			OrderID: event.OrderID,
			UserID:  event.UserID,
			Reason:  failureReason(err),
		}
		if pubErr := s.publisher.PublishInventoryReservationFailed(ctx, failEvent); pubErr != nil {
			log.Printf("publish inventory.reservation_failed failed for order %s: %v", event.OrderID, pubErr)
		}
		return err
	}

	reservedEvent := dto.InventoryReservedEvent{
		OrderID: event.OrderID,
		UserID:  event.UserID,
		Total:   event.Total,
		Items:   event.Items,
	}
	if err := s.publisher.PublishInventoryReserved(ctx, reservedEvent); err != nil {
		log.Printf("publish inventory.reserved failed for order %s: %v", event.OrderID, err)
	}

	return nil
}

func (s *inventoryService) reserveStock(ctx context.Context, event dto.OrderCreatedEvent) error {
	for _, item := range event.Items {
		inventoryItem, err := s.repo.FindByProductID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, repository.ErrInventoryNotFound) {
				return ErrInventoryNotFound
			}
			return err
		}

		if inventoryItem.Quantity < item.Quantity {
			return ErrInsufficientInventory
		}

		inventoryItem.Quantity -= item.Quantity

		if err := s.repo.Upsert(ctx, inventoryItem); err != nil {
			return err
		}
	}

	return nil
}

func failureReason(err error) string {
	switch {
	case errors.Is(err, ErrInventoryNotFound):
		return "inventory_not_found"
	case errors.Is(err, ErrInsufficientInventory):
		return "insufficient_inventory"
	default:
		return "unknown"
	}
}
