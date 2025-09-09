package config

import (
	"os"

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
}

func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		DBType:       getEnv("DB_TYPE", "sqlite"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "user"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "dbname"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "24"),//In hours
		DBPath:       getEnv("DB_PATH", "invoice_creator.db"), // For SQLite
		LogFilePath: "logs/app.log", // Log file path
	}
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
