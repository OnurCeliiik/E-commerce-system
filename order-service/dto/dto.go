package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Items []OrderLineRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderLineRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type OrderLineResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}

type OrderResponse struct {
	ID        uuid.UUID           `json:"id"`
	UserID    uuid.UUID           `json:"user_id"`
	Status    string              `json:"status"`
	Total     float64             `json:"total"`
	Items     []OrderLineResponse `json:"items"`
	CreatedAt time.Time           `json:"created_at"`
}
