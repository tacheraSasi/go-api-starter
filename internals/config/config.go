package config

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	DatabaseURL  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	ServerPort   string
	JWTSecret    string
	JWTExpiresIn string
	DBPath       string // For SQLite
	LogFilePath  string
	GINMode      string
	CORSOrigins  []string
	Timezone     string

	// Redis / asynq
	RedisURL string

	// Asynqmon monitoring dashboard
	MonitoringUser     string
	MonitoringPassword string
}

func LoadConfig() *Config {
	// Try .env from the current directory and a couple of parents so the
	// binaries work whether they are started from the repo root or cmd/<svc>.
	for _, path := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			break
		}
	}
	_ = godotenv.Load()

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
		DatabaseURL:  getEnvAny("", "DATABASE_URL", "DB_URL"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "user"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "dbname"),
		DBSSLMode:    getEnv("DB_SSL_MODE", "disable"),
		ServerPort:   serverPort,
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		JWTExpiresIn: jwtExpiresIn,
		DBPath:       getEnv("DB_PATH", "core.db"), // For SQLite
		LogFilePath:  logFilePath,
		GINMode:      ginMode,
		CORSOrigins:  origins,
		Timezone:     getEnv("TIMEZONE", "UTC"),

		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		MonitoringUser:     getEnv("MONITORING_USER", "admin"),
		MonitoringPassword: getEnv("MONITORING_PASSWORD", "changeme"),
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET must not be empty")
	}
	if strings.TrimSpace(c.ServerPort) == "" {
		return fmt.Errorf("SERVER_PORT must not be empty")
	}
	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return fmt.Errorf("TIMEZONE %q is invalid: %w", c.Timezone, err)
	}
	// When using postgres or mysql, either DATABASE_URL or the individual
	// DB fields must be set.
	if c.DBType == "postgres" || c.DBType == "mysql" {
		if strings.TrimSpace(c.DatabaseURL) == "" &&
			(strings.TrimSpace(c.DBHost) == "" || strings.TrimSpace(c.DBName) == "") {
			return fmt.Errorf("DATABASE_URL or DB_HOST+DB_NAME must be set for %s", c.DBType)
		}
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
