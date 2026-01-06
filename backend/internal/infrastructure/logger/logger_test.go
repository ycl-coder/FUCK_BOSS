package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger_Development(t *testing.T) {
	logger, err := NewLogger("development", nil)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	// Test all log levels
	logger.Debug("debug message", zap.String("key", "value"))
	logger.Info("info message", zap.String("key", "value"))
	logger.Warn("warn message", zap.String("key", "value"))
	logger.Error("error message", zap.String("key", "value"))

	// Sync to flush logs
	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error (may be expected): %v", err)
	}
}

func TestNewLogger_Production(t *testing.T) {
	logger, err := NewLogger("production", nil)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	logger.Info("production log", zap.String("key", "value"))
	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error (may be expected): %v", err)
	}
}

func TestNewLogger_WithConfig(t *testing.T) {
	config := &LogConfig{
		Level:            "debug",
		Format:           "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := NewLogger("production", config)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Debug("debug with config", zap.String("test", "value"))
	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error (may be expected): %v", err)
	}
}

func TestNewLoggerFromConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *LogConfig
		wantErr  bool
	}{
		{
			name: "json format",
			config: &LogConfig{
				Level:  "info",
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "text format",
			config: &LogConfig{
				Level:  "info",
				Format: "text",
			},
			wantErr: false,
		},
		{
			name: "console format",
			config: &LogConfig{
				Level:  "info",
				Format: "console",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: &LogConfig{
				Level:  "invalid",
				Format: "json",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLoggerFromConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLoggerFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("NewLoggerFromConfig() returned nil logger")
			}
			if logger != nil {
				logger.Sync()
			}
		})
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		want    zapcore.Level
		wantErr bool
	}{
		{
			name:    "debug level",
			level:   "debug",
			want:    zapcore.DebugLevel,
			wantErr: false,
		},
		{
			name:    "info level",
			level:   "info",
			want:    zapcore.InfoLevel,
			wantErr: false,
		},
		{
			name:    "warn level",
			level:   "warn",
			want:    zapcore.WarnLevel,
			wantErr: false,
		},
		{
			name:    "warning level",
			level:   "warning",
			want:    zapcore.WarnLevel,
			wantErr: false,
		},
		{
			name:    "error level",
			level:   "error",
			want:    zapcore.ErrorLevel,
			wantErr: false,
		},
		{
			name:    "case insensitive",
			level:   "DEBUG",
			want:    zapcore.DebugLevel,
			wantErr: false,
		},
		{
			name:    "invalid level",
			level:   "invalid",
			want:    zapcore.InfoLevel,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLevel(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLogLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZapLogger_WithContext(t *testing.T) {
	logger, err := NewLogger("development", nil)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Test with request ID
	ctx := WithRequestID(context.Background(), "req-123")
	loggerWithCtx := logger.WithContext(ctx)
	loggerWithCtx.Info("message with request ID")

	// Test with trace ID
	ctx = WithTraceID(context.Background(), "trace-456")
	loggerWithCtx = logger.WithContext(ctx)
	loggerWithCtx.Info("message with trace ID")

	// Test with multiple context values
	ctx = WithRequestID(context.Background(), "req-789")
	ctx = WithTraceID(ctx, "trace-789")
	loggerWithCtx = logger.WithContext(ctx)
	loggerWithCtx.Info("message with multiple context values")

	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error (may be expected): %v", err)
	}
}

func TestZapLogger_WithFields(t *testing.T) {
	logger, err := NewLogger("development", nil)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	loggerWithFields := logger.WithFields(
		zap.String("service", "test"),
		zap.Int("version", 1),
	)
	loggerWithFields.Info("message with fields")

	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error (may be expected): %v", err)
	}
}

func TestExtractContextFields(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		want   int // expected number of fields
	}{
		{
			name:   "empty context",
			ctx:    context.Background(),
			want:   0,
		},
		{
			name:   "with request ID",
			ctx:    WithRequestID(context.Background(), "req-123"),
			want:   1,
		},
		{
			name:   "with trace ID",
			ctx:    WithTraceID(context.Background(), "trace-456"),
			want:   1,
		},
		{
			name:   "with user ID",
			ctx:    WithUserID(context.Background(), "user-789"),
			want:   1,
		},
		{
			name:   "with all context values",
			ctx:    WithUserID(WithTraceID(WithRequestID(context.Background(), "req-123"), "trace-456"), "user-789"),
			want:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := extractContextFields(tt.ctx)
			if len(fields) != tt.want {
				t.Errorf("extractContextFields() returned %d fields, want %d", len(fields), tt.want)
			}
		})
	}
}

func TestContextHelpers(t *testing.T) {
	// Test WithRequestID
	ctx := WithRequestID(context.Background(), "req-123")
	if requestID, ok := ctx.Value(RequestIDKey).(string); !ok || requestID != "req-123" {
		t.Errorf("WithRequestID() failed, got %v", requestID)
	}

	// Test WithTraceID
	ctx = WithTraceID(context.Background(), "trace-456")
	if traceID, ok := ctx.Value(TraceIDKey).(string); !ok || traceID != "trace-456" {
		t.Errorf("WithTraceID() failed, got %v", traceID)
	}

	// Test WithUserID
	ctx = WithUserID(context.Background(), "user-789")
	if userID, ok := ctx.Value(UserIDKey).(string); !ok || userID != "user-789" {
		t.Errorf("WithUserID() failed, got %v", userID)
	}
}

