// Package middleware provides gRPC interceptors for logging and recovery.
package middleware

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fuck_boss/backend/internal/infrastructure/logger"
)

// RecoveryInterceptor returns a gRPC unary server interceptor that recovers from panics.
// It logs the panic with stack trace and returns an internal error.
func RecoveryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Use a named return to allow setting err from defer
		defer func() {
			if r := recover(); r != nil {
				// Get logger with context
				ctxLogger := log.WithContext(ctx)

				// Log panic with stack trace
				ctxLogger.Error("gRPC panic recovered",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)

				// Convert panic to gRPC error
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		// Call handler
		return handler(ctx, req)
	}
}
