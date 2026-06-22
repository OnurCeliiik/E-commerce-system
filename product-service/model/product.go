package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
