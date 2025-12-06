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

## Authentication Endpoints (httpOnly Cookies)

**Note:** Authentication now uses httpOnly cookies instead of Bearer tokens for improved security against XSS attacks.

### Register New User (Sets httpOnly Cookie)
```bash
curl -i -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' \
  -c cookies.txt
```

**Response:**
```json
{
  "message": "User registered successfully",
  "status": "ok",
  "user": {
    "id": 1,
    "email": "test@example.com"
  }
}
```

**Important:** The JWT token is now set as an httpOnly cookie (see `Set-Cookie` header), not in the response body.

### Login (Sets httpOnly Cookie)
```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' \
  -c cookies.txt
```

**Response:**
```json
{
  "message": "Login successful",
  "status": "ok",
  "user": {
    "id": 1,
    "email": "test@example.com"
  }
}
```

The `-c cookies.txt` flag saves the cookie to a file for subsequent requests.

### Logout (Clears httpOnly Cookie)
```bash
curl -i -X POST http://localhost:8080/auth/logout \
  -b cookies.txt \
  -c cookies.txt
```

**Response:**
```json
{
  "message": "Logged out successfully",
  "status": "ok"
}
```

The cookie is cleared by setting `Max-Age=0`.

### Test Invalid Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "wrongpassword"
  }'
```

---

## Protected Endpoints (Require Cookie Authentication)

**Protected endpoints:**
- `POST /llm` - Generate meal suggestions
- `GET /api/profile` - Get user profile
- `GET /api/preferences` - Get user preferences
- `PUT /api/preferences` - Update user preferences
- `GET /api/usage` - Get usage statistics (meal generation count)

### Get User Profile
```bash
# First, login to get cookie (see above)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Then access protected endpoint with cookie
curl http://localhost:8080/api/profile \
  -b cookies.txt
```

The `-b cookies.txt` flag sends the saved cookie with the request.

### Test Unauthorized Access (No Cookie)
```bash
curl http://localhost:8080/api/profile
```

**Response:**
```json
{
  "message": "Authentication required",
  "status": "error"
}
```

### Test After Logout (Cookie Cleared)
```bash
# Logout to clear cookie
curl -X POST http://localhost:8080/auth/logout \
  -b cookies.txt \
  -c cookies.txt

# Try to access protected endpoint (should fail)
curl http://localhost:8080/api/profile \
  -b cookies.txt
```

**Response:**
```json
{
  "message": "Authentication required",
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

### LLM Request (Requires Authentication)
```bash
# First, login to get cookie
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Then make LLM request with cookie
curl -X POST http://localhost:8080/llm \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "message": "Rice, soy sauce, yoghurt, onion, peppers, beansprouts and tofu"
  }'
```

---

## User Preferences

### Get User Preferences
```bash
# Login first
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Get preferences
curl http://localhost:8080/api/preferences \
  -b cookies.txt
```

### Set/Update User Preferences
```bash
# Login first
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Update preferences
curl -X PUT http://localhost:8080/api/preferences \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "dietary_restrictions": "vegetarian,gluten-free",
    "max_cooking_time": 30
  }'
```

### Complete Workflow: Set Preferences + Get Meal Suggestions
```bash
# 1. Login and save cookie
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# 2. Set preferences
curl -X PUT http://localhost:8080/api/preferences \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "dietary_restrictions": "vegetarian",
    "max_cooking_time": 30
  }'

# 3. Get meal suggestions (preferences are automatically included)
curl -X POST http://localhost:8080/llm \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "message": "Rice, tofu, soy sauce, peppers, onions"
  }'
```

---

## Usage Tracking

### Get Usage Statistics
```bash
# Login first
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Get usage stats
curl http://localhost:8080/api/usage \
  -b cookies.txt
```

**Response:**
```json
{
  "success": true,
  "data": {
    "used": 5,
    "remaining": 15,
    "limit": 20
  }
}
```

### Test Usage Limit (20 meals)
```bash
# Login first
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -c cookies.txt

# Make multiple meal requests
for i in {1..21}; do
  echo "Request $i:"
  curl -s -X POST http://localhost:8080/llm \
    -H "Content-Type: application/json" \
    -b cookies.txt \
    -d '{"message": "Test meal generation"}' | jq '.data.usage'
  echo "---"
done

# Request 21 should fail with "Usage limit reached"
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

### Full Authentication Test with httpOnly Cookies
```bash
# 1. Start server (in one terminal)
go run main.go

# 2. Register user and save cookie (in another terminal)
curl -i -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123"}' \
  -c cookies.txt

# 3. Access protected endpoint with cookie
curl http://localhost:8080/api/profile \
  -b cookies.txt

# 4. Try accessing without cookie (should fail)
curl http://localhost:8080/api/profile

# 5. Login with correct credentials
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "test123"}' \
  -c cookies.txt

# 6. Try wrong password
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "wrong"}'

# 7. Logout (clear cookie)
curl -i -X POST http://localhost:8080/auth/logout \
  -b cookies.txt \
  -c cookies.txt

# 8. Try accessing protected endpoint after logout (should fail)
curl http://localhost:8080/api/profile \
  -b cookies.txt
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
