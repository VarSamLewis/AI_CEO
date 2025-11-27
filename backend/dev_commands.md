# Development Commands Reference

## Server Management

### Start Server
```bash
go run main.go
```

### Start Server (with auto-reload using Air)
```bash
# Install air first: go install github.com/cosmtrek/air@latest
air
```

### Build for Production
```bash
go build -o server main.go
./server
```

---

## Database (Turso)

### Check Database Connection
```bash
# Using health endpoint
curl http://localhost:8080/health/db
```

### Connect to Turso CLI
```bash
turso db shell <your-database-name>
```

### View Users Table
```bash
turso db shell <your-database-name> "SELECT * FROM users;"
```

### Delete All Users (for testing)
```bash
turso db shell <your-database-name> "DELETE FROM users;"
```

### Check Database URL
```bash
echo $TURSO_DATABASE_URL
```

---

## Health Checks

### Basic Health
```bash
curl http://localhost:8080/health
```

### Database Health
```bash
curl http://localhost:8080/health/db
```

### LLM Health
```bash
curl http://localhost:8080/health/llm
```

---

## Authentication Endpoints

### Register New User (Returns JWT Token)
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "message": "User registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "test@example.com"
  }
}
```

### Login (Returns JWT Token)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "test@example.com"
  }
}
```

### Test Invalid Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }'
```

### Save Token to Variable
```bash
# Register or login and extract token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  | jq -r '.token')

echo $TOKEN
```

---

## Protected Endpoints (Require JWT)

### Get User Profile
```bash
# First, get a token (see above)
TOKEN="your-jwt-token-here"

curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Using Saved Token Variable
```bash
# After saving token to $TOKEN variable
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Test Unauthorized Access (No Token)
```bash
curl http://localhost:8080/api/profile
```

**Response:**
```json
{
  "error": "Authorization header required",
  "status": "error"
}
```

### Test Invalid Token
```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer invalid-token"
```

**Response:**
```json
{
  "error": "Invalid or expired token",
  "status": "error"
}
```

---

## Other Endpoints

### Echo Test
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello World"
  }'
```

### LLM Request
```bash
curl -X POST http://localhost:8080/llm \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "What is Go?"
  }'
```

---

## Go Commands

### Install Dependencies
```bash
go mod tidy
```

### Add New Dependency
```bash
go get github.com/package/name
```

### Run Tests
```bash
go test ./...
```

### Run E2E Tests
```bash
cd backend
go test ./test -v
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run E2E Tests with Detailed Output
```bash
go test ./test -v -run TestE2E
```

### Format Code
```bash
go fmt ./...
```

### Check for Issues
```bash
go vet ./...
```

---

## Environment Variables

### Check Current Environment
```bash
cat .env
```

### Required Variables
```bash
# .env file should contain:
TURSO_DATABASE_URL=libsql://your-database.turso.io
TURSO_AUTH_TOKEN=your-auth-token
ANTHROPIC_API_KEY=your-api-key
JWT_SECRET=your-secret-key-min-32-characters-long
PORT=8080
```

**Note:** If `JWT_SECRET` is not set, a default development secret will be used. **Always set this in production!**

---

## Git Commands

### Check Status
```bash
git status
```

### Create Commit
```bash
git add .
git commit -m "Your message"
```

### Push to Remote
```bash
git push origin dev
```

---

## Quick Test Flow

### Full Authentication Test with JWT
```bash
# 1. Start server (in one terminal)
go run main.go

# 2. Register user and save token (in another terminal)
TOKEN=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123"}' \
  | jq -r '.token')

echo "Token: $TOKEN"

# 3. Access protected endpoint with token
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"

# 4. Try accessing without token (should fail)
curl http://localhost:8080/api/profile

# 5. Login with correct credentials
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123"}'

# 6. Try wrong password
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "wrong"}'
```

---

## Useful Curl Options

### Pretty Print JSON Response
```bash
curl http://localhost:8080/health | jq
```

### Show Response Headers
```bash
curl -i http://localhost:8080/health
```

### Verbose Output (debugging)
```bash
curl -v http://localhost:8080/health
```

### Save Response to File
```bash
curl http://localhost:8080/health > response.json
```

---

## Debugging

### Check if Server is Running
```bash
lsof -i :8080
# or
netstat -an | grep 8080
```

### Kill Process on Port 8080
```bash
lsof -ti:8080 | xargs kill -9
```

### View Server Logs (if running in background)
```bash
tail -f server.log
```

---

## Production Deployment

### Build Binary
```bash
CGO_ENABLED=0 GOOS=linux go build -o server main.go
```

### Run in Production Mode
```bash
export GIN_MODE=release
./server
```
