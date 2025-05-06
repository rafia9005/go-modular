package handler

import (
	"fmt"
	"go-modular-boilerplate/internal/pkg/bus"
	"go-modular-boilerplate/internal/pkg/jwt"
	"go-modular-boilerplate/internal/pkg/logger"
	"go-modular-boilerplate/modules/auth/domain/service"
	"go-modular-boilerplate/modules/users/domain/entity"
	"go-modular-boilerplate/modules/users/dto/request"
	"go-modular-boilerplate/modules/users/dto/response"
	"net/http"

	"github.com/labstack/echo"
)

// AuthHandler struct handles HTTP request for auth
type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
	event       *bus.EventBus
	jwt         jwt.JWT
}

// creates a new auth handler
func NewAuthHandler(log *logger.Logger, event *bus.EventBus, authService *service.AuthService, jwt jwt.JWT) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
		event:       event,
		jwt:         jwt,
	}
}

// Initialize Event Handle
func (h *AuthHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	h.log.Info("Handling register request")

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.log.Debug("Request validated successfully:", req)

	user := entity.NewUser(req.Name, req.Email, req.Password)
	err := h.authService.CreateUser(c.Request().Context(), user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already in use"})
		}
		h.log.Error("Failed to create user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.log.Debug("User created successfully:", user)

	h.event.Publish(bus.Event{Type: "user.created", Payload: user})
	h.log.Debug("Event 'user.created' published successfully")

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"user": response.FromEntity(user),
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	h.log.Info("Handling login request")

	req := new(request.LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.log.Debug("Request validated successfully:", req)

	user, err := h.authService.ProcessLogin(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound {
			h.log.Warn("User not found:", req.Email)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}
		if err == service.ErrInvalidPassword {
			h.log.Warn("Invalid password for email:", req.Email)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}
		h.log.Error("Failed to process login:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.log.Debug("User authenticated successfully:", user)

	// Return user information or JWT token (if implemented)
	tokenData := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	}

	token, err := h.jwt.GenerateToken(tokenData)

	if err != nil {
		h.log.Error("Failed to generate token:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"token":   token,
			"user":    response.FromEntity(user),
			"message": "Login successful",
		},
	})
}

// RegisterRoutes sets up the auth routes
func (h *AuthHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath + "/auth")
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
}
