// Package logging initializes the OpenTelemetry LoggerProvider and builds
// the slog handler chain.
package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"git.sonicoriginal.software/grpc-foundation/lifecycle"
	"git.sonicoriginal.software/logger/handlers/flat"
	"git.sonicoriginal.software/logger/handlers/json"
	"git.sonicoriginal.software/logger/handlers/structured"
	"git.sonicoriginal.software/logger/handlers/tee"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVar = "OTEL_LOGS_EXPORTER"

// DefaultAttributeLevels defines the standard hierarchical logging structure
// for all services.
var DefaultAttributeLevels = [][]string{
	{}, // Level 0: message only
	{"source", "service", "address", "component", "trace_id"}, // Level 1: service-level context
	{"span_id"}, // Level 2: span context
	{"method"},  // Level 3: operation details
	{},          // Level 4: unknown attributes
}

// Init initializes logging based on OTEL_LOGS_EXPORTER and LOG_FORMAT.
//
// OTEL_LOGS_EXPORTER controls the OTel LoggerProvider:
//   - "otlp": OTLP gRPC exporter to the provided endpoint
//   - "" or "none": no OTel log provider
//
// LOG_FORMAT controls the stdout bridge:
//   - "none": OTel only, no stdout output (production)
//   - "json": tee to stdout in JSON format
//   - "text" or "flat": tee to stdout in flat text format
//   - default (including "structured" or empty): tee to stdout in structured format
//
// Returns the logger, a ShutdownFunc for the LoggerProvider, and any error.
func Init(
	ctx context.Context, res *resource.Resource, serviceName string,
) (*slog.Logger, lifecycle.ShutdownFunc, error) {
	exporterType := os.Getenv(envVar)

	shutdown := func(context.Context) error { return nil }

	switch exporterType {
	case "otlp":
		// OTLP gRPC exporters use lazy connections - they cannot fail at creation time.
		// Endpoint and TLS are configured via OTEL_EXPORTER_OTLP_ENDPOINT env var.
		exporter, _ := otlploggrpc.New(ctx)

		batchProcessor := sdklog.NewBatchProcessor(exporter)
		logProcessor := baggagecopy.NewLogProcessor(nil)
		lp := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(batchProcessor),
			sdklog.WithProcessor(logProcessor),
			sdklog.WithResource(res),
		)
		global.SetLoggerProvider(lp)

		shutdown = lp.Shutdown
	case "none", "":
	default:
		return nil, nil, fmt.Errorf("unsupported %s: %s (supported: otlp, none)", envVar, exporterType)
	}

	// Create handler AFTER SetLoggerProvider — otelslog captures the global
	// provider at creation time, not lazily per record.
	otelHandler := otelslog.NewHandler(serviceName)
	serviceAttr := slog.String("service", serviceName)

	var stdoutHandler slog.Handler
	switch strings.ToLower(os.Getenv("LOG_FORMAT")) {
	case "json":
		stdoutHandler = json.NewHandler()
	case "text", "flat":
		stdoutHandler = flat.NewHandler()
	case "none":
		otelLogger := slog.New(otelHandler)
		logger := otelLogger.With(serviceAttr)
		return logger, shutdown, nil
	default: // "structured" or empty
		stdoutHandler = structured.NewHandler(&structured.Options{
			AttributeLevels: DefaultAttributeLevels,
		})
	}

	teeHandler := tee.NewHandler(otelHandler, stdoutHandler)
	teeLogger := slog.New(teeHandler)
	logger := teeLogger.With(serviceAttr)

	return logger, shutdown, nil
}
