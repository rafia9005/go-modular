package auth

import (
	"go-modular-boilerplate/internal/pkg/bus"
	"go-modular-boilerplate/internal/pkg/logger"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Module struct {
	db     *gorm.DB
	logger *logger.Logger
}

func (m *Module) Name() string {
	return "auth"
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {
	m.db = db
	m.logger = log
	m.event = event

	m.logger.Info("Initializing auth module")

	m.logger.Debug("")
}

func (m *Module) RegisterRoutes(e *echo.Echo, basePath string) {
	m.logger.Info("Registering user routes at %s/auth", basePath)
	
}