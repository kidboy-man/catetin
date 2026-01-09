package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ingunawandra/catetin/internal/controller/dto"
	"github.com/ingunawandra/catetin/internal/controller/http/middleware"
	"github.com/ingunawandra/catetin/internal/service"
	appErrors "github.com/ingunawandra/catetin/pkg/errors"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// POST /api/v1/authentications/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, appErrors.ErrValidation.WithDetails(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	// Call service
	result, err := h.authService.Register(c.Request.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		middleware.AbortWithError(c, err)
		return
	}

	// Build response
	response := &dto.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    result.ExpiresIn,
		User: &dto.UserInfo{
			ID:          result.User.ID.String(),
			FullName:    result.User.FullName,
			Email:       req.Email,
			PhoneNumber: &result.User.PhoneNumber,
			Image:       result.User.Image,
		},
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse("User registered successfully", response))
}

// Login handles user login
// POST /api/v1/authentications/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, appErrors.ErrValidation.WithDetails(map[string]interface{}{
			"validation_errors": err.Error(),
		}))
		return
	}

	// Call service
	result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		middleware.AbortWithError(c, err)
		return
	}

	// Build response
	response := &dto.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    result.ExpiresIn,
		User: &dto.UserInfo{
			ID:          result.User.ID.String(),
			FullName:    result.User.FullName,
			Email:       req.Email,
			PhoneNumber: &result.User.PhoneNumber,
			Image:       result.User.Image,
		},
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse("Login successful", response))
}
