package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/order/catalog"
	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/OnurCeliiik/ecommerce/services/order/metrics"
	"github.com/OnurCeliiik/ecommerce/services/order/model"
	"github.com/OnurCeliiik/ecommerce/services/order/repository"
	"github.com/OnurCeliiik/ecommerce/services/order/users"
	"github.com/google/uuid"
)

// This interface is used to define the methods that the order service needs.
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
}

type ProductCatalog interface {
	GetUnitPrice(ctx context.Context, productID uuid.UUID) (float64, error)
}

type UserDirectory interface {
	GetUserEmail(ctx context.Context, userID uuid.UUID) (string, error)
}

type OrderEventPublisher interface {
	PublishOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error
}

type orderService struct {
	repo      OrderRepository
	catalog   ProductCatalog
	users     UserDirectory
	publisher OrderEventPublisher
}

func NewOrderService(
	repo OrderRepository,
	catalog ProductCatalog,
	users UserDirectory,
	publisher OrderEventPublisher) *orderService {
	return &orderService{
		repo:      repo,
		catalog:   catalog,
		users:     users,
		publisher: publisher,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	customerEmail, err := s.users.GetUserEmail(ctx, userID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	orderID := uuid.New()
	lines := make([]model.OrderLine, 0, len(req.Items))
	var total float64

	for _, item := range req.Items {
		unitPrice, err := s.catalog.GetUnitPrice(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, catalog.ErrProductNotFound) {
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		subtotal := unitPrice * float64(item.Quantity)
		total += subtotal

		lines = append(lines, model.OrderLine{
			ID:        uuid.New(),
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
		})
	}

	order := &model.Order{
		ID:            orderID,
		UserID:        userID,
		CustomerEmail: customerEmail,
		Status:        string(model.OrderStatusPending),
		Total:     total,
		CreatedAt: time.Now(),
		Lines:     lines,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := s.publisher.PublishOrderCreated(ctx, toOrderCreatedEvent(order)); err != nil {
		log.Printf("publish order.created failed for order %s: %v", order.ID, err)
	}

	return toOrderResponse(order), nil
}

func (s *orderService) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	return toOrderResponse(order), nil
}

func (s *orderService) ProcessInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {
	err := s.updateOrderStatusIfNotTerminal(ctx, event.OrderID, string(model.OrderStatusConfirmed))
	if err != nil {
		metrics.RecordInventoryEvent("reserved", "error")
		return err
	}
	metrics.RecordInventoryEvent("reserved", "success")
	return nil
}

func (s *orderService) ProcessInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error {
	err := s.updateOrderStatusIfNotTerminal(ctx, event.OrderID, string(model.OrderStatusFailed))
	if err != nil {
		metrics.RecordInventoryEvent("reservation_failed", "error")
		return err
	}
	metrics.RecordInventoryEvent("reservation_failed", "success")
	return nil
}

func (s *orderService) updateOrderStatusIfNotTerminal(ctx context.Context, orderID uuid.UUID, status string) error {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status == string(model.OrderStatusConfirmed) || order.Status == string(model.OrderStatusFailed) {
		return nil
	}

	return s.repo.UpdateStatus(ctx, orderID, status)
}

func toOrderCreatedEvent(order *model.Order) dto.OrderCreatedEvent {
	items := make([]dto.OrderLineResponse, 0, len(order.Lines))
	for i := range order.Lines {
		line := order.Lines[i]
		items = append(items, dto.OrderLineResponse{
			ID:        line.ID,
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
			UnitPrice: line.UnitPrice,
		})
	}

	return dto.OrderCreatedEvent{
		OrderID:       order.ID,
		UserID:        order.UserID,
		CustomerEmail: order.CustomerEmail,
		Total:         order.Total,
		Items:     items,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}
}

func toOrderResponse(order *model.Order) *dto.OrderResponse {
	items := make([]dto.OrderLineResponse, 0, len(order.Lines))
	for i := range order.Lines {
		line := order.Lines[i]
		items = append(items, dto.OrderLineResponse{
			ID:        line.ID,
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
			UnitPrice: line.UnitPrice,
		})
	}

	return &dto.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		Total:     order.Total,
		Items:     items,
		CreatedAt: order.CreatedAt,
	}
}

func (s *orderService) GetOrders(ctx context.Context, userID uuid.UUID) ([]*dto.OrderResponse, error) {
	orders, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	resp := make([]*dto.OrderResponse, 0, len(orders))
	for i := range orders {
		resp = append(resp, toOrderResponse(orders[i]))
	}

	return resp, nil
}
