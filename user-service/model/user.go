package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName    string    `gorm:"not null"`
	LastName     string    `gorm:"not null"`
	Email        string    `gorm:"not null;uniqueIndex"`
	PasswordHash string    `gorm:"not null;column:password_hash"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
