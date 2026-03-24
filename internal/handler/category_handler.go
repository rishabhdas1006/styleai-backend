package handler

import (
	"errors"
	"net/http"
	"strconv"

	"styleai-backend/internal/common"
	"styleai-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: service}
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {

	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	category, err := h.categoryService.CreateCategory(req.Name)
	if err != nil {

		if errors.Is(err, common.ErrCategoryExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"category": category,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category id",
		})
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {

		if errors.Is(err, common.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}
