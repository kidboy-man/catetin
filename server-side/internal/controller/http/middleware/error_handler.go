package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ingunawandra/catetin/internal/controller/dto"
	appErrors "github.com/ingunawandra/catetin/pkg/errors"
)

// ErrorHandler is a middleware that handles errors returned by handlers
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error (most recent)
		err := c.Errors.Last().Err

		// Check if it's an AppError
		if appErr, ok := appErrors.IsAppError(err); ok {
			// Use AppError details
			response := dto.ErrorResponse{
				Status:  "error",
				Message: appErr.Message,
				Errors: map[string]interface{}{
					"code": appErr.Code,
				},
			}

			// Add additional details if present
			if appErr.Details != nil {
				for k, v := range appErr.Details {
					response.Errors.(map[string]interface{})[k] = v
				}
			}

			c.JSON(appErr.HTTPStatus, response)
			return
		}

		// Handle non-AppError as internal server error
		log.Printf("Unhandled error: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "error",
			Message: "An internal error occurred",
			Errors: map[string]interface{}{
				"code": appErrors.ErrCodeInternal,
			},
		})
	}
}

// AbortWithError is a helper to abort with an AppError
func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

// AbortWithAppError is a helper to abort with a specific AppError
func AbortWithAppError(c *gin.Context, appErr *appErrors.AppError) {
	_ = c.Error(appErr)
	c.Abort()
}
