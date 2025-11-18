# Getting Started with Notes API

A simple and secure REST API for managing notes with user authentication.

## Prerequisites

- Go 1.20+
- MySQL 8.0+
- Docker & Docker Compose (optional)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/herman-xphp/notes-api.git
cd notes-api
```

### 2. Setup Environment Variables

Copy the example environment file and update with your database credentials:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database
DB_HOST=localhost
DB_USER=root
DB_PASS=your_password
DB_NAME=notes_api
DB_PORT=3306

# API
APP_PORT=3000
ENV=development

# JWT
JWT_SECRET=your-secret-key-min-32-characters-long

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Install Dependencies

```bash
go mod download
go mod tidy
```

### 4. Run Database Migrations

```bash
# Using golang-migrate
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/notes_api" up

# Or using GORM auto-migration (if configured)
go run cmd/api/main.go
```

### 5. Build and Run

```bash
# Build
go build -o notes-api ./cmd/api/main.go

# Run
./notes-api
```

Server will start on `http://localhost:3000`

## Quick Start Example

### 1. Register a User

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123"
  }'
```

Response:
```json
{
  "status": "success",
  "message": "user registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "token": "eyJhbGc..."
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123"
  }'
```

### 3. Create a Note

```bash
curl -X POST http://localhost:3000/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "My First Note",
    "content": "This is a sample note content"
  }'
```

### 4. Get All Notes

```bash
curl http://localhost:3000/api/v1/notes \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json"
```

### 5. Get Note by ID

```bash
curl http://localhost:3000/api/v1/notes/1 \
  -H "Authorization: Bearer <your_token>"
```

### 6. Update a Note

```bash
curl -X PUT http://localhost:3000/api/v1/notes/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "Updated Title"
  }'
```

### 7. Delete a Note

```bash
curl -X DELETE http://localhost:3000/api/v1/notes/1 \
  -H "Authorization: Bearer <your_token>"
```

## Docker Setup (Optional)

### Run with Docker Compose

```bash
docker-compose up -d
```

This will start both MySQL and the API server.

## Documentation

- **[API Documentation](./API_DOCUMENTATION.md)** - Complete API reference
- **[Testing Guide](./TESTING.md)** - How to run tests
- **[Deployment Guide](./DEPLOYMENT.md)** - Production deployment
- **[Swagger/OpenAPI](./docs/swagger.yaml)** - Interactive API docs

## Project Structure

```
notes-api/
├── cmd/api/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/                 # Domain models
│   ├── dto/                    # Data transfer objects
│   ├── handler/                # HTTP handlers + tests
│   ├── middleware/             # Middleware (auth, error handling, etc)
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic
│   └── utils/                  # Utilities (JWT, hashing, logging)
├── pkg/
│   └── database/               # Database connection
├── configs/                    # Configuration management
├── migrations/                 # Database migrations
└── docs/                       # API documentation (Swagger)
```

## Security Features

✅ **XSS Prevention** - HTML input escaping  
✅ **SQL Injection Prevention** - GORM parameterized queries  
✅ **Password Security** - Bcrypt hashing with validation  
✅ **JWT Authentication** - 24-hour token expiration  
✅ **Rate Limiting** - Auth endpoints (5/15min), General (10/1min)  
✅ **CORS Protection** - Configurable allowed origins  
✅ **Security Headers** - HSTS, Content-Type, X-Frame-Options  

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestCreateNote_Success ./internal/handler
```

### Code Quality

- **Clean Architecture** - Separation of concerns
- **Interface-based Design** - Easy to mock and test
- **Error Handling** - Structured error responses
- **Logging** - Structured logging with Zerolog
- **Validation** - Input validation with go-playground/validator

## Troubleshooting

### Database Connection Failed

- Check MySQL is running: `mysql -u root -p`
- Verify credentials in `.env`
- Check database exists: `CREATE DATABASE notes_api;`

### Port Already in Use

Change `APP_PORT` in `.env` to an available port.

### JWT Token Expired

- Get a new token by logging in again
- Token expiration is set to 24 hours

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review test cases for usage examples

## License

Apache License 2.0 - See LICENSE file
