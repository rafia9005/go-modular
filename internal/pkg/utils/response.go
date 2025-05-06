package utils

import (
	"net/http"

	"github.com/labstack/echo"
)

// Response is a standard structure for JSON responses
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONResponse sends a standard JSON response
func JSONResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	response := Response{
		Data:    data,
		Message: message,
		Error:   err,
	}
	return c.JSON(statusCode, response)
}

// SuccessResponse is a helper for success responses
func SuccessResponse(c echo.Context, data interface{}, message string) error {
	return JSONResponse(c, http.StatusOK, data, message, "")
}

// ErrorResponse is a helper for error responses
func ErrorResponse(c echo.Context, statusCode int, err string) error {
	return JSONResponse(c, statusCode, nil, "", err)
}

// CustomResponse is a helper for custom responses with flexible status codes
func CustomResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	return JSONResponse(c, statusCode, data, message, err)
}
