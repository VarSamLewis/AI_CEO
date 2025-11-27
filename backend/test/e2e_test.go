package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"backend/auth"
	"backend/handlers"
	"backend/middleware"
	db "backend/database"
)

// Test response structures
type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ProfileResponse struct {
	User struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

// setupTestRouter creates a router with all routes for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Health checks
	r.GET("/health", handlers.HealthCheck)
	r.GET("/health/db", handlers.DBHealthCheck)

	// Auth routes
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)

	// Protected routes
	r.GET("/api/profile", middleware.AuthMiddleware(), handlers.GetProfile)

	return r
}

// setupTestDB initializes the database for testing
func setupTestDB(t *testing.T) {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("Warning: No .env file found: %v", err)
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create users table
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.CreateUsersTable(ctx); err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}
}

// cleanupTestDB removes test data
func cleanupTestDB(t *testing.T, email string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.DB.ExecContext(ctx, "DELETE FROM users WHERE email = ?", email)
	if err != nil {
		t.Logf("Warning: Failed to cleanup test user: %v", err)
	}
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Run tests
	m.Run()

	// Teardown
	if db.DB != nil {
		db.DB.Close()
	}
}

// TestE2E_CompleteAuthFlow tests the entire authentication flow
func TestE2E_CompleteAuthFlow(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()

	testEmail := "e2e_test@example.com"
	testPassword := "testpass123"
	defer cleanupTestDB(t, testEmail)

	t.Run("1. Register new user", func(t *testing.T) {
		payload := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Parse response
		var resp AuthResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Verify response
		if resp.Token == "" {
			t.Error("Expected token in response, got empty string")
		}
		if resp.User.Email != testEmail {
			t.Errorf("Expected email %s, got %s", testEmail, resp.User.Email)
		}
		if resp.User.ID == 0 {
			t.Error("Expected user ID > 0, got 0")
		}
		if resp.Message != "User registered successfully" {
			t.Errorf("Expected success message, got %s", resp.Message)
		}

		t.Logf("✓ User registered successfully with ID: %d", resp.User.ID)
		t.Logf("✓ Token received: %s...", resp.Token[:20])
	})

	t.Run("2. Login with correct credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Parse response
		var resp AuthResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Verify response
		if resp.Token == "" {
			t.Error("Expected token in response, got empty string")
		}
		if resp.User.Email != testEmail {
			t.Errorf("Expected email %s, got %s", testEmail, resp.User.Email)
		}
		if resp.Message != "Login successful" {
			t.Errorf("Expected login success message, got %s", resp.Message)
		}

		t.Logf("✓ Login successful")
		t.Logf("✓ Token received: %s...", resp.Token[:20])
	})

	t.Run("3. Access protected route with valid token", func(t *testing.T) {
		// First login to get token
		payload := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		loginReq.Header.Set("Content-Type", "application/json")
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)

		var loginResp AuthResponse
		json.Unmarshal(loginW.Body.Bytes(), &loginResp)
		token := loginResp.Token

		// Now access protected route
		req := httptest.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		// Parse response
		var resp ProfileResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Verify response
		if resp.User.Email != testEmail {
			t.Errorf("Expected email %s, got %s", testEmail, resp.User.Email)
		}
		if resp.User.ID == 0 {
			t.Error("Expected user ID > 0, got 0")
		}

		t.Logf("✓ Protected route accessed successfully")
		t.Logf("✓ User profile retrieved: ID=%d, Email=%s", resp.User.ID, resp.User.Email)
	})

	t.Run("4. Reject duplicate registration", func(t *testing.T) {
		payload := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusConflict {
			t.Fatalf("Expected status 409 (Conflict), got %d", w.Code)
		}

		// Parse error response
		var resp ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Message != "User already exists" {
			t.Errorf("Expected 'User already exists' error, got %s", resp.Message)
		}

		t.Logf("✓ Duplicate registration correctly rejected")
	})

	t.Run("5. Reject login with wrong password", func(t *testing.T) {
		payload := map[string]string{
			"email":    testEmail,
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 (Unauthorized), got %d", w.Code)
		}

		// Parse error response
		var resp ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Message != "Invalid credentials" {
			t.Errorf("Expected 'Invalid credentials' error, got %s", resp.Message)
		}

		t.Logf("✓ Wrong password correctly rejected")
	})

	t.Run("6. Reject login with non-existent email", func(t *testing.T) {
		payload := map[string]string{
			"email":    "nonexistent@example.com",
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 (Unauthorized), got %d", w.Code)
		}

		// Parse error response
		var resp ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Message != "Invalid credentials" {
			t.Errorf("Expected 'Invalid credentials' error, got %s", resp.Message)
		}

		t.Logf("✓ Non-existent email correctly rejected")
	})

	t.Run("7. Reject protected route without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/profile", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 (Unauthorized), got %d", w.Code)
		}

		// Parse error response
		var resp ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Message != "Authorization header required" {
			t.Errorf("Expected 'Authorization header required' error, got %s", resp.Message)
		}

		t.Logf("✓ Missing token correctly rejected")
	})

	t.Run("8. Reject protected route with invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 (Unauthorized), got %d", w.Code)
		}

		// Parse error response
		var resp ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if resp.Message != "Invalid or expired token" {
			t.Errorf("Expected 'Invalid or expired token' error, got %s", resp.Message)
		}

		t.Logf("✓ Invalid token correctly rejected")
	})

	t.Run("9. Reject malformed authorization header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/profile", nil)
		req.Header.Set("Authorization", "NotBearer token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 (Unauthorized), got %d", w.Code)
		}

		t.Logf("✓ Malformed authorization header correctly rejected")
	})

	t.Run("10. Reject invalid email format", func(t *testing.T) {
		payload := map[string]string{
			"email":    "not-an-email",
			"password": testPassword,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400 (Bad Request), got %d", w.Code)
		}

		t.Logf("✓ Invalid email format correctly rejected")
	})

	t.Run("11. Reject password shorter than 6 characters", func(t *testing.T) {
		payload := map[string]string{
			"email":    "short@example.com",
			"password": "12345",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400 (Bad Request), got %d", w.Code)
		}

		t.Logf("✓ Short password correctly rejected")
	})
}

// TestE2E_DatabaseConnection tests database connectivity
func TestE2E_DatabaseConnection(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()

	t.Run("Database health check", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health/db", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		t.Logf("✓ Database connection healthy")
	})
}
