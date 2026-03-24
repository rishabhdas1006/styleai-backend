package service

import (
	"errors"
	"fmt"
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"

	"gorm.io/gorm"
)

type VariantService struct {
	variantRepo *repository.VariantRepository
	productRepo *repository.ProductRepository
}

func NewVariantService(variantRepo *repository.VariantRepository, productRepo *repository.ProductRepository) *VariantService {
	return &VariantService{variantRepo: variantRepo, productRepo: productRepo}
}

func generateSKU(productID uint, color, size string) string {
	return fmt.Sprintf("P%d-%s-%s",
		productID,
		strings.ToUpper(color[:3]),
		size,
	)
}

func (s *VariantService) CreateVariant(productID uint, size, color string, price float64, stock int) (*models.ProductVariant, error) {

	// normalize input
	size = strings.ToUpper(strings.TrimSpace(size))
	color = strings.ToLower(strings.TrimSpace(color))

	// check duplicates
	existing, err := s.variantRepo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	for _, v := range existing {
		if v.Size == size && v.Color == color {
			return nil, common.ErrVariantExists
		}
	}

	// generate SKU
	sku := generateSKU(productID, color, size)

	variant := &models.ProductVariant{
		ProductID: productID,
		Size:      size,
		Color:     color,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
	}

	err = s.variantRepo.Create(variant)
	if err != nil {

		// fallback for DB constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, common.ErrVariantExists
		}

		return nil, err
	}

	err = s.productRepo.UpdateMinPrice(productID)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *VariantService) GetVariants(productID uint) ([]models.ProductVariant, error) {
	return s.variantRepo.FindByProductID(productID)
}

func (s *VariantService) UpdateVariant(id uint, price *float64, stock *int) (*models.ProductVariant, error) {

	variant, err := s.variantRepo.FindByID(id)
	if err != nil {
		return nil, common.ErrVariantNotFound
	}

	if price != nil {
		variant.Price = *price
	}

	if stock != nil {
		variant.Stock = *stock
	}

	err = s.variantRepo.Update(variant)
	if err != nil {
		return nil, err
	}

	err = s.productRepo.UpdateMinPrice(variant.ProductID)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (s *VariantService) DeleteVariant(id uint) error {
	variant, err := s.variantRepo.FindByID(id)
	if err != nil {
		return common.ErrVariantNotFound
	}

	err = s.variantRepo.Delete(id)
	if err != nil {
		return err
	}

	return s.productRepo.UpdateMinPrice(variant.ProductID)
}
