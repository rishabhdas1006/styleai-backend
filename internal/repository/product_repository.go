package repository

import (
	"strings"
	"styleai-backend/internal/dto"
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

func (r *ProductRepository) FindAll(offset int, filter dto.GetProductsQuery) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	// Base query
	db := r.db.Model(&models.Product{}).
		Joins("LEFT JOIN product_variants v ON v.product_id = products.id").
		Preload("Category").
		Group("products.id")

	// product filters

	if filter.CreatedByID != nil {
		db = db.Where(
			"products.created_by_id = ?",
			*filter.CreatedByID,
		)
	}

	if filter.Category != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(categories.name) = ?", strings.ToLower(filter.Category))
	}

	if filter.Brand != "" {
		db = db.Where("LOWER(brand) = ?", strings.ToLower(filter.Brand))
	}

	if filter.Search != "" {
		filter.Search = "%" + strings.ToLower(filter.Search) + "%"
		db = db.Where("LOWER(name) LIKE ?", filter.Search)
	}

	if filter.Gender != "" {
		db = db.Where("LOWER(products.gender) = ?", strings.ToLower(filter.Gender))
	}

	// variant filters

	if filter.Color != "" {
		db = db.Where("v.color = ?", filter.Color)
	}

	if filter.Size != "" {
		db = db.Where("v.size = ?", filter.Size)
	}

	if filter.MinPrice > 0 {
		db = db.Where("v.price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		db = db.Where("v.price <= ?", filter.MaxPrice)
	}

	// count

	countDB := db.Session(&gorm.Session{}) // clone query

	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// sorting

	switch filter.Sort {
	case "price_asc":
		db = db.Order("products.min_price ASC")
	case "price_desc":
		db = db.Order("products.min_price DESC")
	default:
		db = db.Order("products.id DESC")
	}

	// pagination

	err := db.
		Limit(filter.Limit).
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
