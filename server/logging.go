package server

import (
	"context"
	"log/slog"

	"git.sonicoriginal.software/logger"

	"google.golang.org/grpc"
)

// makeLoggerUnaryInterceptor puts log into every request context. gRPC builds
// each RPC's context from the incoming stream rather than from the context the
// process started with, so a handler calling logger.FromContext would otherwise
// fall back to slog.Default() and its records would never reach the collector.
func makeLoggerUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any,
		_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		return handler(logger.ContextWithLogger(ctx, log), req)
	}
}

// loggerStream carries a context the interceptor has attached the logger to.
// A streaming handler reads its context from the stream rather than from an
// argument, so replacing Context is the only way to reach one.
type loggerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the context carrying the logger.
func (s *loggerStream) Context() context.Context {
	return s.ctx
}

// makeLoggerStreamInterceptor is the streaming counterpart to
// makeLoggerUnaryInterceptor.
func makeLoggerStreamInterceptor(log *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any, ss grpc.ServerStream,
		_ *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		ctx := logger.ContextWithLogger(ss.Context(), log)

		return handler(srv, &loggerStream{ServerStream: ss, ctx: ctx})
	}
}
