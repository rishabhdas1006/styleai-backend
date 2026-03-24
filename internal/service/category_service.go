package service

import (
	"errors"
	"styleai-backend/internal/common"
	"styleai-backend/internal/models"
	"styleai-backend/internal/repository"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: repo}
}

func (s *CategoryService) CreateCategory(name string) (*models.Category, error) {

	existing, err := s.categoryRepo.FindByName(name)

	if err == nil && existing != nil {
		return nil, common.ErrCategoryExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &models.Category{
		Name: name,
	}

	err = s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {

	category, err := s.categoryRepo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	return category, nil
}
