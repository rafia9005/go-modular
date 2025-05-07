package middleware

import (
	"fmt"
	"go-modular/internal/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var jwtService jwt.JWT

func InitializeAuth(service jwt.JWT) {
	jwtService = service
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Authorization header is missing",
				"message": "Unauthorized",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid Authorization header format",
				"message": "Unauthorized",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   fmt.Sprintf("Invalid token: %v", err),
				"message": "Unauthorized",
			})
		}

		c.Set("user", claims)

		return next(c)
	}
}
