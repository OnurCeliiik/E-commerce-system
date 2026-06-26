package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID   uuid.UUID           `json:"order_id"`
	UserID    uuid.UUID           `json:"user_id"`
	Total     float64             `json:"total"`
	Items     []OrderLineResponse `json:"items"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}

type OrderLineResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}
