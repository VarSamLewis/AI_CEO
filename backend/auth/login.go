package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"backend/handlers"
	db "backend/database"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	// 1. Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Find user by email
	ctx_1, cancel_1 := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel_1()

	var userID int64
	var hashedPassword string
	err := db.DB.QueryRowContext(ctx_1, "SELECT id, hashed_password FROM users WHERE email = ?", req.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal whether email exists or not (security best practice)
			handlers.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}

	// 3. Verify password
	if err := VerifyPassword(hashedPassword, req.Password); err != nil {
		handlers.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// 4. Generate JWT token
	token, err := GenerateToken(userID, req.Email)
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 5. Return success with token
	handlers.SuccessResponse(c, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":    userID,
			"email": req.Email,
		},
	})
}
