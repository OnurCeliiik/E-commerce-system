package repository

import (
	"context"
	"strings"
	"time"

	"github.com/OnurCeliiik/ecommerce/services/notification/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type processedNotificationRepository struct {
	db *gorm.DB
}

func NewProcessedNotificationRepository(db *gorm.DB) *processedNotificationRepository {
	return &processedNotificationRepository{
		db: db,
	}
}

func (r *processedNotificationRepository) TryClaim(ctx context.Context, orderID uuid.UUID, eventType string) (bool, error) {
	record := model.ProcessedNotification{
		OrderID:     orderID,
		Outcome:     "processing",
		EventType:   eventType,
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
