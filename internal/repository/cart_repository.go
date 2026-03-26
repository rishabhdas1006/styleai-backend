package repository

import (
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}

	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{UserID: userID}
		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}

	return nil, err
}

func (r *CartRepository) FindItem(cartID, variantID uint) (*models.CartItem, error) {
	var item models.CartItem

	result := r.db.
		Where("cart_id = ? AND variant_id = ?", cartID, variantID).
		First(&item)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &item, nil
}

func (r *CartRepository) GetItemByID(itemID uint) (*models.CartItem, error) {
	var item models.CartItem

	err := r.db.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *CartRepository) CreateItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *CartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *CartRepository) DeleteItem(itemID uint) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}

func (r *CartRepository) GetCartWithItems(userID uint) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.
		Preload("Items").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cart = models.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}

	return &cart, nil
}
