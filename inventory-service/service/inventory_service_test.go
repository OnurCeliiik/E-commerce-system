package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
	"github.com/OnurCeliiik/ecommerce/services/inventory/model"
	"github.com/OnurCeliiik/ecommerce/services/inventory/repository"
	"github.com/OnurCeliiik/ecommerce/services/inventory/service"
	"github.com/google/uuid"
)

type mockInventoryRepo struct {
	findByProductID func(ctx context.Context, productID uuid.UUID) (*model.InventoryItem, error)
	upsert          func(ctx context.Context, item *model.InventoryItem) error
}

func (m *mockInventoryRepo) FindByProductID(ctx context.Context, productID uuid.UUID) (*model.InventoryItem, error) {
	return m.findByProductID(ctx, productID)
}

func (m *mockInventoryRepo) Upsert(ctx context.Context, item *model.InventoryItem) error {
	return m.upsert(ctx, item)
}

type mockInventoryEventPublisher struct {
	publishInventoryReserved          func(ctx context.Context, event dto.InventoryReservedEvent) error
	publishInventoryReservationFailed func(ctx context.Context, event dto.InventoryReservationFailedEvent) error
}

func (m *mockInventoryEventPublisher) PublishInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {
	return m.publishInventoryReserved(ctx, event)
}

func (m *mockInventoryEventPublisher) PublishInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error {
	return m.publishInventoryReservationFailed(ctx, event)
}

type mockProcessedOrderRepo struct {
	tryClaim   func(ctx context.Context, orderID uuid.UUID) (bool, error)
	setOutcome func(ctx context.Context, orderID uuid.UUID, outcome string) error
}

func (m *mockProcessedOrderRepo) TryClaim(ctx context.Context, orderID uuid.UUID) (bool, error) {
	return m.tryClaim(ctx, orderID)
}

func (m *mockProcessedOrderRepo) SetOutcome(ctx context.Context, orderID uuid.UUID, outcome string) error {
	return m.setOutcome(ctx, orderID, outcome)
}

func TestGetInventory_Success(t *testing.T) {
	productID := uuid.New()
	item := &model.InventoryItem{
		ProductID: productID,
		Quantity:  100,
		UpdatedAt: time.Now(),
	}

	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, id uuid.UUID) (*model.InventoryItem, error) {
			if id == productID {
				return item, nil
			}
			return nil, repository.ErrInventoryNotFound
		},
		upsert: func(_ context.Context, _ *model.InventoryItem) error {
			return nil
		},
	}

	publisher := &mockInventoryEventPublisher{
		publishInventoryReserved: func(_ context.Context, _ dto.InventoryReservedEvent) error {
			return nil
		},
		publishInventoryReservationFailed: func(_ context.Context, _ dto.InventoryReservationFailedEvent) error {
			return nil
		},
	}

	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}

	svc := service.NewInventoryService(repo, processedRepo, publisher)

	resp, err := svc.GetInventory(context.Background(), productID)
	if err != nil {
		t.Fatalf("got error getting inventory: %v", err)
	}
	if resp.ProductID != productID {
		t.Fatalf("expected product ID to be %s, got %s", productID, resp.ProductID)
	}
	if resp.Quantity != 100 {
		t.Fatalf("expected quantity to be 100, got %d", resp.Quantity)
	}
}

