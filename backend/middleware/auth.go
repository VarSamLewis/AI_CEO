package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"backend/auth"
	"backend/handlers"
)

// AuthMiddleware verifies JWT token and adds user info to context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// 2. Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		// 3. Verify token
		token := parts[1]
		claims, err := auth.VerifyToken(token)
		if err != nil {
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// 4. Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		// 5. Continue to next handler
		c.Next()
	}
}
