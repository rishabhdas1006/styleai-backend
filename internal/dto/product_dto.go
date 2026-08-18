package dto

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Brand       string `json:"brand" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Gender      string `json:"gender" binding:"required"`

	PrimaryImageURL      string `json:"primary_image_url" binding:"required"`
	PrimaryImagePublicID string `json:"primary_image_public_id" binding:"required"`
}

type GetProductsQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=50"`

	Category string `form:"category"`
	Brand    string `form:"brand"`
	Search   string `form:"search"`
	Gender   string `form:"gender"`

	Mine        bool `form:"mine"`
	CreatedByID *uint

	Color string `form:"color"`
	Size  string `form:"size"`

	MinPrice float64 `form:"minPrice"`
	MaxPrice float64 `form:"maxPrice"`

	Sort string `form:"sort"`
}

type ProductListItem struct {
	ID              uint    `json:"ID"`
	Name            string  `json:"Name"`
	Description     string  `json:"Description"`
	Brand           string  `json:"Brand"`
	CategoryID      uint    `json:"CategoryID"`
	CategoryName    string  `json:"CategoryName"`
	Gender          string  `json:"Gender"`
	MinPrice        float64 `json:"MinPrice"`
	PrimaryImageURL string  `json:"PrimaryImageURL"`
}

type ProductListResponse struct {
	Products   []ProductListItem `json:"products"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int64             `json:"totalPages"`
}
