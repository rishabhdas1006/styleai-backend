package handler

import (
	"net/http"
	"strings"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CloudinaryHandler struct {
	cloudinaryService *service.CloudinaryService
}

func NewCloudinaryHandler(s *service.CloudinaryService) *CloudinaryHandler {
	return &CloudinaryHandler{cloudinaryService: s}
}

type SignatureRequest struct {
	ProductID string `json:"productId"`
	VariantID string `json:"variantId"`
}

type CleanupRequest struct {
	ProductID string `json:"productId" binding:"required"`
	VariantID string `json:"variantId" binding:"required"`
}

func (h *CloudinaryHandler) GetProductImageSignature(c *gin.Context) {
	data, err := h.cloudinaryService.GenerateProductImageSignature()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate upload signature",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

type ProductImageCleanupRequest struct {
	Folder string `json:"folder" binding:"required"`
}

func (h *CloudinaryHandler) CleanupProductImage(c *gin.Context) {
	var req ProductImageCleanupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.cloudinaryService.DeleteProductImageFolder(req.Folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cleanup product image",
		})
		return
	}

	if !strings.HasPrefix(req.Folder, "products/pending/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product image folder",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product image cleaned up successfully",
	})
}

func (h *CloudinaryHandler) GetSignature(c *gin.Context) {

	var req SignatureRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.ProductID == "" || req.VariantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "productId and variantId are required",
		})
		return
	}

	if _, err := uuid.Parse(req.VariantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid variantId format",
		})
		return
	}

	data, err := h.cloudinaryService.GenerateSignature(req.ProductID, req.VariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *CloudinaryHandler) CleanupVariantImages(c *gin.Context) {
	var req CleanupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if _, err := uuid.Parse(req.VariantID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid variantId format",
		})
		return
	}

	if err := h.cloudinaryService.DeleteVariantImages(
		req.ProductID,
		req.VariantID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cleanup variant images",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "variant images cleaned up successfully",
	})
}
