package main

import (
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
	"github.com/tacheraSasi/go-api-starter/pkg/logger"
)

// overseer is the cron scheduler. It does not process tasks itself; it
// enqueues them into Redis on a schedule, where one or more minions pick them
// up. Run exactly one overseer per environment to avoid duplicate schedules.
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

	l := logger.NewServiceLogger("overseer")
	l.WithFields(logrus.Fields{
		"go":       runtime.Version(),
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"pid":      os.Getpid(),
		"timezone": cfg.Timezone,
	}).Info("Starting overseer scheduler")

	redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		l.WithError(err).Fatal("Failed to parse REDIS_URL")
	}
	l.WithField("redis", cfg.RedisURL).Info("Redis connection configured")

	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
		Location: location,
		Logger:   logger.NewAsynqLogger(l),
		LogLevel: asynq.InfoLevel,
		PreEnqueueFunc: func(task *asynq.Task, _ []asynq.Option) {
			l.WithField("task", task.Type()).Debug("Preparing scheduled task")
		},
		PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
			entry := l.WithFields(logrus.Fields{
				"task":    info.Type,
				"task_id": info.ID,
				"queue":   info.Queue,
			})
			if err != nil {
				entry.WithError(err).Error("Failed to enqueue scheduled task")
				return
			}
			entry.Info("Enqueued scheduled task")
		},
	})

	// ── Schedules ──────────────────────────────────────────────────────────
	// Cron format: minute hour day-of-month month day-of-week.
	// Add business schedules here; the matching handler must be registered in
	// cmd/minion. The heartbeat is only a template placeholder.
	register := func(spec string, task *asynq.Task, opts ...asynq.Option) {
		id, registerErr := scheduler.Register(spec, task, opts...)
		if registerErr != nil {
			l.WithError(registerErr).WithField("cron", spec).Fatal("Failed to register schedule")
		}
		l.WithFields(logrus.Fields{
			"task":     task.Type(),
			"cron":     spec,
			"entry_id": id,
		}).Info("Registered scheduled task")
	}

	register("0 * * * *", asynq.NewTask(minions.TaskHeartbeat, nil), asynq.Queue("low"))
	// Examples (add the matching task constant in internals/minions first):
	//   register("0 2 * * *", asynq.NewTask("cleanup:expired_tokens", nil), asynq.Queue("default"))
	//   register("0 8 * * *", asynq.NewTask("notification:daily_digest", nil), asynq.Queue("default"))

	// ── Start ──────────────────────────────────────────────────────────────
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := scheduler.Run(); err != nil {
			l.WithError(err).Fatal("Scheduler failed")
		}
	}()

	l.Info("Overseer scheduler running, waking minions on schedule")
	sig := <-sigCh
	l.WithField("signal", sig.String()).Info("Shutdown signal received")
	scheduler.Shutdown()
	l.Info("Overseer scheduler stopped gracefully")
}
