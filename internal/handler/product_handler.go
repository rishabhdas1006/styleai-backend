package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/internal/dto"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {

	var req dto.CreateProductRequest

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
		req.Gender,
	)

	if err != nil {
		if errors.Is(err, common.ErrInvalidGender) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": dto.ToProductDetailResponse(*product),
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
		"product": dto.ToProductDetailResponse(*product),
	})
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var query dto.GetProductsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid query parameters",
		})
		return
	}

	query.Category = strings.TrimSpace(query.Category)
	query.Brand = strings.TrimSpace(query.Brand)
	query.Search = strings.TrimSpace(query.Search)
	query.Gender = strings.TrimSpace(query.Gender)
	query.Color = strings.TrimSpace(query.Color)
	query.Size = strings.TrimSpace(query.Size)
	query.Sort = strings.TrimSpace(query.Sort)

	validSorts := map[string]bool{
		"price_asc":  true,
		"price_desc": true,
		"newest":     true,
	}

	if !validSorts[query.Sort] {
		query.Sort = ""
	}

	result, err := h.productService.GetProducts(query)

	if err != nil {
		if errors.Is(err, common.ErrMaxPriceLessThanMinPrice) || errors.Is(err, common.ErrInvalidGender) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
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
