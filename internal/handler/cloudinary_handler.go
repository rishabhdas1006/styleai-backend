package handler

import (
	"net/http"
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
