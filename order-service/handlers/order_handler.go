package handlers

import (
	"context"

	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This interface is used to define the methods that the order service must implement.
type OrderService interface {
	CreateOrder(ctx context.Context, userID uuid.UUID, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {}
