package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: service}
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {

	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	product, err := h.productService.CreateProduct(
		req.Name,
		req.Description,
		req.Brand,
		req.CategoryID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, common.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// filters
	category := c.Query("category")
	brand := c.Query("brand")
	search := c.Query("search")

	color := strings.TrimSpace(c.Query("color"))
	size := strings.TrimSpace(c.Query("size"))

	var minPrice, maxPrice float64

	if v := c.Query("minPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid minPrice"})
			return
		}
		minPrice = parsed
	}

	if v := c.Query("maxPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxPrice"})
			return
		}
		maxPrice = parsed
	}

	// sorting
	sort := strings.TrimSpace(c.DefaultQuery("sort", ""))

	validSorts := map[string]bool{
		"price_asc":  true,
		"price_desc": true,
		"newest":     true,
	}

	if !validSorts[sort] {
		sort = ""
	}

	result, err := h.productService.GetProducts(
		page,
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
		if errors.Is(err, common.ErrMaxPriceLessThanMinPrice) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": common.ErrMaxPriceLessThanMinPrice.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
