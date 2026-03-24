package middleware

import (
	"net/http"
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": common.ErrMissingAuthToken.Error(),
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": common.ErrInvalidAuthToken.Error(),
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := utils.ValidateToken(tokenString, secret)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": common.ErrInvalidAuthToken.Error(),
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": common.ErrInvalidAuthToken.Error(),
			})
			c.Abort()
			return
		}

		userID := uint(claims[common.ContextUserID].(float64))
		role := claims["role"].(string)

		c.Set(common.ContextUserID, userID)
		c.Set(common.ContextRole, role)

		c.Next()
	}
}
