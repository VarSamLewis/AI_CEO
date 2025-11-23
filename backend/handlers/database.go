package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	db "backend/database"
)

// Standard response helpers
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

func SuccessResponse(c *gin.Context, data gin.H) {
	response := gin.H{"status": "ok"}
	for k, v := range data {
		response[k] = v
	}
	c.JSON(http.StatusOK, response)
}

func HealthCheck(c *gin.Context) {
	SuccessResponse(c, gin.H{})
}

func DBHealthCheck(c *gin.Context) {
	if err := db.DB.Ping(); err != nil {
		ErrorResponse(c, http.StatusServiceUnavailable, err.Error())
		return
	}
	SuccessResponse(c, gin.H{"database": "connected"})
}

type EchoRequest struct {
	Message string `json:"message" binding:"required"`
}

func Echo(c *gin.Context) {
	var req EchoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"echo": req.Message})
}
