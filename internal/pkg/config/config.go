package config

import (
	"ping-uptime/internal/pkg/jwt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Initialize loads the .env file and environment variables.
func Initialize() error {
	if err := godotenv.Load(); err != nil {
		// .env file is optional — real env vars still work
		log.Println("No .env file found, reading from environment variables")
	}
	return nil
}

// checkKey aborts if the environment variable is not set.
func checkKey(key string) {
	if os.Getenv(key) == "" {
		log.Fatalf("Configuration key %s not found; aborting\n", key)
		os.Exit(1)
	}
}

// GetString returns the value of an environment variable as a string.
func GetString(key string) string {
	checkKey(key)
	return os.Getenv(key)
}

// GetInt returns the value of an environment variable as an int.
func GetInt(key string) int {
	checkKey(key)
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatalf("Configuration key %s is not a valid integer: %v\n", key, err)
		os.Exit(1)
	}
	return val
}

// GetBool returns the value of an environment variable as a bool.
func GetBool(key string) bool {
	checkKey(key)
	val, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.Fatalf("Configuration key %s is not a valid boolean: %v\n", key, err)
		os.Exit(1)
	}
	return val
}

// GetJWTService returns a JWT service using the signature key from the environment.
func GetJWTService() jwt.JWT {
	signatureKey := GetString("JWT_SIGNATURE_KEY")
	if signatureKey == "" {
		panic("JWT signature key not found in configuration")
	}
	return jwt.NewJWTImpl(signatureKey, 7)
}
