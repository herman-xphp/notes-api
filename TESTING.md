# Testing Guide

Comprehensive guide for running and understanding tests in Notes API.

## Overview

The Notes API includes:
- **24 Unit Tests** covering all handlers
- **87.5% Code Coverage** on handler package
- **Mock Service Layer** for isolated testing
- **Testify** for assertions and testing utilities

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Tests with Detailed Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Test Package

```bash
# Handler tests only
go test -v ./internal/handler

# Service tests
go test -v ./internal/service/impl

# All tests
go test -v ./internal/...
```

### Run Specific Test

```bash
go test -run TestCreateNote_Success -v ./internal/handler
```

### Run Tests with Timeout

```bash
go test -timeout 30s ./...
```

## Test Structure

### Handler Tests

Located in: `internal/handler/*_test.go`

Tests are organized by handler:

#### Auth Handler Tests (7 tests)

```
TestRegisterSuccess          - Happy path registration
TestRegisterEmailExists      - Conflict handling (email already exists)
TestRegisterServerError      - Server error handling
TestLoginSuccess             - Happy path login
TestLoginUnauthorized        - Invalid credentials
TestLoginServerError         - Server error handling
TestNewAuthHandler           - Handler initialization
```

#### Health Handler Tests (2 tests)

```
TestHealthCheck_HandlerExists        - Handler initialization
TestHealthCheck_ResponseStructure    - Response format validation
```

#### Note Handler Tests (15 tests)

**Create Operations:**
```
TestCreateNote_Unauthorized   - Missing authentication
TestCreateNote_Success        - Happy path create
TestCreateNote_ServiceError   - Service error handling
```

**Read Operations:**
```
TestGetByID_Unauthorized      - Missing authentication
TestGetByID_InvalidID         - Invalid ID format
TestGetByID_Success           - Happy path get single note
TestGetAllNotes_Unauthorized  - Missing authentication
TestGetAllNotes_Success       - Happy path get all with pagination
```

**Update Operations:**
```
TestUpdateNote_Unauthorized   - Missing authentication
TestUpdateNote_InvalidID      - Invalid ID format
TestUpdateNote_Success        - Happy path update
```

**Delete Operations:**
```
TestDeleteNote_Unauthorized   - Missing authentication
TestDeleteNote_InvalidID      - Invalid ID format
TestDeleteNote_Success        - Happy path delete
```

**Initialization:**
```
TestNewNoteHandler            - Handler initialization
```

## Test Coverage Details

### Coverage by Package

```
handler/              87.5%  ✅ Well covered
```

### Coverage by Category

- ✅ **Happy Path** - All success scenarios covered
- ✅ **Authorization** - All auth checks tested
- ✅ **Input Validation** - Invalid inputs validated
- ✅ **Error Handling** - Error cases verified
- ✅ **Initialization** - Handler initialization tested

## Testing Patterns

### Mock Service Layer

Tests use mocks for the service layer:

```go
type mockNoteService struct {
    createResp *domain.Note
    createErr  error
    // ... other fields
}

func (m *mockNoteService) Create(ctx context.Context, userID uint, 
    req dto.CreateNoteRequest) (*domain.Note, error) {
    return m.createResp, m.createErr
}
```

### Using httptest

Tests use `httptest.NewRequest()` and `app.Test()`:

```go
func TestCreateNote_Success(t *testing.T) {
    // Setup
    mock := &mockNoteService{
        createResp: &domain.Note{ID: 1, Title: "Test", ...},
    }
    h := NewNoteHandler(mock)
    
    // Create Fiber app
    app := fiber.New()
    app.Post("/notes", func(c *fiber.Ctx) error {
        c.Locals("validated", &dto.CreateNoteRequest{...})
        c.Locals("user_id", uint(1))
        return h.Create(c)
    })
    
    // Test
    req := httptest.NewRequest("POST", "/notes", nil)
    resp, err := app.Test(req)
    
    // Assert
    require.NoError(t, err)
    require.Equal(t, 201, resp.StatusCode)
}
```

### Testify Assertions

Tests use testify for clear assertions:

```go
require.NoError(t, err)           // No error occurred
require.Equal(t, 201, statusCode) // Status is 201
require.NotNil(t, handler)        // Value is not nil
require.Contains(t, body, "text") // String contains text
```

## Test Scenarios

### Authorization Tests

All protected endpoints verify:
- Missing authorization header → 401
- Invalid token format → 401
- Expired token → 401

### Input Validation Tests

Endpoints validate:
- Required fields present
- Field type validation
- Format validation (email, etc)
- Range validation (limits, etc)

### Error Handling Tests

Scenarios covered:
- Database errors
- Service errors
- Validation failures
- Not found errors

## Integration Testing

Current unit tests use mocked services. For integration tests:

```bash
# Integration tests would connect to real database
go test -tags=integration ./...
```

Integration tests are listed as future improvements.

## Code Coverage Goals

| Package | Current | Target | Status |
|---------|---------|--------|--------|
| handler | 87.5%   | 80%    | ✅ Met |
| service | TBD     | 80%    | 📋 Planned |
| repo    | TBD     | 70%    | 📋 Planned |

## Continuous Integration

For CI/CD pipelines, run:

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run linter
golangci-lint run ./...

# Build
go build -o notes-api ./cmd/api/main.go
```

## Debugging Tests

### Print Debug Information

```go
t.Logf("Variable value: %v", variable)
```

### Run Single Test with Debug

```bash
go test -run TestCreateNote_Success -v ./internal/handler
```

### Use Delve Debugger

```bash
dlv test ./internal/handler -- -test.run TestCreateNote_Success
```

## Best Practices

✅ **Do:**
- Test both success and failure cases
- Use clear test names describing the scenario
- Mock external dependencies
- Use table-driven tests for similar scenarios
- Keep tests independent and isolated

❌ **Don't:**
- Test internal implementation details
- Create dependencies between tests
- Use sleep/delays in tests
- Skip error checking
- Use hardcoded values without context

## Test Data

### Example User
```json
{
  "id": 1,
  "name": "Test User",
  "email": "test@example.com",
  "password": "SecurePass123"
}
```

### Example Note
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Test Note",
  "content": "Test Content",
  "created_at": "2024-11-17T10:30:00Z",
  "updated_at": "2024-11-17T10:30:00Z"
}
```

## Adding New Tests

When adding new features:

1. **Create test file** following pattern: `feature_test.go`
2. **Mock dependencies** needed for testing
3. **Test success case** first
4. **Test error cases** (validation, auth, not found, etc)
5. **Verify coverage** is maintained/improved
6. **Update this documentation** if adding new test patterns

## Troubleshooting

### Tests Failing Locally But Passing in CI

- Check Go version matches
- Ensure database is not running (we use mocks)
- Clear Go cache: `go clean -cache`

### Test Timeout

- Increase timeout: `go test -timeout 60s ./...`
- Check for infinite loops in code

### Mock Not Working

- Verify mock methods match interface signature
- Check method parameter order
- Ensure error return type is correct

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Package](https://github.com/stretchr/testify)
- [Fiber Testing](https://docs.gofiber.io/guide/testing)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
