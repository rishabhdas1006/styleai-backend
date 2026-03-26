package handler

import (
	"net/http"
	"strconv"
	"styleai-backend/internal/common"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

type AddItemRequest struct {
	VariantID uint `json:"variant_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req AddItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	err := h.cartService.AddItem(userID, req.VariantID, req.Quantity)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "item added"})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, total, err := h.cartService.GetCart(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"cart":  cart,
		"total": total,
	})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	var req UpdateItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.cartService.UpdateItem(userID, uint(itemID), req.Quantity)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("user_id")

	err := h.cartService.RemoveItem(userID, uint(itemID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case common.ErrInvalidQuantity:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case common.ErrInsufficientStock:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case common.ErrUnauthorized:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case common.ErrVariantNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case common.ErrCartItemNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
