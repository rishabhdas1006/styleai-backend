package repository

import (
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type VariantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

func (r *VariantRepository) Create(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *VariantRepository) FindByProductID(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := r.db.Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

func (r *VariantRepository) FindByID(id uint) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.First(&variant, id).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *VariantRepository) Update(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *VariantRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductVariant{}, id).Error
}
