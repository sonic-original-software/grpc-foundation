package server

import (
	"log/slog"
	"time"

	"git.sonicoriginal.software/grpc-foundation/logging"
	"git.sonicoriginal.software/logger"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// New creates a gRPC server with production-grade defaults.
// The provided logger is used for recovery logging and request logging via interceptors.
// Additional server options can be provided to customize behavior.
func New(log *slog.Logger, opts ...grpc.ServerOption) *grpc.Server {
	if log == nil {
		log = logger.NewNullLogger()
	}

	recoveryHandler := recovery.WithRecoveryHandler(makeRecoveryHandler(log))

	// Logging interceptor using slog
	logOpts := []grpclogging.Option{
		grpclogging.WithLogOnEvents(grpclogging.StartCall, grpclogging.FinishCall),
		grpclogging.WithFieldsFromContext(logging.BaggageFields),
	}

	standardOpts := []grpc.ServerOption{
		// Interceptor chain (outermost to innermost)
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryHandler),
			grpclogging.UnaryServerInterceptor(logging.InterceptorLogger(log), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryHandler),
			grpclogging.StreamServerInterceptor(logging.InterceptorLogger(log), logOpts...),
		),
		// OpenTelemetry tracing via stats handler
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		// Connection management
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     getEnvDurationOrDefault(EnvMaxConnectionIdle, DefaultMaxConnectionIdle),
			MaxConnectionAge:      getEnvDurationOrDefault(EnvMaxConnectionAge, DefaultMaxConnectionAge),
			MaxConnectionAgeGrace: getEnvDurationOrDefault(EnvMaxConnectionAgeGrace, DefaultMaxConnectionAgeGrace),
			Time:                  getEnvDurationOrDefault(EnvKeepAliveTime, DefaultKeepAliveTime),
			Timeout:               getEnvDurationOrDefault(EnvKeepAliveTimeout, DefaultKeepAliveTimeout),
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Minute,
			PermitWithoutStream: true,
		}),
		// Message size limits
		grpc.MaxRecvMsgSize(getEnvIntOrDefault(EnvMaxRecvMsgSize, DefaultMaxRecvMsgSize)),
		grpc.MaxSendMsgSize(getEnvIntOrDefault(EnvMaxSendMsgSize, DefaultMaxSendMsgSize)),
	}

	// Append custom options (can override standard options if needed)
	allOpts := append(standardOpts, opts...)
	return grpc.NewServer(allOpts...)
}

// makeRecoveryHandler creates a panic recovery handler that logs and returns Internal error.
func makeRecoveryHandler(log *slog.Logger) func(p any) error {
	return func(p any) error {
		log.Error("Recovered from panic", slog.Any("panic", p))
		return status.Errorf(codes.Internal, "internal server error")
	}
}
