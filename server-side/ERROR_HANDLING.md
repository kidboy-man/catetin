# Custom Error Handling System

## Overview

The application uses a custom error system that provides consistent error handling across all layers with proper HTTP status codes and error codes.

## Architecture

### 1. Custom Error Type (`pkg/errors/errors.go`)

**AppError** - A custom error type that includes:
- `Code` - Unique error code for identifying error types
- `Message` - Human-readable error message
- `HTTPStatus` - HTTP status code for API responses
- `Details` - Additional context (optional)
- `Err` - Underlying error for wrapping (optional)

**Features:**
- ✅ Compatible with native Go `error` interface
- ✅ Supports error wrapping with `Unwrap()`
- ✅ Can add additional details dynamically
- ✅ Automatic HTTP status code mapping

### 2. Error Codes

Predefined error codes for common scenarios:

#### General Errors
- `INTERNAL_ERROR` - Internal server error (500)
- `BAD_REQUEST` - Invalid request (400)
- `UNAUTHORIZED` - Unauthorized access (401)
- `FORBIDDEN` - Access forbidden (403)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Resource conflict (409)
- `VALIDATION_ERROR` - Validation failed (400)

#### Authentication Errors
- `INVALID_CREDENTIALS` - Invalid email/password (401)
- `EMAIL_ALREADY_EXISTS` - Email already registered (409)
- `INVALID_TOKEN` - Invalid auth token (401)
- `EXPIRED_TOKEN` - Expired auth token (401)

#### Resource Errors
- `USER_NOT_FOUND` - User not found (404)
- `RESOURCE_NOT_FOUND` - Generic resource not found (404)
- `VERSION_CONFLICT` - Optimistic locking conflict (409)

#### Business Logic Errors
- `INVALID_INPUT` - Invalid input provided (400)
- `OPERATION_NOT_ALLOWED` - Operation not allowed (403)

### 3. Error Handler Middleware

Located at `internal/controller/http/middleware/error_handler.go`

**Features:**
- Automatically converts `AppError` to proper HTTP responses
- Catches unhandled errors and returns 500
- Logs unexpected errors
- Returns standardized JSON error format

**Helper Functions:**
- `AbortWithError(c, err)` - Abort request with any error
- `AbortWithAppError(c, appErr)` - Abort request with AppError

## Usage Guide

### Creating Errors

#### 1. Use Predefined Errors
```go
import appErrors "github.com/ingunawandra/catetin/pkg/errors"

// Return predefined error
return nil, appErrors.ErrInvalidCredentials

// Add details to predefined error
return nil, appErrors.ErrValidation.WithDetails(map[string]interface{}{
    "field": "email",
    "reason": "invalid format",
})
```

#### 2. Create New Error
```go
// Create a new custom error
return nil, appErrors.New(
    appErrors.ErrCodeInvalidInput,
    "Amount must be greater than zero",
    http.StatusBadRequest,
)
```

#### 3. Wrap Existing Error
```go
// Wrap an existing error with context
user, err := repo.FindByID(ctx, userID)
if err != nil {
    return nil, appErrors.Wrap(
        err,
        appErrors.ErrCodeInternal,
        "Failed to find user",
        http.StatusInternalServerError,
    )
}
```

### Service Layer Example

```go
package service

import (
    appErrors "github.com/ingunawandra/catetin/pkg/errors"
)

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
    // Get user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            // Return custom error for invalid credentials
            return nil, appErrors.ErrInvalidCredentials
        }
        // Wrap unexpected errors
        return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to find user", 500)
    }

    // Verify password
    if !s.passwordHasher.IsValidPassword(user.Password, password) {
        return nil, appErrors.ErrInvalidCredentials
    }

    // Success
    return &LoginResponse{User: user}, nil
}
```

### HTTP Handler Example

```go
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/ingunawandra/catetin/internal/controller/http/middleware"
    appErrors "github.com/ingunawandra/catetin/pkg/errors"
)

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest

    // Validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        // Use middleware helper to abort with error
        middleware.AbortWithAppError(c, appErrors.ErrValidation.WithDetails(map[string]interface{}{
            "validation_errors": err.Error(),
        }))
        return
    }

    // Call service
    result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        // Error middleware will automatically handle this
        middleware.AbortWithError(c, err)
        return
    }

    // Return success response
    c.JSON(http.StatusOK, dto.NewSuccessResponse("Login successful", result))
}
```

## Error Response Format

