package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	//* Database
	DatabaseURL string

	//* Server
	Port     int
	GRPCPort string

	//* Features
	LogLevel string
	APIKey   string

	//* JWT
	JWTSecret          string
	JWTRefreshSecret   string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func LoadConfigDev() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:erke@localhost:5432/auth_service?sslmode=disable"),
		Port:               getEnvAsInt("PORT", 8080),
		GRPCPort:           getEnv("GRPC_PORT", "50051"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		APIKey:             getEnv("API_KEY", ""),
		JWTSecret:          getEnv("JWT_ACCESS_SECRET", "your-access-secret-key-here"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key-here"),
		AccessTokenExpiry:  getEnvAsDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", 168*time.Hour), // 7 days
	}
}

// Добавьте метод для duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Остальные методы остаются такими же...
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
