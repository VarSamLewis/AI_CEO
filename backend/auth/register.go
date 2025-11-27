package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"backend/handlers"
	db "backend/database"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 2. Check if user exists
	ctx_1, cancel_1 := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel_1()

	var exists bool
	err := db.DB.QueryRowContext(ctx_1, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", req.Email).Scan(&exists)
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		handlers.ErrorResponse(c, http.StatusConflict, "User already exists")
		return
	}

	// 3. Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// 4. Insert into database
	ctx_2, cancel_2 := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel_2()

	result, err := db.DB.ExecContext(ctx_2,
		"INSERT INTO users (email, hashed_password, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
		req.Email, hashedPassword,
	)
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// 5. Get the newly created user ID
	userID, err := result.LastInsertId()
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user ID")
		return
	}

	// 6. Generate JWT token (auto-login after registration)
	token, err := GenerateToken(userID, req.Email)
	if err != nil {
		handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 7. Return success with token
	handlers.SuccessResponse(c, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": gin.H{
			"id":    userID,
			"email": req.Email,
		},
	})
}
