package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName    string    `gorm:"not null"`
	LastName     string    `gorm:"not null"`
	Email        string    `gorm:"not null;uniqueIndex"`
	PasswordHash string    `gorm:"not null;column:password_hash"`
	Role         Role      `gorm:"not null;default:customer"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
