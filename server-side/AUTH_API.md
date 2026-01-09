# Authentication API Documentation

## Overview
RESTful authentication API with email/password registration and login.

## Base URL
```
http://localhost:8080/api/v1
```

## Endpoints

### 1. Health Check
Check if the API is running.

**Endpoint**: `GET /health`

**Response**:
```json
{
  "status": "healthy",
  "service": "catetin-api"
}
```

---

### 2. Register
Register a new user account.

**Endpoint**: `POST /api/v1/authentications/register`

**Request Body**:
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Validation Rules**:
- `full_name`: Required, minimum 2 characters, maximum 100 characters
- `email`: Required, valid email format
- `password`: Required, minimum 6 characters, maximum 100 characters

**Success Response** (201 Created):
```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "John Doe",
      "email": "john.doe@example.com",
      "phone_number": "john.doe@example.com",
      "image": null
    }
  }
}
```

**Error Responses**:

- **400 Bad Request** - Invalid request payload
```json
{
  "status": "error",
  "message": "Invalid request payload",
  "errors": "validation error details"
}
```

- **409 Conflict** - Email already registered
```json
{
  "status": "error",
  "message": "Email already registered"
}
```

- **500 Internal Server Error** - Server error
```json
{
  "status": "error",
  "message": "Failed to register user",
  "errors": "error details"
}
```

---

### 3. Login
Authenticate with email and password.

**Endpoint**: `POST /api/v1/authentications/login`

**Request Body**:
```json
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

**Validation Rules**:
- `email`: Required, valid email format
- `password`: Required

**Success Response** (200 OK):
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "full_name": "John Doe",
      "email": "john.doe@example.com",
      "phone_number": "john.doe@example.com",
      "image": null
    }
  }
}
```

**Error Responses**:

- **400 Bad Request** - Invalid request payload
```json
{
  "status": "error",
  "message": "Invalid request payload",
  "errors": "validation error details"
}
```

- **401 Unauthorized** - Invalid credentials
```json
{
  "status": "error",
  "message": "Invalid email or password"
}
```

- **500 Internal Server Error** - Server error
```json
{
  "status": "error",
  "message": "Failed to login",
  "errors": "error details"
}
```

---

## Token Information

### Access Token
- **Purpose**: Used to authenticate API requests
- **Default Expiration**: 60 minutes (configurable via `JWT_ACCESS_TOKEN_DURATION`)
- **Usage**: Include in `Authorization` header as `Bearer <token>`

### Refresh Token
- **Purpose**: Used to obtain new access tokens without re-login
- **Default Expiration**: 30 days (configurable via `JWT_REFRESH_TOKEN_DURATION`)
- **Note**: Refresh token endpoint not yet implemented

---

## Testing with cURL

### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/authentications/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/authentications/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Health check
```bash
curl http://localhost:8080/health
```

---

## Testing with Postman

1. **Import Collection**: Create a new collection in Postman
2. **Set Base URL**: Create an environment variable `base_url = http://localhost:8080`
3. **Create Requests**:
   - **Register**: POST `{{base_url}}/api/v1/authentications/register`
   - **Login**: POST `{{base_url}}/api/v1/authentications/login`
   - **Health**: GET `{{base_url}}/health`

---

## Configuration

Required environment variables in `.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=catetin
DB_SSLMODE=disable

# JWT
JWT_SECRET_KEY=your_secret_key_minimum_32_characters_long
JWT_ACCESS_TOKEN_DURATION=60
JWT_REFRESH_TOKEN_DURATION=30

# Server
PORT=8080
ENV=development
```

---

## Security Notes

1. **Password Storage**: Passwords are hashed using bcrypt with cost factor 10
2. **JWT Signing**: Tokens are signed with HS256 (HMAC-SHA256)
3. **Email Uniqueness**: Email addresses must be unique per account
4. **Soft Delete Support**: Deleted accounts can be recreated with the same email
5. **Token Expiration**: Access tokens expire after configured duration

---

## Next Steps

Future enhancements:
- [ ] Token refresh endpoint
- [ ] Password reset flow
- [ ] Email verification
- [ ] OAuth2 providers (Google, Facebook)
- [ ] Two-factor authentication (2FA)
- [ ] Rate limiting
- [ ] API key authentication for external services
