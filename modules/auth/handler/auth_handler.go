package handler

import (
	"go-modular-boilerplate/internal/pkg/bus"
	"go-modular-boilerplate/internal/pkg/logger"
)

type AuthHadnler struct {
	log   *logger.Logger
	event *bus.EventBus
}


