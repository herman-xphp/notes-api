# Notes API

A simple, secure, and production-ready REST API for managing notes with user authentication.

## Overview

Notes API is a complete backend solution for a notes management application. It features user authentication with JWT, comprehensive testing, security best practices, and production-ready code.

**Status**: ✅ Production Ready | **Tests**: 24/24 Passing | **Coverage**: 87.5%

## Features

✅ **User Authentication** - Register, login with JWT tokens  
✅ **Note Management** - CRUD operations with pagination  
✅ **Secure** - XSS prevention, SQL injection mitigation, rate limiting  
✅ **Well-Tested** - 24 unit tests with 87.5% code coverage  
✅ **Documented** - Swagger/OpenAPI specs, comprehensive guides  
✅ **Clean Code** - Clean architecture, proper separation of concerns  
✅ **Production-Ready** - Structured logging, graceful shutdown, error handling  

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.22+ |
| **Framework** | Fiber v2 |
| **Database** | MySQL 8.0+ with GORM |
| **Authentication** | JWT (golang-jwt/jwt/v5) |
| **Password Hash** | Bcrypt |
| **Logging** | Zerolog |
| **Testing** | Testify |
| **Validation** | go-playground/validator |

## Quick Start

### 1. Install & Setup

```bash
# Clone repository
git clone https://github.com/herman-xphp/notes-api.git
cd notes-api

# Copy environment template
cp .env.example .env

# Install dependencies
go mod download

# Start the server
go run cmd/api/main.go
```

Server runs on `http://localhost:3000`

### 2. Try the API

```bash
# Register
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"SecurePass123"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"SecurePass123"}' \
  | jq -r '.data.token')

# Create note
curl -X POST http://localhost:3000/api/v1/notes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Note","content":"Hello World"}'
```

## API Endpoints

### Auth
```
POST   /api/v1/auth/register    - Register new user
POST   /api/v1/auth/login       - Login & get JWT token
```

### Health
```
GET    /api/v1/health           - Health check
```

### Notes (Protected with JWT)
```
POST   /api/v1/notes            - Create note
GET    /api/v1/notes            - Get all notes (paginated)
GET    /api/v1/notes/{id}       - Get note by ID
PUT    /api/v1/notes/{id}       - Update note
DELETE /api/v1/notes/{id}       - Delete note
```

## Documentation

Comprehensive documentation available in these guides:

📖 **[Getting Started](./GETTING_STARTED.md)**
- Installation instructions
- Environment setup
- Quick start examples
- Docker setup

📚 **[API Documentation](./API_DOCUMENTATION.md)**
- Complete endpoint reference
- Request/response formats
- Error codes and messages
- cURL examples

🧪 **[Testing Guide](./TESTING.md)**
- Running tests
- Test coverage details
- Test structure and patterns
- Adding new tests

🚀 **[Deployment Guide](./DEPLOYMENT.md)**
- Binary deployment
- Docker deployment
- Kubernetes setup
- SSL/TLS configuration
- Backup strategies

📋 **[Swagger Docs](./docs/swagger.yaml)**
- Interactive OpenAPI documentation
- Can be viewed at https://editor.swagger.io

## Project Structure

```
notes-api/
├── cmd/api/
│   └── main.go                      # Entry point
├── internal/
│   ├── domain/                      # Domain models (Note, User)
│   ├── dto/                         # Data transfer objects
│   ├── handler/                     # HTTP handlers + tests
│   │   ├── auth_handler.go
│   │   ├── auth_handler_test.go     # 7 tests
│   │   ├── note_handler.go
│   │   ├── note_handler_test.go     # 15 tests
│   │   ├── health_handler.go
│   │   └── health_handler_test.go   # 2 tests
│   ├── middleware/
│   │   ├── auth.go                  # JWT auth
│   │   ├── error_handler.go         # Global error handler
│   │   ├── ratelimit.go             # Rate limiting
│   │   └── ...
│   ├── repository/                  # Data access layer
│   ├── service/                     # Business logic
│   └── utils/                       # JWT, hashing, logging
├── pkg/database/
│   └── mysql.go                     # DB connection
├── configs/
│   └── config.go                    # Config management
├── migrations/                      # Database migrations
├── docs/                            # Generated Swagger docs
└── test.http                        # HTTP test file
```

## Security Features

