package service

import (
	"context"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/OnurCeliiik/ecommerce/services/order/model"
	"github.com/google/uuid"
)

// This interface is used to define the methods that the order service needs.
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
}

type ProductCatalog interface {
	GetUnitPrice(ctx context.Context, productID uuid.UUID) (float64, error)
}

type orderService struct {
	repo    OrderRepository
	catalog ProductCatalog
}

func NewOrderService(repo OrderRepository, catalog ProductCatalog) *orderService {
	return &orderService{
		repo:    repo,
		catalog: catalog,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {

	orderID := uuid.New()
	lines := make([]model.OrderLine, 0, len(req.Items))
	var total float64

	for _, item := range req.Items {
		unitPrice, err := s.catalog.GetUnitPrice(ctx, item.ProductID)
		if err != nil {
			return nil, ErrProductNotFound
		}

		subtotal := unitPrice * float64(item.Quantity)
		total += subtotal

		lines = append(lines, model.OrderLine{
			ID: uuid.New(),
			OrderID: orderID,
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			UnitPrice: unitPrice,
		})
	}

	order := &model.Order{
		ID: orderID,
		UserID: userID,
		Status: string(model.OrderStatusPending),
		Total: total,
		CreatedAt: time.Now(),
		Lines: lines,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return toOrderResponse(order), nil
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
