package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

// Init sets up the logger (dev or prod)
func Init(env string) error {
	var err error
	if env == "dev" {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	return err
}

func Sync() {
	log.Sync()
}

// Simple wrapper: converts key-value args into zap fields
func toFields(args ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue // ignore bad keys
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	return fields
}

// === Wrapper Methods ===

// Info logs an info-level message
func Info(msg string, args ...any) {
	log.Info(msg, toFields(args...)...)
}

// Warn logs a warning-level message
func Warn(msg string, args ...any) {
	log.Warn(msg, toFields(args...)...)
}

// Error logs an error-level message
func Error(msg string, args ...any) {
	log.Error(msg, toFields(args...)...)
}

// Debug logs a debug-level message
func Debug(msg string, args ...any) {
	log.Debug(msg, toFields(args...)...)
}

// Fatal logs a fatal-level message and exits
func Fatal(msg string, args ...any) {
	log.Fatal(msg, toFields(args...)...)
}
