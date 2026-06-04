package main

import (
	"ping-uptime/internal/app"
	"ping-uptime/internal/pkg/config"
	"ping-uptime/internal/pkg/logger"
	"ping-uptime/internal/pkg/middleware"
	"ping-uptime/modules/auth"
	user "ping-uptime/modules/users"
	"log"
	"os"
)

func main() {

	// Load configuration from .env
	if err := config.Initialize(); err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	// initialize logger
	logCfg := logger.DefaultConfig()

	// Start the application
	application, err := app.NewApp(&logCfg)
	if err != nil {
		log.Fatalf("Error creating application : %v", err)
		os.Exit(1)
	}

	// Initialize Auth middleware
	jwtSignatureKey := config.GetJWTService()
	middleware.InitializeAuth(jwtSignatureKey)

	// register modules
	application.RegisterModule(user.NewModule())
	application.RegisterModule(auth.NewModule())

	// initialize the application
	if err := application.Initialize(); err != nil {
		log.Fatalf("Error initializing application : %v", err)
		os.Exit(1)
	}

	// Start the application
	application.Start()
}
