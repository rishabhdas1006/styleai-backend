package repository

import (
	"strings"
	"styleai-backend/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Variants").
		Preload("Images").
		First(&product, id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll(
	offset, limit int,
	category, brand, search, sort string,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	db := r.db.Model(&models.Product{}).Preload("Category")

	// filtering

	if category != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(categories.name) = ?", strings.ToLower(category))
	}

	if brand != "" {
		db = db.Where("LOWER(brand) = ?", strings.ToLower(brand))
	}

	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
		db = db.Where("LOWER(name) LIKE ?", search)
	}

	// count

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// sorting

	switch sort {
	case "price_asc":
		db = db.Order("min_price ASC")
	case "price_desc":
		db = db.Order("min_price DESC")
	default:
		db = db.Order("id DESC")
	}

	// pagination

	err := db.
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *ProductRepository) UpdateMinPrice(productID uint) error {

	var minPrice float64

	err := r.db.
		Model(&models.ProductVariant{}).
		Where("product_id = ?", productID).
		Select("COALESCE(MIN(price), 0)").
		Scan(&minPrice).Error

	if err != nil {
		return err
	}

	return r.db.
		Model(&models.Product{}).
		Where("id = ?", productID).
		Update("min_price", minPrice).Error
}
