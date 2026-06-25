package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv       string
	Port         string
	DatabaseURL  string
	JWTSecret    string
	JWTTTL       time.Duration
	CORSOrigin   string
	GeminiAPIKey string
}

func Load() Config {
	jwtTTLMinutes := getEnvAsInt("JWT_TTL_MINUTES", 60)

	return Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-development"),
		JWTTTL:       time.Duration(jwtTTLMinutes) * time.Minute,
		CORSOrigin:   getEnv("CORS_ALLOWED_ORIGIN", "*"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
