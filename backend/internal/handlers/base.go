package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string, err error) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
}

// SendSuccess sends a success response
func SendSuccess(c echo.Context, message string, data interface{}, statusCode int) error {
	return c.JSON(statusCode, SuccessResponse(message, data))
}

// SendError sends an error response
func SendError(c echo.Context, message string, err error, statusCode int) error {
	return c.JSON(statusCode, ErrorResponse(message, err))
}

// BadRequest sends a bad request response
func BadRequest(c echo.Context, message string) error {
	return SendError(c, message, nil, http.StatusBadRequest)
}

// InternalError sends an internal server error response
func InternalError(c echo.Context, message string, err error) error {
	return SendError(c, message, err, http.StatusInternalServerError)
}

// NotFound sends a not found response
func NotFound(c echo.Context, message string) error {
	return SendError(c, message, nil, http.StatusNotFound)
}
