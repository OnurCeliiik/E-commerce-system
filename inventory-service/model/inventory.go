package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	ProductID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Quantity  int       `gorm:"not null;default:0"`
	UpdatedAt time.Time
}
