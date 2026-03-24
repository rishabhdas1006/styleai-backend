package handler

import (
	"errors"
	"net/http"
	"strconv"

	"styleai-backend/internal/common"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type VariantHandler struct {
	variantService *service.VariantService
}

func NewVariantHandler(service *service.VariantService) *VariantHandler {
	return &VariantHandler{variantService: service}
}

type CreateVariantRequest struct {
	Size  string  `json:"size" binding:"required"`
	Color string  `json:"color" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}

type UpdateVariantRequest struct {
	Price *float64 `json:"price"`
	Stock *int     `json:"stock"`
}

func (h *VariantHandler) CreateVariant(c *gin.Context) {

	productIDParam := c.Param("id")
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req CreateVariantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	variant, err := h.variantService.CreateVariant(
		uint(productID),
		req.Size,
		req.Color,
		req.Price,
		req.Stock,
	)

	if err != nil {

		if errors.Is(err, common.ErrVariantExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create variant",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"variant": variant,
	})
}

func (h *VariantHandler) GetVariants(c *gin.Context) {

	productIDParam := c.Param("id")
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	variants, err := h.variantService.GetVariants(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch variants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"variants": variants,
	})
}

func (h *VariantHandler) UpdateVariant(c *gin.Context) {

	variantIDParam := c.Param("id")
	variantID, err := strconv.Atoi(variantIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	var req UpdateVariantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Price == nil && req.Stock == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least one field (price or stock) must be provided",
		})
		return
	}

	variant, err := h.variantService.UpdateVariant(
		uint(variantID),
		req.Price,
		req.Stock,
	)

	if err != nil {

		if errors.Is(err, common.ErrVariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update variant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"variant": variant,
	})
}

func (h *VariantHandler) DeleteVariant(c *gin.Context) {

	variantIDParam := c.Param("id")
	variantID, err := strconv.Atoi(variantIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid variant id"})
		return
	}

	err = h.variantService.DeleteVariant(uint(variantID))

	if err != nil {

		if errors.Is(err, common.ErrVariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete variant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "variant deleted successfully",
	})
}
