package utils

import (
	"net/http"

	"github.com/labstack/echo"
)

// Response is a standard structure for JSON responses.
type Response struct{}

// JSONResponse sends a standard JSON response.
func (r *Response) JSONResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	response := map[string]interface{}{
		"data":    data,
		"message": message,
		"error":   err,
	}
	return c.JSON(statusCode, response)
}

// SuccessResponse is a helper for success responses.
func (r *Response) SuccessResponse(c echo.Context, data interface{}, message string) error {
	return r.JSONResponse(c, http.StatusOK, data, message, "")
}

// ErrorResponse is a helper for error responses.
func (r *Response) ErrorResponse(c echo.Context, statusCode int, err string) error {
	return r.JSONResponse(c, statusCode, nil, "", err)
}

// CreatedResponse is a helper for responses with HTTP 201 Created.
func (r *Response) CreatedResponse(c echo.Context, data interface{}, message string) error {
	return r.JSONResponse(c, http.StatusCreated, data, message, "")
}

// NoContentResponse is a helper for HTTP 204 No Content responses.
func (r *Response) NoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// ForbiddenResponse is a helper for HTTP 403 Forbidden responses.
func (r *Response) ForbiddenResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusForbidden, nil, "", err)
}

// UnauthorizedResponse is a helper for HTTP 401 Unauthorized responses.
func (r *Response) UnauthorizedResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusUnauthorized, nil, "", err)
}

// BadRequestResponse is a helper for HTTP 400 Bad Request responses.
func (r *Response) BadRequestResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusBadRequest, nil, "", err)
}

// InternalServerErrorResponse is a helper for HTTP 500 Internal Server Error responses.
func (r *Response) InternalServerErrorResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusInternalServerError, nil, "", err)
}

// ConflictResponse is a helper for HTTP 409 Conflict responses.
func (r *Response) ConflictResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusConflict, nil, "", err)
}

// NotFoundResponse is a helper for HTTP 404 Not Found responses.
func (r *Response) NotFoundResponse(c echo.Context, err string) error {
	return r.JSONResponse(c, http.StatusNotFound, nil, "", err)
}

// CustomResponse is a helper for custom responses with flexible status codes.
func (r *Response) CustomResponse(c echo.Context, statusCode int, data interface{}, message string, err string) error {
	return r.JSONResponse(c, statusCode, data, message, err)
}
