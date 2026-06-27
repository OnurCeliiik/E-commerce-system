package repository

import (
	"context"
	"errors"

	"github.com/OnurCeliiik/ecommerce/services/order/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var order model.Order

	err := r.db.WithContext(ctx).Preload("Lines").Where("id = ?", id).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	var orders []*model.Order

	err := r.db.WithContext(ctx).Preload("Lines").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
