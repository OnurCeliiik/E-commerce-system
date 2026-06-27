package database

import (
	"fmt"
	"os"

	"github.com/OnurCeliiik/ecommerce/services/order/model"
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
	// Backfill-safe migration for customer_email on existing databases.
	if err := db.Exec(`ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_email text DEFAULT ''`).Error; err != nil {
		return err
	}
	if err := db.Exec(`UPDATE orders SET customer_email = '' WHERE customer_email IS NULL`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.Order{},
		&model.OrderLine{},
	)
}
