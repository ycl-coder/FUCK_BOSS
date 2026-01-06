// Package logger provides structured logging using zap.
// It supports different log levels, formats (JSON/Console), and context-aware logging.
package logger

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface for structured logging.
type Logger interface {
	// Debug logs a debug message with optional fields.
	Debug(msg string, fields ...zap.Field)

	// Info logs an info message with optional fields.
	Info(msg string, fields ...zap.Field)

	// Warn logs a warning message with optional fields.
	Warn(msg string, fields ...zap.Field)

	// Error logs an error message with optional fields.
	Error(msg string, fields ...zap.Field)

	// WithContext returns a logger with context fields (e.g., request ID).
	WithContext(ctx context.Context) Logger

	// WithFields returns a logger with additional fields.
	WithFields(fields ...zap.Field) Logger

	// Sync flushes any buffered log entries.
	Sync() error
}

// zapLogger is the zap-based implementation of Logger.
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new logger based on the environment and log configuration.
// environment can be "development" or "production".
// If logConfig is nil, default configuration is used.
func NewLogger(environment string, logConfig *LogConfig) (Logger, error) {
	var config zap.Config

	if environment == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	// Apply log configuration if provided
	if logConfig != nil {
		// Set log level
		level, err := parseLogLevel(logConfig.Level)
		if err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		config.Level = zap.NewAtomicLevelAt(level)

		// Set log format
		if logConfig.Format == "text" || logConfig.Format == "console" {
			config.Encoding = "console"
			config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		} else {
			config.Encoding = "json"
		}

		// Set output paths
		if len(logConfig.OutputPaths) > 0 {
			config.OutputPaths = logConfig.OutputPaths
		}
		if len(logConfig.ErrorOutputPaths) > 0 {
			config.ErrorOutputPaths = logConfig.ErrorOutputPaths
		}
	}

	// Ensure caller information is enabled (default is false, but explicit is better)
	config.DisableCaller = false
	// Stacktrace is only shown for errors by default, which is fine
	config.DisableStacktrace = false

	// Build logger
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Add caller skip to show the actual caller location instead of logger.go
	// Skip 1 level: logger.go wrapper methods (Debug, Info, Warn, Error)
	// This ensures the caller shows the actual code location (e.g., main.go:53)
	// instead of the wrapper location (logger.go:139)
	logger = logger.WithOptions(zap.AddCallerSkip(1))

	return &zapLogger{logger: logger}, nil
}

// NewLoggerFromConfig creates a new logger from LogConfig.
// It automatically determines the environment based on the log format.
func NewLoggerFromConfig(logConfig *LogConfig) (Logger, error) {
	environment := "production"
	if logConfig != nil && (logConfig.Format == "text" || logConfig.Format == "console") {
		environment = "development"
	}
	return NewLogger(environment, logConfig)
}

// LogConfig represents logging configuration.
type LogConfig struct {
	// Level is the log level (debug, info, warn, error).
	Level string

	// Format is the log format (json, text, console).
	Format string

	// OutputPaths is a list of paths to write logging output to.
	OutputPaths []string

	// ErrorOutputPaths is a list of paths to write error level logs to.
	ErrorOutputPaths []string
}

// parseLogLevel parses a log level string and returns the corresponding zapcore.Level.
func parseLogLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// Debug logs a debug message with optional fields.
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message with optional fields.
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message with optional fields.
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message with optional fields.
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// WithContext returns a logger with context fields.
// It extracts request ID and other context values if available.
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	fields := extractContextFields(ctx)
	return &zapLogger{logger: l.logger.With(fields...)}
}

// WithFields returns a logger with additional fields.
func (l *zapLogger) WithFields(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// Sync flushes any buffered log entries.
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// extractContextFields extracts fields from context.
// Currently supports request ID from context value with key "request_id".
func extractContextFields(ctx context.Context) []zap.Field {
	var fields []zap.Field

	// Extract request ID if available
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	// Extract trace ID if available
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// Extract user ID if available
	if userID, ok := ctx.Value("user_id").(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	return fields
}

// RequestIDKey is the context key for request ID.
const RequestIDKey = "request_id"

// TraceIDKey is the context key for trace ID.
const TraceIDKey = "trace_id"

// UserIDKey is the context key for user ID.
const UserIDKey = "user_id"

// WithRequestID adds request ID to context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithTraceID adds trace ID to context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithUserID adds user ID to context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
