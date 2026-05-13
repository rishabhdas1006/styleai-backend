package dto

import (
	"styleai-backend/internal/models"
)

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
}

type GetProductsQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=50"`

	Category string `form:"category"`
	Brand    string `form:"brand"`
	Search   string `form:"search"`
	Gender   string `form:"gender"`

	Color string `form:"color"`
	Size  string `form:"size"`

	MinPrice float64 `form:"minPrice"`
	MaxPrice float64 `form:"maxPrice"`

	Sort string `form:"sort"`
}

type ProductListResponse struct {
	Products   []models.Product `json:"products"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int64            `json:"totalPages"`
}
