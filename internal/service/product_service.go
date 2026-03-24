package service

import (
	"errors"
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
	Products []models.Product `json:"products"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
	Total    int64            `json:"total"`
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
) (*ProductListResponse, error) {

	if page < 1 {
		page = 1
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.FindAll(
		offset,
		limit,
		category,
		brand,
		search,
		sort,
	)

	if err != nil {
		return nil, err
	}

	return &ProductListResponse{
		Products: products,
		Page:     page,
		Limit:    limit,
		Total:    total,
	}, nil
}
