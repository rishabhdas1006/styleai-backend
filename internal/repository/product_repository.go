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

func (r *ProductRepository) WithTx(tx *gorm.DB) *ProductRepository {
	return &ProductRepository{db: tx}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Images").
		First(&product, id).Error

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll(
	offset, limit int,
	category, brand, search, sort, gender string,
	color, size string,
	minPrice, maxPrice float64,
) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	// Base query
	db := r.db.Model(&models.Product{}).
		Joins("LEFT JOIN product_variants v ON v.product_id = products.id").
		Preload("Category").
		Group("products.id")

	// product filters

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

	if gender != "" {
		db = db.Where("LOWER(products.gender) = ?", strings.ToLower(gender))
	}

	// variant filters

	if color != "" {
		db = db.Where("v.color = ?", color)
	}

	if size != "" {
		db = db.Where("v.size = ?", size)
	}

	if minPrice > 0 {
		db = db.Where("v.price >= ?", minPrice)
	}

	if maxPrice > 0 {
		db = db.Where("v.price <= ?", maxPrice)
	}

	// count

	countDB := db.Session(&gorm.Session{}) // clone query

	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// sorting

	switch sort {
	case "price_asc":
		db = db.Order("products.min_price ASC")
	case "price_desc":
		db = db.Order("products.min_price DESC")
	default:
		db = db.Order("products.id DESC")
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
