package service

import (
	"errors"
	"fmt"
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"
	"styleai-backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VariantService struct {
	variantRepo      *repository.VariantRepository
	variantImageRepo *repository.VariantImageRepository
	productRepo      *repository.ProductRepository
	db               *gorm.DB
	cldManager       *utils.CloudinaryManager
}

func NewVariantService(
	variantRepo *repository.VariantRepository,
	imageRepo *repository.VariantImageRepository,
	productRepo *repository.ProductRepository,
	db *gorm.DB,
	cldManager *utils.CloudinaryManager,
) *VariantService {
	return &VariantService{
		variantRepo:      variantRepo,
		variantImageRepo: imageRepo,
		productRepo:      productRepo,
		db:               db,
		cldManager:       cldManager,
	}
}

func generateSKU(productID uint, color, size string) string {
	return fmt.Sprintf("P%d-%s-%s",
		productID,
		strings.ToUpper(color[:3]),
		size,
	)
}

type ImageInput struct {
	URL      string
	PublicID string
}

func (s *VariantService) CreateVariant(
	variantID string,
	productID uint,
	size, color string,
	price float64,
	stock int,
	images []ImageInput,
) (*models.ProductVariant, error) {

	if _, err := uuid.Parse(variantID); err != nil {
		return nil, common.ErrInvalidVariantID
	}

	size = strings.ToUpper(strings.TrimSpace(size))
	color = strings.ToLower(strings.TrimSpace(color))

	existing, err := s.variantRepo.FindByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants: %w", err)
	}

	for _, v := range existing {
		if v.Size == size && v.Color == color {
			return nil, common.ErrVariantExists
		}
	}

	sku := generateSKU(productID, color, size)

	variant := &models.ProductVariant{
		ID:        variantID,
		ProductID: productID,
		Size:      size,
		Color:     color,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
	}

	prefix := fmt.Sprintf("products/%d/%s", productID, variantID)
	tx := s.db.Begin()

	if err := s.variantRepo.WithTx(tx).Create(variant); err != nil {
		tx.Rollback()
		_ = s.cldManager.DeleteByPrefix(prefix)

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, common.ErrVariantExists
		}

		return nil, fmt.Errorf("create variant failed: %w", err)
	}

	var variantImages []models.VariantImage
	for i, img := range images {
		if img.URL == "" || img.PublicID == "" {
			tx.Rollback()
			_ = s.cldManager.DeleteByPrefix(prefix)
			return nil, common.ErrInvalidImageData
		}

		variantImages = append(variantImages, models.VariantImage{
			VariantID: variantID,
			ImageURL:  img.URL,
			PublicID:  img.PublicID,
			Position:  i,
		})
	}

	if len(variantImages) > 0 {
		if err := s.variantImageRepo.WithTx(tx).CreateImages(variantImages); err != nil {
			tx.Rollback()
			_ = s.cldManager.DeleteByPrefix(prefix)
			return nil, fmt.Errorf("create images failed: %w", err)
		}
	}

	if err := s.productRepo.WithTx(tx).UpdateMinPrice(productID); err != nil {
		tx.Rollback()
		_ = s.cldManager.DeleteByPrefix(prefix)
		return nil, fmt.Errorf("update min price failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		_ = s.cldManager.DeleteByPrefix(prefix)
		return nil, common.ErrTransactionFailed
	}

	return variant, nil
}

func (s *VariantService) GetVariants(productID uint) ([]models.ProductVariant, error) {
	variants, err := s.variantRepo.FindByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants: %w", err)
	}
	return variants, nil
}

func (s *VariantService) UpdateVariant(id string, price *float64, stock *int) (*models.ProductVariant, error) {

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

	if err := s.variantRepo.Update(variant); err != nil {
		return nil, fmt.Errorf("update variant failed: %w", err)
	}

	if err := s.productRepo.UpdateMinPrice(variant.ProductID); err != nil {
		return nil, fmt.Errorf("update min price failed: %w", err)
	}

	return variant, nil
}

func (s *VariantService) DeleteVariant(id string) error {

	variant, err := s.variantRepo.FindByID(id)
	if err != nil {
		return common.ErrVariantNotFound
	}

	if err := s.variantRepo.Delete(id); err != nil {
		return fmt.Errorf("delete variant failed: %w", err)
	}

	if err := s.productRepo.UpdateMinPrice(variant.ProductID); err != nil {
		return fmt.Errorf("update min price failed: %w", err)
	}

	return nil
}
