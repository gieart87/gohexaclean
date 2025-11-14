package logger

import (
	"os"

	"github.com/gieart87/gohexaclean/internal/infra/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	var zapConfig zap.Config

	if cfg.Format == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set output
	if cfg.Output == "stdout" {
		zapConfig.OutputPaths = []string{"stdout"}
	}

	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// NewDefaultLogger creates a default logger for development
func NewDefaultLogger() *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build(zap.AddCallerSkip(1))
	return &Logger{Logger: logger}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.With(zap.Any(key, value))}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{Logger: l.With(zapFields...)}
}

// Close closes the logger
func (l *Logger) Close() error {
	return l.Sync()
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(cfg *config.LoggerConfig) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewDefaultLogger()
	}
	return globalLogger
}

// Helper functions for global logger
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
	os.Exit(1)
}
