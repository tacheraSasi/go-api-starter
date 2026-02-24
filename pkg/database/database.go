package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig holds database connection parameters
type DBConfig struct {
	Type     string // Database type: "mysql", "postgres", "sqlite"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // For PostgreSQL
	FilePath string // For SQLite
}

var DB *gorm.DB

// Connect establishes a connection to the specified database
func Connect(config DBConfig) error {
	var dialector gorm.Dialector
	var dsn string

	switch config.Type {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dsn = config.FilePath
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to %s database: %v", config.Type, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Printf("%s database connected successfully", config.Type)
	return nil
}

// GetDB returns the initialized database instance
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate runs auto migration for the given models
func AutoMigrate(models ...any) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.AutoMigrate(models...)
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}
	return sqlDB.Close()
}