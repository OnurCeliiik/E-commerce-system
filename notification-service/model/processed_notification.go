package model

import (
	"time"

	"github.com/google/uuid"
)

type ProcessedNotification struct {
	OrderID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	EventType   string    `gorm:"primaryKey"`
	Outcome     string
	ProcessedAt time.Time
}
