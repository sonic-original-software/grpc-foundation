// Package otel provides OpenTelemetry initialization for tracing, metrics, and logging.
package otel

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"git.sonicoriginal.software/grpc-foundation/lifecycle"
	"git.sonicoriginal.software/grpc-foundation/otel/logging"
	"git.sonicoriginal.software/grpc-foundation/otel/metrics"
	"git.sonicoriginal.software/grpc-foundation/otel/tracing"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Init initializes all OpenTelemetry providers: logging first (so tracing and
// metrics initialization can log), then tracing, then metrics.
//
// Returns the logger and a ShutdownFunc that tears down all providers in
// reverse order (metrics, tracing, logging).
func Init(ctx context.Context, serviceName string) (*slog.Logger, lifecycle.ShutdownFunc, error) {
	serviceVersion := os.Getenv("SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "dev"
	}

	// resource.Merge only fails when resources have incompatible schema URLs.
	// Since we use NewSchemaless (no schema URL), this cannot fail.
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	var shutdowns []lifecycle.ShutdownFunc

	// Logging first so tracing and metrics init can log.
	log, logShutdown, err := logging.Init(ctx, res, endpoint, serviceName)
	if err != nil {
		return nil, nil, fmt.Errorf("logging init: %w", err)
	}
	shutdowns = append(shutdowns, logShutdown)

	tracerShutdown, err := tracing.Init(ctx, res, endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("tracing init: %w", err)
	}
	shutdowns = append(shutdowns, tracerShutdown)

	meterShutdown, err := metrics.Init(ctx, res, endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("metrics init: %w", err)
	}
	shutdowns = append(shutdowns, meterShutdown)

	// Shutdown in reverse init order: metrics, tracing, logging last
	// so other providers can log during their shutdown.
	shutdown := func(ctx context.Context) error {
		var errs []error

		for i := len(shutdowns) - 1; i >= 0; i-- {
			if err := shutdowns[i](ctx); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("otel shutdown errors: %v", errs)
		}

		return nil
	}

	return log, shutdown, nil
}
