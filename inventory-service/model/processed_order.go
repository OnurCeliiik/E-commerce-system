package model

import (
	"time"

	"github.com/google/uuid"
)

type ProcessedOrder struct {
	OrderID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Outcome     string    `gorm:"not null"` // "reserved" or "failed"
	ProcessedAt time.Time
}
