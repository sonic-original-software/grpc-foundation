//revive:disable:package-comments
package client

import (
	"log/slog"

	"git.sonicoriginal.software/grpc-foundation/logging"
	"git.sonicoriginal.software/logger"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// New creates a gRPC client connection with standardized options.
//
// Standard options include:
//   - Insecure credentials (for internal service mesh communication)
//   - OpenTelemetry trace propagation (automatic distributed tracing)
//   - Client-side observability logging
//
// If log is nil, a null logger is used that discards all output.
//
// Additional options can be passed to customize the connection:
//
//	conn, err := client.New("service:50051", log, nil, nil,
//	    grpc.WithBlock(),
//	    grpc.WithTimeout(5*time.Second),
//	)
//
// For external services requiring TLS, pass custom credentials:
//
//	tlsConfig := &tls.Config{...}
//	creds := credentials.NewTLS(tlsConfig)
//	conn, err := client.New("external.service.com:443", log, nil, nil,
//	    grpc.WithTransportCredentials(creds), // Overrides insecure default
//	)
func New(
	address string,
	log *slog.Logger,
	unaryInterceptors []grpc.UnaryClientInterceptor,
	streamInterceptors []grpc.StreamClientInterceptor,
	opts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	if log == nil {
		log = logger.NewNullLogger()
	}

	logOpts := []grpclogging.Option{
		grpclogging.WithLogOnEvents(grpclogging.StartCall, grpclogging.FinishCall),
	}

	// Build standard options for internal service communication
	standardOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(
			append(
				[]grpc.UnaryClientInterceptor{
					grpclogging.UnaryClientInterceptor(
						logging.InterceptorLogger(log), logOpts...,
					),
				},
				unaryInterceptors...,
			)...,
		),
		grpc.WithChainStreamInterceptor(
			append(
				[]grpc.StreamClientInterceptor{
					grpclogging.StreamClientInterceptor(
						logging.InterceptorLogger(log), logOpts...,
					),
				},
				streamInterceptors...,
			)...,
		),
	}

	// Append any custom options (later options override earlier ones)
	standardOpts = append(standardOpts, opts...)

	return grpc.NewClient(address, standardOpts...)
}
