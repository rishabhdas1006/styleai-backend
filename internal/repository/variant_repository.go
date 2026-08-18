package repository

import (
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type VariantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

func (r *VariantRepository) WithTx(tx *gorm.DB) *VariantRepository {
	return &VariantRepository{db: tx}
}

func (r *VariantRepository) Create(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *VariantRepository) FindByProductID(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant

	err := r.db.
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("product_id = ?", productID).
		Find(&variants).Error

	return variants, err
}

func (r *VariantRepository) FindByID(id string) (*models.ProductVariant, error) {
	var variant models.ProductVariant

	err := r.db.First(&variant, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrVariantNotFound
		}
		return nil, err
	}

	return &variant, nil
}

func (r *VariantRepository) Update(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *VariantRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductVariant{}, "id = ?", id).Error
}

func (r *VariantRepository) ReduceStock(variantID string, qty int) error {
	result := r.db.Model(&models.ProductVariant{}).
		Where("id = ? AND stock >= ?", variantID, qty).
		Update("stock", gorm.Expr("stock - ?", qty))

	if result.RowsAffected == 0 {
		return common.ErrInsufficientStock
	}

	return result.Error
}
