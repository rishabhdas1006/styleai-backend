package repository

import (
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type VariantImageRepository struct {
	db *gorm.DB
}

func NewVariantImageRepository(db *gorm.DB) *VariantImageRepository {
	return &VariantImageRepository{db: db}
}

func (r *VariantImageRepository) WithTx(tx *gorm.DB) *VariantImageRepository {
	return &VariantImageRepository{db: tx}
}

func (r *VariantImageRepository) CreateImages(images []models.VariantImage) error {
	return r.db.Create(&images).Error
}

func (r *VariantImageRepository) FindByVariantID(variantID string) ([]models.VariantImage, error) {
	var images []models.VariantImage
	err := r.db.Where("variant_id = ?", variantID).Order("position ASC").Find(&images).Error
	return images, err
}

func (r *VariantImageRepository) DeleteByVariantID(variantID string) error {
	return r.db.Where("variant_id = ?", variantID).Delete(&models.VariantImage{}).Error
}