### Success Response
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": {
    // response data
  }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Human-readable error message",
  "errors": {
    "code": "ERROR_CODE",
    // additional details if provided
  }
}
```

## Examples

### Example 1: Email Already Exists (409)
```json
{
  "status": "error",
  "message": "Email already registered",
  "errors": {
    "code": "EMAIL_ALREADY_EXISTS"
  }
}
```

### Example 2: Invalid Credentials (401)
```json
{
  "status": "error",
  "message": "Invalid email or password",
  "errors": {
    "code": "INVALID_CREDENTIALS"
  }
}
```

### Example 3: Validation Error with Details (400)
```json
{
  "status": "error",
  "message": "Validation failed",
  "errors": {
    "code": "VALIDATION_ERROR",
    "validation_errors": "Field 'email' is required"
  }
}
```

### Example 4: Internal Error (500)
```json
{
  "status": "error",
  "message": "An internal error occurred",
  "errors": {
    "code": "INTERNAL_ERROR"
  }
}
```

## Benefits

1. **Consistency**: All errors follow the same format across the entire API
2. **Type Safety**: Error codes are typed constants, preventing typos
3. **Automatic HTTP Mapping**: No need to manually set HTTP status codes in handlers
4. **Error Wrapping**: Maintain error context while adding custom messages
5. **Flexible Details**: Add additional context to errors dynamically
6. **Clean Handlers**: Controllers don't need complex error checking logic
7. **Centralized Logic**: Error handling logic is in one place (middleware)
8. **Go-Compatible**: Works seamlessly with standard Go error handling

## Best Practices

### ✅ DO:

1. **Use predefined errors when possible**
   ```go
   return nil, appErrors.ErrInvalidCredentials
   ```

2. **Add details for context**
   ```go
   return nil, appErrors.ErrValidation.WithDetails(map[string]interface{}{
       "field": "amount",
       "min": 0,
   })
   ```

3. **Wrap errors with context**
   ```go
   return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to process payment", 500)
   ```

4. **Use middleware helpers in handlers**
   ```go
   middleware.AbortWithError(c, err)
   ```

5. **Return appropriate error codes**
   - 400 for validation errors
   - 401 for authentication errors
   - 403 for authorization errors
   - 404 for not found errors
   - 409 for conflicts
   - 500 for internal errors

### ❌ DON'T:

1. **Don't create errors with wrong status codes**
   ```go
   // BAD: Using 200 for an error
   appErrors.New(appErrors.ErrCodeInvalidInput, "Bad input", 200)
   ```

2. **Don't expose sensitive information in error messages**
   ```go
   // BAD: Exposing internal details
   return nil, appErrors.New(..., fmt.Sprintf("SQL Error: %v", err), 500)

   // GOOD: Generic message
   return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to process request", 500)
   ```

3. **Don't manually set HTTP status in handlers**
   ```go
   // BAD: Manually handling errors in controller
   if err != nil {
       c.JSON(http.StatusInternalServerError, ...)
       return
   }

   // GOOD: Let middleware handle it
   if err != nil {
       middleware.AbortWithError(c, err)
       return
   }
   ```

4. **Don't return generic errors**
   ```go
   // BAD
   return nil, errors.New("something went wrong")

   // GOOD
   return nil, appErrors.New(appErrors.ErrCodeInternal, "Failed to process request", 500)
   ```

## Adding New Error Codes

To add a new error code:

1. Add the constant in `pkg/errors/errors.go`:
   ```go
   const (
       // ... existing codes
       ErrCodePaymentFailed ErrorCode = "PAYMENT_FAILED"
   )
   ```

2. Create a predefined error:
   ```go
   var (
       // ... existing errors
       ErrPaymentFailed = New(
           ErrCodePaymentFailed,
           "Payment processing failed",
           http.StatusPaymentRequired,
       )
   )
   ```

3. Use it in your service:
   ```go
   return nil, appErrors.ErrPaymentFailed
   ```

## Testing

When testing error handling:

```go
func TestLogin_InvalidCredentials(t *testing.T) {
    // ... setup

    _, err := authService.Login(ctx, "test@example.com", "wrong")

    // Check if it's the right error
    assert.NotNil(t, err)

    appErr, ok := appErrors.IsAppError(err)
    assert.True(t, ok)
    assert.Equal(t, appErrors.ErrCodeInvalidCredentials, appErr.Code)
    assert.Equal(t, http.StatusUnauthorized, appErr.HTTPStatus)
}
```

## Summary

The custom error handling system provides:
- ✅ Consistent error responses
- ✅ Type-safe error codes
- ✅ Automatic HTTP status mapping
- ✅ Clean separation of concerns
- ✅ Go-idiomatic error handling
- ✅ Easy to extend and maintain
