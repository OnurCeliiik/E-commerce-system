package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

var (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	Status    string      `gorm:"not null" json:"status"`
	Total     float64     `gorm:"not null" json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Lines     []OrderLine `gorm:"foreignKey:OrderID"`
}

type OrderLine struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"not null" json:"unit_price"`
}
