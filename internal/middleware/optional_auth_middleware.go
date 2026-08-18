package middleware

import (
	"strings"

	"styleai-backend/internal/common"
	"styleai-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func OptionalAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// No token provided → continue as guest
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")

		// Invalid header format → continue as guest
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := utils.ValidateToken(tokenString, secret)

		// Invalid token → continue as guest
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		userIDValue, ok := claims[common.ContextUserID]
		if !ok {
			c.Next()
			return
		}

		roleValue, ok := claims[common.ContextRole]
		if !ok {
			c.Next()
			return
		}

		userID, ok := userIDValue.(float64)
		if !ok {
			c.Next()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Next()
			return
		}

		c.Set(common.ContextUserID, uint(userID))
		c.Set(common.ContextRole, role)

		c.Next()
	}
}