| Feature | Implementation |
|---------|-----------------|
| **XSS Prevention** | HTML entity escaping on input |
| **SQL Injection** | GORM parameterized queries |
| **Password Security** | Bcrypt hashing, strength validation |
| **Authentication** | JWT with 24-hour expiration |
| **Authorization** | Per-user resource access control |
| **Rate Limiting** | Auth: 5/15min, General: 10/1min |
| **CORS** | Configurable allowed origins |
| **Security Headers** | HSTS, Content-Type, X-Frame-Options |
| **Request Limits** | Timeout, body size, etc |

## Testing

### Run Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Specific package
go test -v ./internal/handler
```

### Test Coverage

- **24 Total Tests** - All passing ✅
- **87.5% Coverage** - Handler package
- **7 Auth Tests** - Register, login, errors
- **2 Health Tests** - Check status, response
- **15 Note Tests** - CRUD operations, auth, validation

## Configuration

Create `.env` file from `.env.example`:

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

# JWT (min 32 chars)
JWT_SECRET=your-secret-key-minimum-32-characters-long

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## Build & Run

### Development

```bash
go run cmd/api/main.go
```

### Build Binary

```bash
go build -o notes-api ./cmd/api/main.go
./notes-api
```

### Docker

```bash
docker build -t notes-api .
docker run -p 3000:3000 notes-api
```

## Deployment

See [Deployment Guide](./DEPLOYMENT.md) for:
- System service setup
- Docker Compose
- Kubernetes
- SSL/TLS configuration
- Database migration
- Monitoring setup

## Code Quality

✅ **Clean Architecture** - Clear separation of concerns  
✅ **Interface-Based** - Easy to mock and test  
✅ **Error Handling** - Structured AppError types  
✅ **Input Validation** - go-playground/validator  
✅ **Logging** - Structured with Zerolog  
✅ **Security** - Multiple layers of protection  

## Performance

- **Connection Pooling** - 10-100 idle/open connections
- **Request Timeout** - 10 second limit
- **Body Size Limit** - 4 MB maximum
- **Rate Limiting** - Prevent abuse
- **Graceful Shutdown** - 10 second timeout for cleanup

## Monitoring & Logs

### Health Endpoint

```bash
curl http://localhost:3000/api/v1/health
```

### Structured Logs

All logs output as JSON in production:
```json
{
  "level": "info",
  "time": "2024-11-17T10:30:00Z",
  "message": "Server running",
  "port": "3000"
}
```

## Future Enhancements

- [ ] Integration tests
- [ ] Database transactions
- [ ] Soft delete support
- [ ] Refresh tokens
- [ ] Search/Filter functionality
- [ ] Export notes (PDF, Markdown)
- [ ] Real-time updates (WebSocket)
- [ ] File attachments
- [ ] Tags/Categories
- [ ] Sharing notes

## Troubleshooting

### Database Connection Failed
- Check MySQL is running
- Verify credentials in `.env`
- Ensure database exists

### Port Already in Use
- Change `APP_PORT` in `.env`
- Or: `kill $(lsof -t -i :3000)`

### Tests Failing
- Run `go mod tidy`
- Clear cache: `go clean -cache`
- Ensure environment is isolated

See [Deployment Guide](./DEPLOYMENT.md#troubleshooting) for more issues.

## Development Tips

### Code Formatting
```bash
go fmt ./...
```

### Run Linter
```bash
golangci-lint run ./...
```

### Generate Swagger Docs
```bash
swag init -g cmd/api/main.go
```

### View Test Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## API Response Format

All responses follow this format:

```json
{
  "status": "success|error",
  "message": "Descriptive message",
  "data": { ... }
}
```

Success: `status: "success"`  
Error: `status: "error"`  

## Contributing

1. Fork repository
2. Create feature branch
3. Make changes with tests
4. Run `go test ./...`
5. Submit pull request

## License

Apache License 2.0 - See [LICENSE](./LICENSE) file

## Support

- 📖 **Documentation** - See guides above
- 🐛 **Issues** - Report on GitHub
- 💬 **Discussions** - Use GitHub Discussions
- 📧 **Email** - Check repository for contact

## Author

**Herman** - Senior Backend Developer

- Clean architecture enthusiast
- Go expertise
- RESTful API design
- Security-first mindset

---

**Built with ❤️ for simplicity, security, and scalability**
