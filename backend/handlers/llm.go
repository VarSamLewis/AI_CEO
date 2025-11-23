package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"backend/llm"
)

type LLMRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id,omitempty"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

func HandleLLMRequest(c *gin.Context) {
	var req LLMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call the LLM service
	response, err := llm.CallAnthropic(req.Message)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Optionally save to DB in the future
	// db.SaveConversation(req.UserID, req.Message, response)

	SuccessResponse(c, gin.H{"response": response})
}

func LLMHealthCheck(c *gin.Context) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		ErrorResponse(c, http.StatusServiceUnavailable, "ANTHROPIC_API_KEY not configured")
		return
	}

	SuccessResponse(c, gin.H{"llm": "configured"})
}
