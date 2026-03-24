package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// filters
	category := c.Query("category")
	brand := c.Query("brand")
	search := c.Query("search")

	// sorting
	sort := c.DefaultQuery("sort", "newest")

	result, err := h.productService.GetProducts(page, limit, category, brand, search, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
