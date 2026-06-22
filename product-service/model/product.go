package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	Price       float64   `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
