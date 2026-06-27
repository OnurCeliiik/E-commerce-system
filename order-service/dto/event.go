package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID       uuid.UUID           `json:"order_id"`
	UserID        uuid.UUID           `json:"user_id"`
	CustomerEmail string              `json:"customer_email"`
	Total         float64             `json:"total"`
	Items         []OrderLineResponse `json:"items"`
	Status        string              `json:"status"`
	CreatedAt     time.Time           `json:"created_at"`
}

type InventoryReservedEvent struct {
	OrderID       uuid.UUID           `json:"order_id"`
	UserID        uuid.UUID           `json:"user_id"`
	CustomerEmail string              `json:"customer_email"`
	Total         float64             `json:"total"`
	Items         []OrderLineResponse `json:"items"`
}

type InventoryReservationFailedEvent struct {
	OrderID       uuid.UUID `json:"order_id"`
	UserID        uuid.UUID `json:"user_id"`
	CustomerEmail string    `json:"customer_email"`
	Reason        string    `json:"reason"`
}
