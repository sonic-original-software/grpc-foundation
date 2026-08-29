package server

import (
	"context"
	"log/slog"

	"git.sonicoriginal.software/logger"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// withTrace adds the trace identifiers to log when ctx carries a valid span
// context. The stats handler starts the server span before the interceptor
// chain runs, so a request that is being traced already has one here.
//
// otelslog puts these on the OTel record it exports, but the handlers writing
// to stdout never see them, so they have to be slog attributes to be printed.
func withTrace(ctx context.Context, log *slog.Logger) *slog.Logger {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return log
	}

	return log.With(
		slog.String("trace_id", spanContext.TraceID().String()),
		slog.String("span_id", spanContext.SpanID().String()),
	)
}

// makeLoggerUnaryInterceptor puts log into every request context. gRPC builds
// each RPC's context from the incoming stream rather than from the context the
// process started with, so a handler calling logger.FromContext would otherwise
// fall back to slog.Default() and its records would never reach the collector.
func makeLoggerUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any,
		_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (any, error) {
		return handler(logger.ContextWithLogger(ctx, withTrace(ctx, log)), req)
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
		ctx := logger.ContextWithLogger(ss.Context(), withTrace(ss.Context(), log))

		return handler(srv, &loggerStream{ServerStream: ss, ctx: ctx})
	}
}
