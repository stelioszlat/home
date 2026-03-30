package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	Chat   ChatConfig
	Socket SocketConfig
	Token  TokenConfig
	DB     DatabaseConfig
}

type ServerConfig struct {
	Port           string
	Host           string
	AllowedOrigins []string
	LogLever       string
	GinMode        string
}

type ChatConfig struct {
	ApiKey string
}

type SocketConfig struct {
	WSPingInterval          time.Duration
	StatusBroadcastInterval time.Duration
}

type TokenConfig struct {
	AccessToken string
}

type DatabaseConfig struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSslMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8060")
	host := getEnv("HOST", "127.0.0.1")
	allowedOriginsStr := getEnv("ALLOWED_ORIGINS", "*")
	logLevel := getEnv("LOG_LEVEL", "debug")
	ginMode := getEnv("GIN_MODE", "debug")

	chatApiKey := getEnv("CHAT_API_KEY", "")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "home")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbSslMode := getEnv("DB_SSLMODE", "disable")

	accessToken := getEnv("ACCESS_TOKEN", "")
	wsPingInterval := getEnvDuration("WS_PING_INTERVAL", 30*time.Second)
	statusBroadcastInterval := getEnvDuration("STATUS_BROADCAST_INTERVAL", 5*time.Second)

	// Parse allowed origins
	allowedOrigins := []string{allowedOriginsStr}
	if allowedOriginsStr == "*" {
		allowedOrigins = []string{"*"}
	}

	config := Config{
		Server: ServerConfig{
			Port:           port,
			Host:           host,
			AllowedOrigins: allowedOrigins,
			LogLever:       logLevel,
			GinMode:        ginMode,
		},
		Chat: ChatConfig{
			ApiKey: chatApiKey,
		},
		Socket: SocketConfig{
			WSPingInterval:          wsPingInterval,
			StatusBroadcastInterval: statusBroadcastInterval,
		},
		Token: TokenConfig{
			AccessToken: accessToken,
		},
		DB: DatabaseConfig{
			DBHost:     dbHost,
			DBPort:     dbPort,
			DBName:     dbName,
			DBUser:     dbUser,
			DBPassword: dbPassword,
			DBSslMode:  dbSslMode,
		},
	}

	return &config, nil
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to get duration from environment variable
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
