package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/order/catalog"
	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/OnurCeliiik/ecommerce/services/order/model"
	"github.com/OnurCeliiik/ecommerce/services/order/repository"
	"github.com/OnurCeliiik/ecommerce/services/order/service"
	"github.com/OnurCeliiik/ecommerce/services/order/users"
	"github.com/google/uuid"
)

type mockOrderRepository struct {
	create       func(ctx context.Context, order *model.Order) error
	findByID     func(ctx context.Context, id uuid.UUID) (*model.Order, error)
	updateStatus func(ctx context.Context, id uuid.UUID, status string) error
	findByUserID func(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return m.create(ctx, order)
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	return m.findByID(ctx, id)
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.updateStatus(ctx, id, status)
}

func (m *mockOrderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return m.findByUserID(ctx, userID)
}

type mockProductCatalog struct {
	getUnitPrice func(ctx context.Context, productID uuid.UUID) (float64, error)
}

func (m *mockProductCatalog) GetUnitPrice(ctx context.Context, productID uuid.UUID) (float64, error) {
	return m.getUnitPrice(ctx, productID)
}

type mockUserDirectory struct {
	getUserEmail func(ctx context.Context, userID uuid.UUID) (string, error)
}

func (m *mockUserDirectory) GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	return m.getUserEmail(ctx, userID)
}

type mockOrderPublisher struct {
	publishOrderCreated func(ctx context.Context, event dto.OrderCreatedEvent) error
}

func (m *mockOrderPublisher) PublishOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error {
	return m.publishOrderCreated(ctx, event)
}

