package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	// Extract user info from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userEmail, _ := c.Get("user_email")

	// Return user profile
	SuccessResponse(c, gin.H{
		"user": gin.H{
			"id":    userID,
			"email": userEmail,
		},
	})
}
