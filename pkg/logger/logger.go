package logger

import (
    "os"
    "path/filepath"

    "github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger
type Logger struct {
    *logrus.Logger
}

// NewLogger initializes a new Logrus logger with JSON formatting and file output
func NewLogger(logFilePath string) (*Logger, error) {
    // Create a new Logrus logger
    log := logrus.New()

    // Set JSON formatter for structured logging
    log.SetFormatter(&logrus.JSONFormatter{
        PrettyPrint: true,
    })

    log.SetLevel(logrus.InfoLevel)

    if logFilePath != "" {
        if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
            return nil, err
        }

        file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            return nil, err
        }

        log.SetOutput(file)
    } else {
        log.SetOutput(os.Stdout)
    }

    return &Logger{log}, nil
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
    return l.Logger.WithFields(fields)
}