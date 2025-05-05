package auth

import (
	"go-modular-boilerplate/internal/pkg/bus"
	"go-modular-boilerplate/internal/pkg/logger"

	"gorm.io/gorm"
)

type Module struct {
	db     *gorm.DB
	logger *logger.Logger
}

func (m *Module) Initialize(db *gorm.DB, log *logger.Logger, event *bus.EventBus) error {

}
