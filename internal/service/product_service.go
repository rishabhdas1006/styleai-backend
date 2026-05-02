package service

import (
	"errors"
	"math"
	"strings"
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"

	"gorm.io/gorm"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: repo}
}

type ProductListResponse struct {
	Products   []models.Product `json:"products"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int64            `json:"totalPages"`
}

func (s *ProductService) CreateProduct(name, description, brand string, categoryID uint) (*models.Product, error) {

	product := &models.Product{
		Name:        name,
		Description: description,
		Brand:       brand,
		CategoryID:  categoryID,
	}

	err := s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProducts(
	page, limit int,
	category, brand, search, sort string,
	color, size string,
	minPrice, maxPrice float64,
) (*ProductListResponse, error) {

	if page < 1 {
		page = 1
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	category = strings.TrimSpace(category)
	brand = strings.TrimSpace(brand)
	search = strings.TrimSpace(search)
	color = strings.TrimSpace(color)
	size = strings.TrimSpace(size)
	sort = strings.TrimSpace(sort)

	if minPrice < 0 {
		minPrice = 0
	}

	if maxPrice > 0 && maxPrice < minPrice {
		return nil, common.ErrMaxPriceLessThanMinPrice
	}

	products, total, err := s.productRepo.FindAll(
		offset,
		limit,
		category,
		brand,
		search,
		sort,
		color,
		size,
		minPrice,
		maxPrice,
	)

	if err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	return &ProductListResponse{
		Products:   products,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
