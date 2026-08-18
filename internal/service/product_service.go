package service

import (
	"errors"
	"math"
	"strings"
	"styleai-backend/internal/common"
	"styleai-backend/internal/dto"
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

func (s *ProductService) CreateProduct(
	name string,
	description string,
	brand string,
	categoryID uint,
	gender string,
	primaryImageURL string,
	primaryImagePublicID string,
	createdByID uint,
) (*models.Product, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	brand = strings.ToLower(strings.TrimSpace(brand))
	gender = strings.ToLower(strings.TrimSpace(gender))

	if !common.IsValidGender(common.Gender(gender)) {
		return nil, common.ErrInvalidGender
	}

	if primaryImageURL == "" || primaryImagePublicID == "" {
		return nil, common.ErrPrimaryImageRequired
	}

	product := &models.Product{
		Name:                 name,
		Description:          description,
		Brand:                brand,
		CategoryID:           categoryID,
		Gender:               gender,
		CreatedByID:          createdByID,
		PrimaryImageURL:      primaryImageURL,
		PrimaryImagePublicID: primaryImagePublicID,
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

func (s *ProductService) GetProducts(filter dto.GetProductsQuery) (*dto.ProductListResponse, error) {

	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit <= 0 || filter.Limit > 50 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit

	filter.Category = strings.ToLower(strings.TrimSpace(filter.Category))
	filter.Brand = strings.ToLower(strings.TrimSpace(filter.Brand))
	filter.Search = strings.ToLower(strings.TrimSpace(filter.Search))
	filter.Gender = strings.ToLower(strings.TrimSpace(filter.Gender))
	filter.Color = strings.ToLower(strings.TrimSpace(filter.Color))
	filter.Size = strings.ToLower(strings.TrimSpace(filter.Size))
	filter.Sort = strings.TrimSpace(filter.Sort)

	if filter.MinPrice < 0 {
		filter.MinPrice = 0
	}

	if filter.MaxPrice > 0 && filter.MaxPrice < filter.MinPrice {
		return nil, common.ErrMaxPriceLessThanMinPrice
	}

	if filter.Gender != "" && !common.IsValidGender(common.Gender(filter.Gender)) {
		return nil, common.ErrInvalidGender
	}

	products, total, err := s.productRepo.FindAll(offset, filter)

	if err != nil {
		return nil, err
	}

	productItems := make([]dto.ProductListItem, 0, len(products))

	for _, product := range products {
		productItems = append(productItems, dto.ProductListItem{
			ID:              product.ID,
			Name:            product.Name,
			Description:     product.Description,
			Brand:           product.Brand,
			CategoryID:      product.CategoryID,
			CategoryName:    product.Category.Name,
			Gender:          product.Gender,
			MinPrice:        product.MinPrice,
			PrimaryImageURL: product.PrimaryImageURL,
		})
	}

	totalPages := int64(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.ProductListResponse{
		Products:   productItems,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
