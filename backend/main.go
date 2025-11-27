package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/handlers"
	db "backend/database"
	"backend/auth"
	"backend/middleware"
)

type Config struct {
	Port            string
	DatabaseURL     string
	AnthropicAPIKey string
}

func loadConfig() *Config {
	return &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		DatabaseURL:     os.Getenv("TURSO_DATABASE_URL"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Load configuration
	cfg := loadConfig()
  
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialise db connection: %v", err)
	}
	defer db.DB.Close()

	// Create users table with context
	if err := db.CreateUsersTable(context.Background()); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
  
	
	r := gin.Default()

	// Health Check
	r.GET("/health", handlers.HealthCheck)
	r.GET("/health/db", handlers.DBHealthCheck)
	r.GET("/health/llm", handlers.LLMHealthCheck)

	// Echo endpoint
	r.POST("/echo", handlers.Echo)

	// LLM endpoint
	r.POST("/llm", handlers.HandleLLMRequest)

	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)

	// Protected routes (require authentication)
	r.GET("/api/profile", middleware.AuthMiddleware(), handlers.GetProfile)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

