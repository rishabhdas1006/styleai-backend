package repository

import (
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) WithTx(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{db: tx}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) CreateItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order

	err := r.db.
		Preload("Items").
		Where("user_id = ?", userID).
		Find(&orders).Error

	return orders, err
}

func (r *OrderRepository) FindByID(orderID uint) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("Items").
		First(&order, orderID).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}
