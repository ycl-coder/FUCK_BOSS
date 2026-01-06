// Package middleware provides gRPC interceptors for logging and recovery.
package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"fuck_boss/backend/internal/infrastructure/logger"
)

// LoggingInterceptor returns a gRPC unary server interceptor that logs all requests.
// It logs request method, duration, and response status.
func LoggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Generate request ID
		requestID := generateRequestID()
		ctx = logger.WithRequestID(ctx, requestID)

		// Get logger with context
		ctxLogger := log.WithContext(ctx)

		// Extract metadata for logging
		md, _ := metadata.FromIncomingContext(ctx)
		clientIP := extractClientIPFromMetadata(md)

		// Log request start
		ctxLogger.Info("gRPC request started",
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.Any("metadata", md),
		)

		// Handle request and measure duration
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		// Extract status code
		statusCode := "OK"
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code().String()
			} else {
				statusCode = "Unknown"
			}
		}

		// Log request completion
		if err != nil {
			ctxLogger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", statusCode),
				zap.Error(err),
				zap.Duration("duration", duration),
			)
		} else {
			ctxLogger.Info("gRPC request completed",
				zap.String("method", info.FullMethod),
				zap.String("status_code", statusCode),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

// generateRequestID generates a unique request ID.
// Uses timestamp and a simple counter for uniqueness.
// In production, you might want to use UUID or other unique identifier.
func generateRequestID() string {
	// Simple implementation using timestamp and nanosecond
	// In production, use a proper UUID library like github.com/google/uuid
	now := time.Now()
	return now.Format("20060102150405") + "-" +
		formatInt64(now.UnixNano()%1000000, 6) + "-" +
		formatInt64(int64(len(now.String())), 4)
}

// formatInt64 formats an int64 to a string with leading zeros.
func formatInt64(n int64, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = string(rune('0'+(n%10))) + s
		n /= 10
	}
	return s
}

// extractClientIPFromMetadata extracts client IP from gRPC metadata.
func extractClientIPFromMetadata(md metadata.MD) string {
	if md == nil {
		return "unknown"
	}

	// Try X-Forwarded-For first
	if forwardedFor := md.Get("x-forwarded-for"); len(forwardedFor) > 0 {
		return forwardedFor[0]
	}

	// Try X-Real-IP
	if realIP := md.Get("x-real-ip"); len(realIP) > 0 {
		return realIP[0]
	}

	return "unknown"
}
