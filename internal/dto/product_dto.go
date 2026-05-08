package dto

import (
	"time"

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
	Products   []ProductSummaryResponse `json:"products"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	Total      int64                    `json:"total"`
	TotalPages int64                    `json:"totalPages"`
}

type ProductSummaryResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Brand       string                  `json:"brand"`
	Gender      string                  `json:"gender"`
	MinPrice    float64                 `json:"minPrice"`
	Category    ProductCategoryResponse `json:"category"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type ProductDetailResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Brand       string                   `json:"brand"`
	Gender      string                   `json:"gender"`
	MinPrice    float64                  `json:"minPrice"`
	Category    ProductCategoryResponse  `json:"category"`
	Variants    []ProductVariantResponse `json:"variants"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

type ProductCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductVariantResponse struct {
	ID     string                 `json:"id"`
	Size   string                 `json:"size"`
	Color  string                 `json:"color"`
	SKU    string                 `json:"sku"`
	Price  float64                `json:"price"`
	Stock  int                    `json:"stock"`
	Images []ProductImageResponse `json:"images,omitempty"`
}

type ProductImageResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

func ToProductSummaryResponse(product models.Product) ProductSummaryResponse {
	return ProductSummaryResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Brand:       product.Brand,
		Gender:      product.Gender,
		MinPrice:    product.MinPrice,
		Category:    ToProductCategoryResponse(product.Category),
		CreatedAt:   product.CreatedAt,
	}
}

func ToProductDetailResponse(product models.Product) ProductDetailResponse {
	variants := make([]ProductVariantResponse, 0, len(product.Variants))
	for _, variant := range product.Variants {
		variants = append(variants, ToProductVariantResponse(variant))
	}

	return ProductDetailResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Brand:       product.Brand,
		Gender:      product.Gender,
		MinPrice:    product.MinPrice,
		Category:    ToProductCategoryResponse(product.Category),
		Variants:    variants,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ToProductListResponse(products []models.Product, page, limit int, total, totalPages int64) ProductListResponse {
	productResponses := make([]ProductSummaryResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, ToProductSummaryResponse(product))
	}

	return ProductListResponse{
		Products:   productResponses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

func ToProductCategoryResponse(category models.Category) ProductCategoryResponse {
	return ProductCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func ToProductVariantResponse(variant models.ProductVariant) ProductVariantResponse {
	images := make([]ProductImageResponse, 0, len(variant.Images))
	for _, image := range variant.Images {
		images = append(images, ToProductImageResponse(image))
	}

	return ProductVariantResponse{
		ID:     variant.ID,
		Size:   variant.Size,
		Color:  variant.Color,
		SKU:    variant.SKU,
		Price:  variant.Price,
		Stock:  variant.Stock,
		Images: images,
	}
}

func ToProductImageResponse(image models.VariantImage) ProductImageResponse {
	return ProductImageResponse{
		ID:       image.ID,
		URL:      image.ImageURL,
		Position: image.Position,
	}
}