func TestCreateOrder_Success(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	var savedOrder *model.Order
	var publishedEvent dto.OrderCreatedEvent

	repo := &mockOrderRepository{
		create: func(_ context.Context, order *model.Order) error {
			savedOrder = order
			return nil
		},
	}
	catalogMock := &mockProductCatalog{
		getUnitPrice: func(_ context.Context, id uuid.UUID) (float64, error) {
			if id == productID {
				return 10, nil
			}
			return 0, catalog.ErrProductNotFound
		},
	}
	userDir := &mockUserDirectory{
		getUserEmail: func(_ context.Context, id uuid.UUID) (string, error) {
			if id == userID {
				return "buyer@example.com", nil
			}
			return "", users.ErrUserNotFound
		},
	}
	publisher := &mockOrderPublisher{
		publishOrderCreated: func(_ context.Context, event dto.OrderCreatedEvent) error {
			publishedEvent = event
			return nil
		},
	}

	svc := service.NewOrderService(repo, catalogMock, userDir, publisher)
	resp, err := svc.CreateOrder(context.Background(), userID, dto.CreateOrderRequest{
		Items: []dto.OrderLineRequest{{ProductID: productID, Quantity: 2}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != string(model.OrderStatusPending) {
		t.Fatalf("expected pending status, got %s", resp.Status)
	}
	if resp.Total != 20 {
		t.Fatalf("expected total 20, got %f", resp.Total)
	}
	if savedOrder == nil {
		t.Fatal("expected order to be saved")
	}
	if savedOrder.CustomerEmail != "buyer@example.com" {
		t.Fatalf("expected customer email on order, got %s", savedOrder.CustomerEmail)
	}
	if publishedEvent.OrderID != savedOrder.ID {
		t.Fatalf("expected published event for order %s, got %s", savedOrder.ID, publishedEvent.OrderID)
	}
	if publishedEvent.CustomerEmail != "buyer@example.com" {
		t.Fatalf("expected customer email in event, got %s", publishedEvent.CustomerEmail)
	}
}

func TestCreateOrder_UserNotFound(t *testing.T) {
	createCalled := false

	repo := &mockOrderRepository{
		create: func(_ context.Context, _ *model.Order) error {
			createCalled = true
			return nil
		},
	}
	svc := service.NewOrderService(
		repo,
		&mockProductCatalog{},
		&mockUserDirectory{
			getUserEmail: func(_ context.Context, _ uuid.UUID) (string, error) {
				return "", users.ErrUserNotFound
			},
		},
		&mockOrderPublisher{},
	)

	_, err := svc.CreateOrder(context.Background(), uuid.New(), dto.CreateOrderRequest{
		Items: []dto.OrderLineRequest{{ProductID: uuid.New(), Quantity: 1}},
	})
	if !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if createCalled {
		t.Fatal("expected Create not to be called when user not found")
	}
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	createCalled := false

	repo := &mockOrderRepository{
		create: func(_ context.Context, _ *model.Order) error {
			createCalled = true
			return nil
		},
	}
	svc := service.NewOrderService(
		repo,
		&mockProductCatalog{
			getUnitPrice: func(_ context.Context, _ uuid.UUID) (float64, error) {
				return 0, catalog.ErrProductNotFound
			},
		},
		&mockUserDirectory{
			getUserEmail: func(_ context.Context, _ uuid.UUID) (string, error) {
				return "buyer@example.com", nil
			},
		},
		&mockOrderPublisher{},
	)

	_, err := svc.CreateOrder(context.Background(), uuid.New(), dto.CreateOrderRequest{
		Items: []dto.OrderLineRequest{{ProductID: uuid.New(), Quantity: 1}},
	})
	if !errors.Is(err, service.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
	if createCalled {
		t.Fatal("expected Create not to be called when product not found")
	}
}

func TestGetOrder_Success(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, id uuid.UUID) (*model.Order, error) {
			if id != orderID {
				return nil, repository.ErrOrderNotFound
			}
			return &model.Order{
				ID:        orderID,
				UserID:    userID,
				Status:    string(model.OrderStatusPending),
				Total:     15,
				CreatedAt: time.Now(),
			}, nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	resp, err := svc.GetOrder(context.Background(), userID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != orderID {
		t.Fatalf("expected order %s, got %s", orderID, resp.ID)
	}
}

func TestGetOrder_WrongUser(t *testing.T) {
	orderID := uuid.New()

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, _ uuid.UUID) (*model.Order, error) {
			return &model.Order{
				ID:     orderID,
				UserID: uuid.New(),
				Status: string(model.OrderStatusPending),
			}, nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	_, err := svc.GetOrder(context.Background(), uuid.New(), orderID)
	if !errors.Is(err, service.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	repo := &mockOrderRepository{
		findByID: func(_ context.Context, _ uuid.UUID) (*model.Order, error) {
			return nil, repository.ErrOrderNotFound
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	_, err := svc.GetOrder(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, service.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestGetOrders_Success(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	repo := &mockOrderRepository{
		findByUserID: func(_ context.Context, id uuid.UUID) ([]*model.Order, error) {
			if id != userID {
				return nil, repository.ErrOrderNotFound
			}
			return []*model.Order{{
				ID:        orderID,
				UserID:    userID,
				Status:    string(model.OrderStatusConfirmed),
				Total:     42,
				CreatedAt: time.Now(),
			}}, nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	orders, err := svc.GetOrders(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) != 1 || orders[0].ID != orderID {
		t.Fatalf("expected one order %s, got %+v", orderID, orders)
	}
}

func TestProcessInventoryReserved_ConfirmsPendingOrder(t *testing.T) {
	orderID := uuid.New()
	var updatedStatus string

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, id uuid.UUID) (*model.Order, error) {
			return &model.Order{
				ID:     orderID,
				Status: string(model.OrderStatusPending),
			}, nil
		},
		updateStatus: func(_ context.Context, id uuid.UUID, status string) error {
			if id != orderID {
				t.Fatalf("expected order %s, got %s", orderID, id)
			}
			updatedStatus = status
			return nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{OrderID: orderID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedStatus != string(model.OrderStatusConfirmed) {
		t.Fatalf("expected confirmed, got %q", updatedStatus)
	}
}

func TestProcessInventoryReserved_SkipsTerminalOrder(t *testing.T) {
	orderID := uuid.New()
	updateCalled := false

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, _ uuid.UUID) (*model.Order, error) {
			return &model.Order{
				ID:     orderID,
				Status: string(model.OrderStatusConfirmed),
			}, nil
		},
		updateStatus: func(_ context.Context, _ uuid.UUID, _ string) error {
			updateCalled = true
			return nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{OrderID: orderID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updateCalled {
		t.Fatal("expected UpdateStatus not to be called for terminal order")
	}
}

func TestProcessInventoryReservationFailed_MarksFailed(t *testing.T) {
	orderID := uuid.New()
	var updatedStatus string

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, _ uuid.UUID) (*model.Order, error) {
			return &model.Order{
				ID:     orderID,
				Status: string(model.OrderStatusPending),
			}, nil
		},
		updateStatus: func(_ context.Context, _ uuid.UUID, status string) error {
			updatedStatus = status
			return nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	err := svc.ProcessInventoryReservationFailed(context.Background(), dto.InventoryReservationFailedEvent{
		OrderID: orderID,
		Reason:  "insufficient_inventory",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedStatus != string(model.OrderStatusFailed) {
		t.Fatalf("expected failed, got %q", updatedStatus)
	}
}

func TestProcessInventoryReservationFailed_SkipsTerminalOrder(t *testing.T) {
	orderID := uuid.New()
	updateCalled := false

	repo := &mockOrderRepository{
		findByID: func(_ context.Context, _ uuid.UUID) (*model.Order, error) {
			return &model.Order{
				ID:     orderID,
				Status: string(model.OrderStatusFailed),
			}, nil
		},
		updateStatus: func(_ context.Context, _ uuid.UUID, _ string) error {
			updateCalled = true
			return nil
		},
	}
	svc := service.NewOrderService(repo, &mockProductCatalog{}, &mockUserDirectory{}, &mockOrderPublisher{})

	err := svc.ProcessInventoryReservationFailed(context.Background(), dto.InventoryReservationFailedEvent{
		OrderID: orderID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updateCalled {
		t.Fatal("expected UpdateStatus not to be called for terminal order")
	}
}
