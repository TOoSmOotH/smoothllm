package config

import (
	"os"
	"time"
)

type Config struct {
	GinMode string
	Port    string

	DBPath          string
	DBEncryptionKey string

	JWTSecret             string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration

	AllowedOrigins []string
}

func LoadConfig() *Config {
	return &Config{
		GinMode: getEnv("GIN_MODE", "debug"),
		Port:    getEnv("PORT", "8080"),

		DBPath:          getEnv("DB_PATH", "./data/app.db"),
		DBEncryptionKey: getEnv("DB_ENCRYPTION_KEY", "your-secret-key-32-char-long-for-aes256"),

		JWTSecret:             getEnv("JWT_SECRET", "your-jwt-secret-key-minimum-32-characters"),
		JWTAccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m"),
		JWTRefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h"),

		AllowedOrigins: []string{
			getEnv("FRONTEND_URL", "http://localhost:5173"),
			"http://localhost:3000",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
