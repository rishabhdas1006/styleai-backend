package handler

import (
	"errors"
	"net/http"
	"strconv"

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

	userID := c.GetUint(common.ContextUserID)

	product, err := h.productService.CreateProduct(
		req.Name,
		req.Description,
		req.Brand,
		req.CategoryID,
		req.Gender,
		req.PrimaryImageURL,
		req.PrimaryImagePublicID,
		userID,
	)

	if err != nil {
		if errors.Is(err, common.ErrInvalidGender) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, common.ErrPrimaryImageRequired) {
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
	var query dto.GetProductsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid query parameters",
		})
		return
	}

	userID := c.GetUint(common.ContextUserID)

	println(userID)

	if query.Mine {
		query.CreatedByID = &userID
	}

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
