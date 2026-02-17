package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigKey string

const (
	DBHostKey       ConfigKey = "DB_HOST"
	DBPortKey       ConfigKey = "DB_PORT"
	DBUserKey       ConfigKey = "DB_USER"
	DBPasswordKey   ConfigKey = "DB_PASSWORD"
	DBNameKey       ConfigKey = "DB_NAME"
	ServerPortKey   ConfigKey = "SERVER_PORT"
	JWTSecretKey    ConfigKey = "JWT_SECRET"
	JWTExpiresInKey ConfigKey = "JWT_EXPIRES_IN"
)

type Config struct {
	DBType       string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	ServerPort   string
	JWTSecret    string
	JWTExpiresIn string
	DBPath       string // For SQLite
	LogFilePath  string
	GINMode      string
	CORSOrigins  []string
}

func LoadConfig() *Config {
	godotenv.Load()

	serverPort := getEnvAny("8080", "SERVER_PORT", "PORT")
	jwtExpiresIn := getEnvAny("24", "JWT_EXPIRES_IN", "JWT_EXPIRY_HOURS")
	logFilePath := getEnvAny("logs/app.log", "LOG_FILE_PATH")
	ginMode := getEnvAny("release", "GIN_MODE")
	corsAllowedOrigins := getEnvAny("*", "CORS_ALLOWED_ORIGINS")

	origins := splitAndTrim(corsAllowedOrigins)
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	return &Config{
		DBType:       getEnv("DB_TYPE", "sqlite"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "user"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "dbname"),
		ServerPort:   serverPort,
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		JWTExpiresIn: jwtExpiresIn,
		DBPath:       getEnv("DB_PATH", "core.db"), // For SQLite
		LogFilePath:  logFilePath,
		GINMode:      ginMode,
		CORSOrigins:  origins,
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET must not be empty")
	}
	if strings.TrimSpace(c.ServerPort) == "" {
		return fmt.Errorf("SERVER_PORT must not be empty")
	}
	return nil
}

func (c *Config) Get(key ConfigKey) string {
	values := map[ConfigKey]string{
		DBHostKey:       c.DBHost,
		DBPortKey:       c.DBPort,
		DBUserKey:       c.DBUser,
		DBPasswordKey:   c.DBPassword,
		DBNameKey:       c.DBName,
		ServerPortKey:   c.ServerPort,
		JWTExpiresInKey: c.JWTExpiresIn,
		JWTSecretKey:    c.JWTSecret,
	}
	return values[key]
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(string(key))
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAny(defaultValue string, keys ...string) string {
	for _, key := range keys {
		if value, exists := os.LookupEnv(key); exists {
			trimmed := strings.TrimSpace(value)
			if trimmed != "" {
				return trimmed
			}
		}
	}
	return defaultValue
}

func splitAndTrim(input string) []string {
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
