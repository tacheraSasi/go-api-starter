package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/minions"
	"github.com/tacheraSasi/go-api-starter/pkg/database"
	"github.com/tacheraSasi/go-api-starter/pkg/logger"
)

// minion is the background worker. It consumes asynq tasks from Redis. The
// overseer enqueues recurring work, and the core service can enqueue ad-hoc
// work (emails, SMS, exports, etc.).
func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load timezone: %v\n", err)
		os.Exit(1)
	}
	time.Local = location

	l := logger.NewServiceLogger("minion")
	l.WithFields(logrus.Fields{
		"go":       runtime.Version(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"pid":      os.Getpid(),
		"timezone": cfg.Timezone,
	}).Info("Starting minion worker")

	// Database access is optional for notification-only workers but most real
	// workers need it for repository-backed tasks.
	if err := database.Connect(database.DBConfig{
		Type:        cfg.DBType,
		DatabaseURL: cfg.DatabaseURL,
		Host:        cfg.DBHost,
		Port:        cfg.DBPort,
		User:        cfg.DBUser,
		Password:    cfg.DBPassword,
		DBName:      cfg.DBName,
		SSLMode:     cfg.DBSSLMode,
		FilePath:    cfg.DBPath,
	}); err != nil {
		l.WithError(err).Fatal("Database connection failed")
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			l.WithError(closeErr).Error("Failed to close database")
		}
	}()
	l.WithField("type", cfg.DBType).Info("Database connected")

	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		l.WithError(err).Fatal("Failed to parse REDIS_URL")
	}
	l.WithField("redis", cfg.RedisURL).Info("Redis connection configured")

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		Logger:   logger.NewAsynqLogger(l),
		LogLevel: asynq.InfoLevel,
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, taskErr error) {
			retried, _ := asynq.GetRetryCount(ctx)
			maxRetry, _ := asynq.GetMaxRetry(ctx)
			taskID, _ := asynq.GetTaskID(ctx)
			queue, _ := asynq.GetQueueName(ctx)

			l.WithFields(logrus.Fields{
				"task":      task.Type(),
				"task_id":   taskID,
				"queue":     queue,
				"retry":     retried,
				"max_retry": maxRetry,
			}).WithError(taskErr).Error("Task failed")
		}),
	})

	// ── Task handlers ──────────────────────────────────────────────────────
	// Register one handler per task type. Keep the binaries thin: handlers
	// should delegate to internals/services or dedicated minions types.
	notificationHandler := minions.NewNotificationHandler(l)
	maintenanceHandler := minions.NewMaintenanceHandler(l)

	// Repositories are created here for handlers that need database access:
	//   tokenRepo := repositories.NewTokenRepository(database.GetDB())

	mux := asynq.NewServeMux()
	mux.Use(structuredLoggingMiddleware(l))
	mux.HandleFunc(minions.TaskSendSMS, notificationHandler.HandleSendSMS)
	mux.HandleFunc(minions.TaskSendEmail, notificationHandler.HandleSendEmail)
	mux.HandleFunc(minions.TaskHeartbeat, maintenanceHandler.HandleHeartbeat)

	l.WithFields(logrus.Fields{
		"concurrency": 10,
		"queues":      "critical(6), default(3), low(1)",
		"handlers":    3,
	}).Info("Task handlers registered")

	// ── Start ──────────────────────────────────────────────────────────────
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Run(mux); err != nil {
			l.WithError(err).Fatal("Worker server failed")
		}
	}()

	l.Info("Minion worker started, waiting for tasks")
	sig := <-sigCh
	l.WithField("signal", sig.String()).Info("Shutdown signal received")
	srv.Shutdown()
	l.Info("Minion worker stopped gracefully")
}

// structuredLoggingMiddleware wraps every task with structured logs and panic
// recovery so one bad task cannot take the worker process down.
func structuredLoggingMiddleware(l *logrus.Entry) func(asynq.Handler) asynq.Handler {
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) (err error) {
			taskID, _ := asynq.GetTaskID(ctx)
			queue, _ := asynq.GetQueueName(ctx)
			retried, _ := asynq.GetRetryCount(ctx)
			maxRetry, _ := asynq.GetMaxRetry(ctx)

			entry := l.WithFields(logrus.Fields{
				"task":    t.Type(),
				"task_id": taskID,
				"queue":   queue,
				"retry":   fmt.Sprintf("%d/%d", retried, maxRetry),
			})

			entry.Debug("Processing task")
			defer func() {
				if r := recover(); r != nil {
					entry.WithField("panic", r).Error("Task panicked")
					err = fmt.Errorf("panic in task %s: %v", t.Type(), r)
				}
			}()

			start := time.Now()
			err = h.ProcessTask(ctx, t)
			entry = entry.WithField("elapsed", time.Since(start).String())
			if err != nil {
				entry.WithError(err).Error("Task failed")
				return err
			}
			entry.Debug("Task completed")
			return nil
		})
	}
}
