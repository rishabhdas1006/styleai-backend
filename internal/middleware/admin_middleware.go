package middleware

import (
	"net/http"

	"styleai-backend/internal/common"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get(common.ContextRole)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": common.ErrForbidden.Error(),
			})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": common.ErrForbidden.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
