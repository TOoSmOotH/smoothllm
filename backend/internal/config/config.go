package config

import (
	"os"
	"strings"
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
	FrontendURL    string
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

		AllowedOrigins: getOriginsEnv("CORS_ORIGINS", "*"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:5173"),
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

func getOriginsEnv(key string, defaultValue string) []string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	// Handle wildcard
	if value == "*" {
		return []string{"*"}
	}

	// Split comma-separated origins and trim whitespace
	origins := strings.Split(value, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{"*"}
	}

	return result
}
