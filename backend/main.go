package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"

	"backend/handlers"
)

type EchoRequest struct {
	Message string `json:"message" binding:"required"`
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	//Echo endpoint
	r.POST("/echo", func(c *gin.Context) {
		var req EchoRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"echo":req.Message,
		})
	})

	// LLM endpoint
	r.POST("/llm", handlers.HandleLLMRequest)

	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)

}

