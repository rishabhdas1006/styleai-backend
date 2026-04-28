package handler

import (
	"net/http"
	"strconv"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := h.orderService.Checkout(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"order": order,
	})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	orders, err := h.orderService.GetOrdersByUser(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	idParam := c.Param("id")
	orderID, _ := strconv.Atoi(idParam)

	order, err := h.orderService.GetOrderByID(userID, uint(orderID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}
