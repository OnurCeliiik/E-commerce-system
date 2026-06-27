package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/order/catalog"
	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/OnurCeliiik/ecommerce/services/order/model"
	"github.com/google/uuid"
)

// This interface is used to define the methods that the order service needs.
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type ProductCatalog interface {
	GetUnitPrice(ctx context.Context, productID uuid.UUID) (float64, error)
}

type OrderEventPublisher interface {
	PublishOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error
}

type orderService struct {
	repo      OrderRepository
	catalog   ProductCatalog
	publisher OrderEventPublisher
}

func NewOrderService(
	repo OrderRepository,
	catalog ProductCatalog,
	publisher OrderEventPublisher) *orderService {
	return &orderService{
		repo:      repo,
		catalog:   catalog,
		publisher: publisher,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {

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
		ID:        orderID,
		UserID:    userID,
		Status:    string(model.OrderStatusPending),
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

func (s *orderService) ProcessInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {
	return s.repo.UpdateStatus(ctx, event.OrderID, string(model.OrderStatusConfirmed))
}

func (s *orderService) ProcessInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error {
	return s.repo.UpdateStatus(ctx, event.OrderID, string(model.OrderStatusFailed))
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
		OrderID:   order.ID,
		UserID:    order.UserID,
		Total:     order.Total,
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
