package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mayurbotre/task-management-system/internal/middleware"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Envelope struct {
	Status     string         `json:"status"`              // "success" | "error"
	StatusCode int            `json:"statusCode"`          // HTTP status code
	Data       any            `json:"data,omitempty"`      // payload on success
	Error      *ErrorBody     `json:"error,omitempty"`     // error object on failure
	Meta       map[string]any `json:"meta,omitempty"`      // pagination or other metadata
	RequestID  string         `json:"requestId,omitempty"` // from middleware
}

func respond(c *gin.Context, status int, data any, meta map[string]any) {
	rid := c.GetString(middleware.RequestIDKey)
	c.JSON(status, Envelope{
		Status:     "success",
		StatusCode: status,
		Data:       data,
		Meta:       meta,
		RequestID:  rid,
	})
}

func respondError(c *gin.Context, status int, code, message string, details any) {
	rid := c.GetString(middleware.RequestIDKey)
	c.AbortWithStatusJSON(status, Envelope{
		Status:     "error",
		StatusCode: status,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: rid,
	})
}

// Convenience wrappers
func ok(c *gin.Context, data any)      { respond(c, http.StatusOK, data, nil) }
func created(c *gin.Context, data any) { respond(c, http.StatusCreated, data, nil) }
func okWithMeta(c *gin.Context, data any, meta map[string]any) {
	respond(c, http.StatusOK, data, meta)
}
