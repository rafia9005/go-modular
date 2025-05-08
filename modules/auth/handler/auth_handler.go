package handler

import (
	"fmt"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/jwt"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/auth/domain/service"
	"go-modular/modules/users/domain/entity"
	"go-modular/modules/users/dto/request"
	"go-modular/modules/users/dto/response"
	"net/http"

	"github.com/labstack/echo"
)

// AuthHandler struct handles HTTP request for auth.
type AuthHandler struct {
	authService *service.AuthService
	log         *logger.Logger
	event       *bus.EventBus
	jwt         jwt.JWT
	r           *utils.Response
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(log *logger.Logger, event *bus.EventBus, authService *service.AuthService, jwt jwt.JWT) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
		event:       event,
		jwt:         jwt,
		r:           &utils.Response{},
	}
}

// Initialize Event Handle.
func (h *AuthHandler) Handle(event bus.Event) {
	fmt.Printf("User created: %v", event.Payload)
}

// Register handles user registration.
func (h *AuthHandler) Register(c echo.Context) error {
	h.log.Info("Handling register request")

	req := new(request.CreateUserRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user := entity.NewUser(req.Name, req.Email, req.Password)
	err := h.authService.CreateUser(c.Request().Context(), user)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			h.log.Warn("Email already in use:", req.Email)
			return h.r.ErrorResponse(c, http.StatusConflict, "Email already in use")
		}
		h.log.Error("Failed to create user:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User created successfully:", user)

	h.event.Publish(bus.Event{Type: "user.created", Payload: user})
	h.log.Debug("Event 'user.created' published successfully")

	return h.r.SuccessResponse(c, map[string]interface{}{
		"user": response.FromEntity(user),
	}, "User registered successfully")
}

// Login handles user login.
func (h *AuthHandler) Login(c echo.Context) error {
	h.log.Info("Handling login request")

	req := new(request.LoginRequest)
	if err := c.Bind(req); err != nil {
		h.log.Error("Failed to bind request:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		h.log.Error("Validation failed:", err)
		return h.r.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	h.log.Debug("Request validated successfully:", req)

	user, err := h.authService.ProcessLogin(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrInvalidPassword {
			h.log.Warn("Invalid email or password for:", req.Email)
			return h.r.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		}
		h.log.Error("Failed to process login:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	h.log.Debug("User authenticated successfully:", user)

	tokenData := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	}

	token, err := h.jwt.GenerateToken(tokenData)
	if err != nil {
		h.log.Error("Failed to generate token:", err)
		return h.r.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"token": token,
		"user":  response.FromEntity(user),
	}, "Login successful")
}

// RegisterRoutes sets up the auth routes.
func (h *AuthHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	group := e.Group(basePath + "/auth")
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
}
