# Swagger UI Quick Guide

## 🚀 Quick Start

### Access Swagger UI

```
http://localhost:8080/api/v1/swagger
```

## 📖 Using Swagger UI

### 1. Testing Endpoints

#### Without Authentication

1. **Find the endpoint** you want to test (e.g., `POST /users`)
2. **Click "Try it out"**
3. **Fill in the request body**:
```json
{
  "email": "test@example.com",
  "name": "Test User",
  "password": "password123"
}
```
4. **Click "Execute"**
5. **View the response**

#### With Authentication (Protected Endpoints)

1. **First, login to get a token**:
   - Go to `POST /auth/login`
   - Click "Try it out"
   - Enter credentials:
   ```json
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```
   - Click "Execute"
   - **Copy the token** from response

2. **Authorize with the token**:
   - Click the **"Authorize"** button (🔒 icon at top)
   - Enter: `Bearer YOUR_TOKEN_HERE`
   - Click "Authorize"
   - Click "Close"

3. **Test protected endpoints**:
   - Now you can test endpoints like `GET /users`, `GET /users/{id}`, etc.
   - The token will be automatically included in requests

### 2. Viewing Schemas

**Request Schemas**:
- Scroll to "Schemas" section
- Click on schema name (e.g., `CreateUserRequest`)
- View required fields, types, and validation rules

**Response Schemas**:
- Each endpoint shows response schema
- Example values are provided
- Click "Model" to see structure

### 3. Understanding Responses

**HTTP Status Codes**:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid token)
- `404` - Not Found
- `409` - Conflict (e.g., email already exists)

**Response Body**:
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com",
    "name": "Test User",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🔐 Authentication Flow

### Complete Flow Example

```bash
# 1. Create a user
POST /api/v1/users
{
  "email": "newuser@example.com",
  "name": "New User",
  "password": "securepass123"
}

# Response: 201 Created
{
  "success": true,
  "message": "User created successfully",
  "data": { ... }
}

# 2. Login
POST /api/v1/auth/login
{
  "email": "newuser@example.com",
  "password": "securepass123"
}

# Response: 200 OK
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { ... }
  }
}

# 3. Use token for protected endpoints
GET /api/v1/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Response: 200 OK
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

## 📝 Common Use Cases

### Use Case 1: Create and Get User

```
1. POST /users - Create user
   → Get user ID from response

2. POST /auth/login - Login
   → Copy token

3. Click "Authorize" → Paste token

4. GET /users/{id} - Get the user
   → Paste user ID in path parameter
   → Execute
```

### Use Case 2: Update User

```
1. Login and authorize (steps above)

2. PUT /users/{id}
   → Enter user ID
   → Enter request body:
   {
     "name": "Updated Name"
   }
   → Execute
```

### Use Case 3: List Users with Pagination

```
1. Login and authorize

2. GET /users
   → Set query parameters:
     - page: 1
     - limit: 10
   → Execute

3. View paginated response with meta info
```

## 🎨 Swagger UI Features

### Collapsible Sections

- **Click tag name** (e.g., "Users") to expand/collapse
- **Click endpoint** to expand details
- **Click schema** to view model

### Filter Endpoints

- Use the **search box** at top
- Filter by tag, path, or description

### Download Spec

- Click **"OpenAPI spec"** link at top
- Download as JSON or YAML
- Use for code generation or client SDKs

### Try Different Servers

- Select server from dropdown
- Useful for testing dev vs prod

## 🔧 Troubleshooting

### "Unauthorized" Error

**Problem**: Getting 401 when testing protected endpoints

**Solution**:
1. Make sure you've clicked "Authorize"
2. Check token format: must be `Bearer YOUR_TOKEN`
3. Token might be expired - login again

### "Bad Request" Error

**Problem**: Getting 400 when testing endpoints

**Solution**:
1. Check required fields are filled
2. Validate email format
3. Check min/max length requirements
4. View error response for details

### Can't See Response

**Problem**: Response not showing after Execute

**Solution**:
1. Check browser console for errors
2. Make sure server is running
3. Check CORS configuration
4. Try refreshing the page

## 📚 Advanced Tips

### 1. Import to Postman

1. Click "OpenAPI spec" link
2. Copy URL: `http://localhost:8080/api/v1/swagger/spec`
3. In Postman: Import → Link → Paste URL
4. All endpoints imported automatically!

### 2. Generate Client SDK

```bash
# TypeScript
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8080/api/v1/swagger/spec \
  -g typescript-axios \
  -o ./clients/typescript

# Python
openapi-generator-cli generate \
  -i http://localhost:8080/api/v1/swagger/spec \
  -g python \
  -o ./clients/python

# Go
oapi-codegen -package client \
  http://localhost:8080/api/v1/swagger/spec \
  > ./clients/go/client.go
```

### 3. Test with curl

Copy the curl command from Swagger UI:
1. Execute a request in Swagger
2. Find "curl" tab in response
3. Copy the command
4. Run in terminal

Example:
```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/users' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@example.com",
    "name": "User",
    "password": "pass123"
  }'
```

## 🎓 Learning More

### OpenAPI Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger Editor](https://editor.swagger.io/)
- [OpenAPI Guide](https://oai.github.io/Documentation/)

### Best Practices

1. **Always test in order**:
   - Create → Login → Use token → Test protected endpoints

2. **Use realistic data**:
   - Valid email formats
   - Strong passwords
   - Meaningful names

3. **Check validation rules**:
   - Min/max length
   - Required fields
   - Format requirements

4. **Save successful responses**:
   - Copy IDs for later use
   - Save tokens for testing
   - Note success patterns

## 🔗 Quick Links

- **Swagger UI**: http://localhost:8080/api/v1/swagger
- **OpenAPI Spec**: http://localhost:8080/api/v1/swagger/spec
- **API Workflow**: [API_FIRST_WORKFLOW.md](API_FIRST_WORKFLOW.md)
- **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)

---

**Happy Testing!** 🚀

For issues or questions, see [CONTRIBUTING.md](../CONTRIBUTING.md)
