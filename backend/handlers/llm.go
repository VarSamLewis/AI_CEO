package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"backend/llm"
	db "backend/database"
)

type LLMRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id,omitempty"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

type UserPreferences struct {
	DietaryRestrictions string `json:"dietary_restrictions"`
	MaxCookingTime      int    `json:"max_cooking_time"`
}

type UserUsage struct {
	MealCount int
	MaxMeals  int
}

// getUserPreferences fetches user preferences from database
func getUserPreferences(userID int64) (*UserPreferences, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var prefJSON string
	err := db.DB.QueryRowContext(ctx,
		"SELECT user_preference FROM user_preference WHERE user_id = ?",
		userID,
	).Scan(&prefJSON)

	if err == sql.ErrNoRows {
		// No preferences set, return nil
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON preferences
	var prefs UserPreferences
	if err := json.Unmarshal([]byte(prefJSON), &prefs); err != nil {
		return nil, err
	}

	return &prefs, nil
}

func getUserUsage(userID int64) (*UserUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usage UserUsage
	err := db.DB.QueryRowContext(ctx,
		"SELECT meal_count, max_meals FROM users_tracking WHERE user_id = ?",
		userID,
	).Scan(&usage.MealCount, &usage.MaxMeals)

	if err == sql.ErrNoRows {
		_, err = db.DB.ExecContext(ctx,
			"INSERT INTO users_tracking (user_id, meal_count, max_meals) VALUES (?, 0, 20)",
			userID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create tracking record: %w", err)
		}
		return &UserUsage{MealCount: 0, MaxMeals: 20}, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage: %w", err)
	}

	return &usage, nil
  }

func incrementUserUsage(userID int64) error {
  	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  	defer cancel()

  	_, err := db.DB.ExecContext(ctx,
  		`UPDATE users_tracking
  		 SET meal_count = meal_count + 1,
  		     updated_at = datetime('now')
  		 WHERE user_id = ?`,
  		userID,
  	)

  	if err != nil {
  		return fmt.Errorf("failed to increment usage: %w", err)
  	}

  	return nil
  }


func HandleLLMRequest(c *gin.Context) {
  	var req LLMRequest
  	if err := c.ShouldBindJSON(&req); err != nil {
  		ErrorResponse(c, http.StatusBadRequest, err.Error())
  		return
  	}

  	// Get user ID from JWT token (set by AuthMiddleware)
  	userID, exists := c.Get("user_id")
  	if !exists {
  		ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
  		return
  	}

  	// ✅ CHECK USAGE LIMIT - FETCH FROM DB
  	usage, err := getUserUsage(userID.(int64))
  	if err != nil {
  		ErrorResponse(c, http.StatusInternalServerError, "Failed to check usage limits")
  		return
  	}

  	if usage.MealCount >= usage.MaxMeals {
  		ErrorResponse(c, http.StatusForbidden, fmt.Sprintf(
  			"Usage limit reached. You've used %d/%d free meal generations.",
  			usage.MealCount, usage.MaxMeals,
  		))
  		return
  	}

  	// Fetch user preferences
  	prefs, err := getUserPreferences(userID.(int64))
  	if err != nil {
  		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch preferences")
  		return
  	}

  	// Load system prompt from environment
  	systemPrompt := os.Getenv("LLM_SYSTEM_PROMPT")
  	if systemPrompt == "" {
  		systemPrompt = "You are a helpful meal planning assistant."
  	}

  	// Build user message with preferences
  	userMessage := buildMealPrompt(req.Message, prefs)

  	// Call the LLM service with both system and user prompts
  	response, err := llm.CallAnthropic(systemPrompt, userMessage)
  	if err != nil {
  		ErrorResponse(c, http.StatusInternalServerError, err.Error())
  		return
  	}

  	// ✅ INCREMENT USAGE AFTER SUCCESSFUL LLM CALL
  	if err := incrementUserUsage(userID.(int64)); err != nil {
  		// Log error but don't fail the request since user got their response
  		fmt.Printf("Warning: Failed to increment usage for user %d: %v\n", userID, err)
  	}

  	SuccessResponse(c, gin.H{
  		"response": response,
  		"usage": gin.H{
  			"used":      usage.MealCount + 1,
  			"remaining": usage.MaxMeals - (usage.MealCount + 1),
  			"limit":     usage.MaxMeals,
  		},
  	})
  }


// buildMealPrompt constructs a meal planning prompt with user preferences
func buildMealPrompt(ingredients string, prefs *UserPreferences) string {
	prompt := ingredients

	if prefs != nil {
		if prefs.DietaryRestrictions != "" {
			prompt += fmt.Sprintf("\nDietary restrictions: %s", prefs.DietaryRestrictions)
		}
		if prefs.MaxCookingTime > 0 {
			prompt += fmt.Sprintf("\nMaximum cooking time: %d minutes", prefs.MaxCookingTime)
		}
	}

	return prompt
}

func LLMHealthCheck(c *gin.Context) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		ErrorResponse(c, http.StatusServiceUnavailable, "ANTHROPIC_API_KEY not configured")
		return
	}

	SuccessResponse(c, gin.H{"llm": "configured"})
}

// GetUsage returns the user's meal generation usage statistics
func GetUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	usage, err := getUserUsage(userID.(int64))
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch usage")
		return
	}

	SuccessResponse(c, gin.H{
		"used":      usage.MealCount,
		"remaining": usage.MaxMeals - usage.MealCount,
		"limit":     usage.MaxMeals,
	})
}
