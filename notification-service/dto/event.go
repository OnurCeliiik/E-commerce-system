package dto

import "github.com/google/uuid"

type InventoryReservedEvent struct {
	OrderID uuid.UUID       `json:"order_id"`
	UserID  uuid.UUID       `json:"user_id"`
	Total   float64         `json:"total"`
	Items   []OrderLineItem `json:"items"`
}

type InventoryReservationFailedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
	Reason  string    `json:"reason"`
}

type OrderLineItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
}
