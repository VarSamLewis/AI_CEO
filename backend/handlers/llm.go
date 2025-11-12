package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the LLM service
	response, err := llm.CallAnthropic(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Optionally save to DB in the future
	// db.SaveConversation(req.UserID, req.Message, response)

	c.JSON(http.StatusOK, LLMResponse{Response: response})
}
