# API Documentation

Complete API reference for Notes API.

## Base URL

```
http://localhost:3000/api/v1
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

Tokens are obtained by logging in or registering.

## Response Format

All responses follow this format:

```json
{
  "status": "success|error",
  "message": "descriptive message",
  "data": {}
}
```

## Endpoints

### Authentication

#### Register User

Register a new user account.

```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "message": "user registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or validation failed
- `409 Conflict` - Email already exists
- `500 Internal Server Error` - Server error

**Password Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number

---

#### Login

Authenticate and get a JWT token.

```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "login successful",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Invalid email or password
- `500 Internal Server Error` - Server error

**Token Details:**
- Format: JWT (JSON Web Token)
- Expiration: 24 hours
- Algorithm: HS256

---

### Health Check

#### Check API Health

Check if the API and database are healthy.

```http
GET /health
```

**Success Response (200):**
```json
{
  "status": "OK"
}
```

**Error Response (503):**
```json
{
  "status": "unhealthy",
  "message": "database ping failed"
}
```

---

### Notes

#### Create Note

Create a new note.

```http
POST /notes
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Note",
  "content": "Note content here"
}
```

**Parameters:**
- `title` (string, required) - Note title
- `content` (string, optional) - Note content

**Success Response (201):**
```json
{
  "status": "success",
  "message": "note created",
  "data": {
    "id": 1,
    "title": "My Note",
    "content": "Note content here",
    "created_at": "2024-11-17T10:30:00Z",
    "updated_at": "2024-11-17T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation failed
- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Server error

---

#### Get All Notes

Retrieve all notes for the authenticated user with pagination.

```http
GET /notes?page=1&limit=10
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (integer, default: 1) - Page number
- `limit` (integer, default: 10) - Items per page

**Success Response (200):**
```json
{
  "status": "success",
  "message": "notes fetched",
  "data": {
    "data": [
      {
        "id": 1,
        "title": "First Note",
        "content": "Content 1",
        "created_at": "2024-11-17T10:30:00Z",
        "updated_at": "2024-11-17T10:30:00Z"
      },
      {
        "id": 2,
        "title": "Second Note",
        "content": "Content 2",
        "created_at": "2024-11-17T10:35:00Z",
        "updated_at": "2024-11-17T10:35:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid pagination params
- `401 Unauthorized` - Missing or invalid token

---

#### Get Note by ID

Retrieve a specific note.

```http
GET /notes/{id}
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer) - Note ID

**Success Response (200):**
```json
{
  "status": "success",
  "message": "note fetched",
  "data": {
    "id": 1,
    "title": "My Note",
    "content": "Note content here",
    "created_at": "2024-11-17T10:30:00Z",
    "updated_at": "2024-11-17T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found

---

#### Update Note

Update a note (partial update supported).

```http
PUT /notes/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

**Path Parameters:**
- `id` (integer) - Note ID

**Request Body:**
- `title` (string, optional) - New title
- `content` (string, optional) - New content
- At least one field must be provided

**Success Response (200):**
```json
{
  "status": "success",
  "message": "note updated",
  "data": {
    "id": 1,
    "title": "Updated Title",
    "content": "Updated content",
    "created_at": "2024-11-17T10:30:00Z",
    "updated_at": "2024-11-17T11:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request or validation failed
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found

---

#### Delete Note

Delete a note.

```http
DELETE /notes/{id}
Authorization: Bearer <token>
```

**Path Parameters:**
- `id` (integer) - Note ID

**Success Response (200):**
```json
{
  "status": "success",
  "message": "note deleted"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid note ID
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Note not found

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request or validation failed |
| 401 | Unauthorized - Authentication failed or missing token |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Database connection failed |

## Rate Limiting

Rate limits are applied to prevent abuse:

- **Auth Endpoints** (`/auth/register`, `/auth/login`): 5 requests per 15 minutes
- **General Endpoints**: 10 requests per 1 minute

When rate limit exceeded: `429 Too Many Requests`

## Error Messages

Common error messages:

```json
{
  "status": "error",
  "message": "invalid email or password"
}
```

```json
{
  "status": "error",
  "message": "email already exists"
}
```

```json
{
  "status": "error",
  "message": "password must be at least 8 characters long and contain uppercase, lowercase, and number"
}
```

```json
{
  "status": "error",
  "message": "unauthorized"
}
```

```json
{
  "status": "error",
  "message": "resource not found"
}
```

## Examples Using cURL

### Register and Get Token

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123"
  }' | jq -r '.data.token')

echo $TOKEN
```

### Create Note with Token

```bash
curl -X POST http://localhost:3000/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "API Test",
    "content": "Testing the API"
  }'
```

### Get All Notes

```bash
curl http://localhost:3000/api/v1/notes \
  -H "Authorization: Bearer $TOKEN"
```

## OpenAPI/Swagger Documentation

Interactive API documentation is available at:
- **Format**: OpenAPI 2.0 (Swagger)
- **Location**: `/docs/swagger.yaml` or `/docs/swagger.json`
- **Viewer**: https://editor.swagger.io

## CORS

CORS is enabled for specified origins. Configure in `.env`:

```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## Content Types

- Request: `application/json`
- Response: `application/json`

All requests must include `Content-Type: application/json` header.
