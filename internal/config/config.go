package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func Load() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	expiry := 8 * time.Hour
	if v := os.Getenv("JWT_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			expiry = d
		}
	}

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        port,
		JWTSecret:   secret,
		JWTExpiry:   expiry,
	}
}
