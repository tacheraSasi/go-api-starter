package logger

import (
	"os"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// NewServiceLogger creates a console logger for a standalone service
// (core, minion, overseer). It mirrors the cloud service logging style:
// human-readable text, full timestamps, and a service:<instance> field.
func NewServiceLogger(service string) *logrus.Entry {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		ForceColors:     true,
	})
	log.SetLevel(logrus.DebugLevel)

	instance, _ := os.Hostname()
	if len(instance) > 8 {
		instance = instance[:8]
	}

	return log.WithField("service", service+":"+instance)
}

// AsynqLogger adapts a logrus entry to asynq's Logger interface.
type AsynqLogger struct {
	l *logrus.Entry
}

var _ asynq.Logger = (*AsynqLogger)(nil)

// NewAsynqLogger wraps a logrus entry for use by asynq servers and schedulers.
func NewAsynqLogger(l *logrus.Entry) *AsynqLogger {
	return &AsynqLogger{l: l.WithField("component", "asynq")}
}

func (a *AsynqLogger) Debug(args ...any) { a.l.Debug(args...) }
func (a *AsynqLogger) Info(args ...any)  { a.l.Info(args...) }
func (a *AsynqLogger) Warn(args ...any)  { a.l.Warn(args...) }
func (a *AsynqLogger) Error(args ...any) { a.l.Error(args...) }
func (a *AsynqLogger) Fatal(args ...any) { a.l.Fatal(args...) }
