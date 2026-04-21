package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/tasks"
	"github.com/tacheraSasi/go-api-starter/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	// ── Database (needed for subscription tasks) ─────────────────────
	if err := database.Connect(database.DBConfig{
		Type:     cfg.DBType,
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
		FilePath: cfg.DBPath,
	}); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := database.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.BlacklistedToken{},
	); err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	// ── Redis connection ─────────────────────────────────────────────
	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL:", err)
	}

	// ── asynq.Server (worker) ────────────────────────────────────────
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6, // 60 %
			"default":  3, // 30 %
			"low":      1, // 10 %
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeNotificationSMS, tasks.HandleSMSTask)
	mux.HandleFunc(tasks.TypeNotificationEmail, tasks.HandleEmailTask)

	log.Println("Minion worker starting (concurrency=10)")

	// Run worker (blocks until SIGINT/SIGTERM)
	if err := srv.Run(mux); err != nil {
		log.Fatal("Worker failed:", err)
	}
}
