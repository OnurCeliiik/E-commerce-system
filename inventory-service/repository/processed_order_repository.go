package repository

import (
	"context"
	"strings"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/inventory/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type processedOrderRepository struct {
	db *gorm.DB
}

func NewProcessedOrderRepository(db *gorm.DB) *processedOrderRepository {
	return &processedOrderRepository{db: db}
}

func (r *processedOrderRepository) TryClaim(ctx context.Context, orderID uuid.UUID) (bool, error) {
	record := model.ProcessedOrder{
		OrderID:     orderID,
		Outcome:     "processing",
		ProcessedAt: time.Now(),
	}
	err := r.db.WithContext(ctx).Create(&record).Error
	if err != nil {
		if isUniqueViolation(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

func (r *processedOrderRepository) SetOutcome(ctx context.Context, orderID uuid.UUID, outcome string) error {
	return r.db.WithContext(ctx).Model(&model.ProcessedOrder{}).
		Where("order_id = ?", orderID).
		Update("outcome", outcome).Error
}
