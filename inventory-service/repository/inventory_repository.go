package repository

import (
	"context"
	"errors"

	"github.com/OnurCeliiik/ecommerce/services/inventory/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInventoryNotFound = errors.New("inventory not found")

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*model.InventoryItem, error) {
	var item model.InventoryItem

	err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInventoryNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *inventoryRepository) Upsert(ctx context.Context, item *model.InventoryItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}