func TestGetInventory_NotFound(t *testing.T) {
	productID := uuid.New()
	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, _ uuid.UUID) (*model.InventoryItem, error) {
			return nil, repository.ErrInventoryNotFound
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReserved: func(_ context.Context, _ dto.InventoryReservedEvent) error {
			return nil
		},
		publishInventoryReservationFailed: func(_ context.Context, _ dto.InventoryReservationFailedEvent) error {
			return nil
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	svc := service.NewInventoryService(repo, processedRepo, publisher)
	_, err := svc.GetInventory(context.Background(), productID)
	if !errors.Is(err, service.ErrInventoryNotFound) {
		t.Fatalf("expected ErrInventoryNotFound, got %v", err)
	}
}

func TestUpdateInventory_Success(t *testing.T) {
	productID := uuid.New()
	quantity := 100
	item := &model.InventoryItem{
		ProductID: productID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
	}
	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, id uuid.UUID) (*model.InventoryItem, error) {
			if id == productID {
				return item, nil
			}
			return nil, repository.ErrInventoryNotFound
		},
		upsert: func(_ context.Context, item *model.InventoryItem) error {
			if item.ProductID == productID && item.Quantity == quantity {
				return nil
			}
			return errors.New("quantity mismatch")
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReserved: func(_ context.Context, _ dto.InventoryReservedEvent) error {
			return nil
		},
		publishInventoryReservationFailed: func(_ context.Context, _ dto.InventoryReservationFailedEvent) error {
			return nil
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	svc := service.NewInventoryService(repo, processedRepo, publisher)
	resp, err := svc.UpdateInventory(context.Background(), productID, dto.UpdateInventoryRequest{
		Quantity: quantity,
	})
	if err != nil {
		t.Fatalf("got error updating inventory: %v", err)
	}
	if resp.ProductID != productID {
		t.Fatalf("expected product ID to be %s, got %s", productID, resp.ProductID)
	}
	if resp.Quantity != quantity {
		t.Fatalf("expected quantity to be %d, got %d", quantity, resp.Quantity)
	}
}

func TestUpdateInventory_InvalidQuantity(t *testing.T) {
	productID := uuid.New()
	quantity := -1
	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, _ uuid.UUID) (*model.InventoryItem, error) {
			return nil, repository.ErrInventoryNotFound
		},
	}
	publisher := &mockInventoryEventPublisher{}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	svc := service.NewInventoryService(repo, processedRepo, publisher)
	_, err := svc.UpdateInventory(context.Background(), productID, dto.UpdateInventoryRequest{Quantity: quantity})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestProcessOrderCreated_ReservesAndPublishes(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	stock := map[uuid.UUID]int{productID: 10}
	var reservedEvent dto.InventoryReservedEvent
	var outcome string

	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, id uuid.UUID) (*model.InventoryItem, error) {
			qty, ok := stock[id]
			if !ok {
				return nil, repository.ErrInventoryNotFound
			}
			return &model.InventoryItem{ProductID: id, Quantity: qty, UpdatedAt: time.Now()}, nil
		},
		upsert: func(_ context.Context, item *model.InventoryItem) error {
			stock[item.ProductID] = item.Quantity
			return nil
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, id uuid.UUID) (bool, error) {
			if id != orderID {
				t.Fatalf("expected order %s, got %s", orderID, id)
			}
			return true, nil
		},
		setOutcome: func(_ context.Context, id uuid.UUID, o string) error {
			if id != orderID {
				t.Fatalf("expected order %s, got %s", orderID, id)
			}
			outcome = o
			return nil
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReserved: func(_ context.Context, event dto.InventoryReservedEvent) error {
			reservedEvent = event
			return nil
		},
		publishInventoryReservationFailed: func(_ context.Context, _ dto.InventoryReservationFailedEvent) error {
			t.Fatal("PublishInventoryReservationFailed should not be called")
			return nil
		},
	}

	svc := service.NewInventoryService(repo, processedRepo, publisher)
	event := dto.OrderCreatedEvent{
		OrderID:       orderID,
		UserID:        userID,
		CustomerEmail: "buyer@example.com",
		Total:         20,
		Items: []dto.OrderLineItem{
			{ProductID: productID, Quantity: 2, UnitPrice: 10},
		},
	}

	err := svc.ProcessOrderCreated(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stock[productID] != 8 {
		t.Fatalf("expected stock 8, got %d", stock[productID])
	}
	if outcome != "reserved" {
		t.Fatalf("expected outcome reserved, got %q", outcome)
	}
	if reservedEvent.OrderID != orderID {
		t.Fatalf("expected reserved event for order %s, got %s", orderID, reservedEvent.OrderID)
	}
}

func TestProcessOrderCreated_Duplicate(t *testing.T) {
	findCalled := false
	publishCalled := false

	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, _ uuid.UUID) (*model.InventoryItem, error) {
			findCalled = true
			return nil, nil
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return false, nil
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReserved: func(_ context.Context, _ dto.InventoryReservedEvent) error {
			publishCalled = true
			return nil
		},
	}

	svc := service.NewInventoryService(repo, processedRepo, publisher)
	err := svc.ProcessOrderCreated(context.Background(), dto.OrderCreatedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
	})
	if err != nil {
		t.Fatalf("expected no error on duplicate, got %v", err)
	}
	if findCalled {
		t.Fatal("expected FindByProductID not to be called for duplicate")
	}
	if publishCalled {
		t.Fatal("expected PublishInventoryReserved not to be called for duplicate")
	}
}

func TestProcessOrderCreated_InsufficientStock(t *testing.T) {
	productID := uuid.New()
	var failedEvent dto.InventoryReservationFailedEvent
	var outcome string

	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, id uuid.UUID) (*model.InventoryItem, error) {
			return &model.InventoryItem{ProductID: id, Quantity: 1, UpdatedAt: time.Now()}, nil
		},
		upsert: func(_ context.Context, _ *model.InventoryItem) error {
			t.Fatal("Upsert should not be called when stock is insufficient")
			return nil
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, o string) error {
			outcome = o
			return nil
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReservationFailed: func(_ context.Context, event dto.InventoryReservationFailedEvent) error {
			failedEvent = event
			return nil
		},
	}

	svc := service.NewInventoryService(repo, processedRepo, publisher)
	err := svc.ProcessOrderCreated(context.Background(), dto.OrderCreatedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
		Items:   []dto.OrderLineItem{{ProductID: productID, Quantity: 5}},
	})
	if !errors.Is(err, service.ErrInsufficientInventory) {
		t.Fatalf("expected ErrInsufficientInventory, got %v", err)
	}
	if outcome != "failed" {
		t.Fatalf("expected outcome failed, got %q", outcome)
	}
	if failedEvent.Reason != "insufficient_inventory" {
		t.Fatalf("expected insufficient_inventory reason, got %q", failedEvent.Reason)
	}
}

func TestProcessOrderCreated_InventoryNotFound(t *testing.T) {
	productID := uuid.New()
	var failedEvent dto.InventoryReservationFailedEvent

	repo := &mockInventoryRepo{
		findByProductID: func(_ context.Context, _ uuid.UUID) (*model.InventoryItem, error) {
			return nil, repository.ErrInventoryNotFound
		},
	}
	processedRepo := &mockProcessedOrderRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
		setOutcome: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	publisher := &mockInventoryEventPublisher{
		publishInventoryReservationFailed: func(_ context.Context, event dto.InventoryReservationFailedEvent) error {
			failedEvent = event
			return nil
		},
	}

	svc := service.NewInventoryService(repo, processedRepo, publisher)
	err := svc.ProcessOrderCreated(context.Background(), dto.OrderCreatedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
		Items:   []dto.OrderLineItem{{ProductID: productID, Quantity: 1}},
	})
	if !errors.Is(err, service.ErrInventoryNotFound) {
		t.Fatalf("expected ErrInventoryNotFound, got %v", err)
	}
	if failedEvent.Reason != "inventory_not_found" {
		t.Fatalf("expected inventory_not_found reason, got %q", failedEvent.Reason)
	}
}
