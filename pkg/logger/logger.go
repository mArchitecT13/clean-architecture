package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger interface defines the logging methods
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithContext(ctx context.Context) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// logger implements the Logger interface
type logger struct {
	logrus *logrus.Logger
}

// New creates a new logger instance
func New() Logger {
	l := logrus.New()

	// Set output to stdout
	l.SetOutput(os.Stdout)

	// Set log level from environment or default to info
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}

	// Set formatter to JSON for structured logging
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})

	return &logger{logrus: l}
}

// Debug logs debug level message
func (l *logger) Debug(args ...interface{}) {
	l.logrus.Debug(args...)
}

// Info logs info level message
func (l *logger) Info(args ...interface{}) {
	l.logrus.Info(args...)
}

// Warn logs warning level message
func (l *logger) Warn(args ...interface{}) {
	l.logrus.Warn(args...)
}

// Error logs error level message
func (l *logger) Error(args ...interface{}) {
	l.logrus.Error(args...)
}

// Fatal logs fatal level message and exits
func (l *logger) Fatal(args ...interface{}) {
	l.logrus.Fatal(args...)
}

// Debugf logs formatted debug level message
func (l *logger) Debugf(format string, args ...interface{}) {
	l.logrus.Debugf(format, args...)
}

// Infof logs formatted info level message
func (l *logger) Infof(format string, args ...interface{}) {
	l.logrus.Infof(format, args...)
}

// Warnf logs formatted warning level message
func (l *logger) Warnf(format string, args ...interface{}) {
	l.logrus.Warnf(format, args...)
}

// Errorf logs formatted error level message
func (l *logger) Errorf(format string, args ...interface{}) {
	l.logrus.Errorf(format, args...)
}

// Fatalf logs formatted fatal level message and exits
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logrus.Fatalf(format, args...)
}

// WithContext returns a logger with context
func (l *logger) WithContext(ctx context.Context) Logger {
	return &logger{logrus: l.logrus.WithContext(ctx).Logger}
}

// WithField returns a logger with a single field
func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{logrus: l.logrus.WithField(key, value).Logger}
}

// WithFields returns a logger with multiple fields
func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{logrus: l.logrus.WithFields(fields).Logger}
}
