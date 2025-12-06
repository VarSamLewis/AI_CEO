package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"backend/auth"
	"backend/handlers"
)

// AuthMiddleware verifies JWT token from httpOnly cookie and adds user info to context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get token from cookie
		token, err := c.Cookie("token")
		if err != nil {
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		// 2. Verify token
		claims, err := auth.VerifyToken(token)
		if err != nil {
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// 3. Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		// 4. Continue to next handler
		c.Next()
	}
}
