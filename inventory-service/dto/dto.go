package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

type InventoryResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}
