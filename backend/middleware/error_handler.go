package middleware

import (
	"errors"
	"go-web/services"
	"go-web/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorResponse represents a standardized error response format.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewErrorResponseWithDetails creates a new ErrorResponse instance with details.
func NewErrorResponseWithDetails(code int, message, details string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var errorMappings = map[error]int{
	services.ErrPermissionDenied:        http.StatusForbidden,
	&services.UserExistsError{}:         http.StatusConflict,
	&services.InvalidCredentialsError{}: http.StatusUnauthorized,
	gorm.ErrRecordNotFound:              http.StatusNotFound,
}

// ErrorHandler is a middleware to handle errors in a centralized way.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during the request
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Log the error
			utils.Logger.Error("Request error", zap.Error(err))

			// Default to internal server error
			status := http.StatusInternalServerError
			msg := "Internal Server Error"

			// Check for specific error types
			for e, s := range errorMappings {
				if errors.Is(err, e) {
					status = s
					msg = err.Error()
					break
				}
			}

			// Handle validation errors
			if validationErrors, ok := err.(interface{ Errors() []error }); ok {
				var errorMessages []string
				for _, e := range validationErrors.Errors() {
					errorMessages = append(errorMessages, e.Error())
				}
				c.JSON(http.StatusBadRequest, NewErrorResponseWithDetails(http.StatusBadRequest, "Validation failed", strings.Join(errorMessages, "; ")))
				return
			}

			// Handle binding errors (e.g., JSON unmarshalling errors)
			if strings.Contains(err.Error(), "bind") || strings.Contains(err.Error(), "json") {
				status = http.StatusBadRequest
				msg = "Invalid request format"
			}

			c.JSON(status, NewErrorResponse(status, msg))
		}
	}
}