package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "backend/database"
)

type PreferencesRequest struct {
	DietaryRestrictions string `json:"dietary_restrictions"`
	MaxCookingTime      int    `json:"max_cooking_time"`
}

type PreferencesResponse struct {
	DietaryRestrictions string `json:"dietary_restrictions"`
	MaxCookingTime      int    `json:"max_cooking_time"`
}

// GetPreferences retrieves user preferences
func GetPreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var prefJSON string
	err := db.DB.QueryRowContext(ctx,
		"SELECT user_preference FROM user_preference WHERE user_id = ?",
		userID,
	).Scan(&prefJSON)

	if err == sql.ErrNoRows {
		// No preferences set yet
		SuccessResponse(c, gin.H{
			"dietary_restrictions": "",
			"max_cooking_time":     0,
		})
		return
	}

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}

	var prefs PreferencesResponse
	if err := json.Unmarshal([]byte(prefJSON), &prefs); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to parse preferences")
		return
	}

	SuccessResponse(c, gin.H{
		"dietary_restrictions": prefs.DietaryRestrictions,
		"max_cooking_time":     prefs.MaxCookingTime,
	})
}

// UpdatePreferences sets or updates user preferences
func UpdatePreferences(c *gin.Context) {
	var req PreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Convert preferences to JSON
	prefJSON, err := json.Marshal(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to encode preferences")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Upsert preferences (insert or update)
	_, err = db.DB.ExecContext(ctx,
		`INSERT INTO user_preference (user_id, user_preference, updated_at)
		 VALUES (?, ?, datetime('now'))
		 ON CONFLICT(user_id) DO UPDATE SET
		 user_preference = excluded.user_preference,
		 updated_at = excluded.updated_at`,
		userID, string(prefJSON),
	)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to save preferences")
		return
	}

	SuccessResponse(c, gin.H{
		"message":              "Preferences updated successfully",
		"dietary_restrictions": req.DietaryRestrictions,
		"max_cooking_time":     req.MaxCookingTime,
	})
}
