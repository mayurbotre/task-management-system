package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortJSON(c *gin.Context, status int, msg string) { // keep name for minimal diff
	code := "ERROR"
	switch status {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusNotFound:
		code = "NOT_FOUND"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
	case http.StatusForbidden:
		code = "FORBIDDEN"
	case http.StatusConflict:
		code = "CONFLICT"
	case http.StatusUnprocessableEntity:
		code = "UNPROCESSABLE_ENTITY"
	case http.StatusInternalServerError:
		code = "INTERNAL"
	}
	respondError(c, status, code, msg, nil)
}
