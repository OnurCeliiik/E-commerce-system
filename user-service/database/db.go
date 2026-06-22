package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/OnurCeliiik/ecommerce/services/user/model"
	"github.com/OnurCeliiik/ecommerce/services/user/utils"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}

// SeedAdmin creates a bootstrap admin when ADMIN_EMAIL and ADMIN_PASSWORD are set.
func SeedAdmin(db *gorm.DB) error {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return nil
	}

	var existing model.User
	err := db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	admin := &model.User{
		ID:           uuid.New(),
		FirstName:    "Admin",
		LastName:     "User",
		Email:        email,
		PasswordHash: hash,
		Role:         model.RoleAdmin,
	}

	if err := db.Create(admin).Error; err != nil {
		return err
	}

	return nil
}
