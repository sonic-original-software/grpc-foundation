// Package tracing initializes the OpenTelemetry TracerProvider.
package tracing

import (
	"context"
	"fmt"
	"os"

	"git.sonicoriginal.software/grpc-foundation/lifecycle"

	"go.opentelemetry.io/contrib/processors/baggagecopy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const envVar = "OTEL_TRACES_EXPORTER"

// Init initializes the OTel TracerProvider based on OTEL_TRACES_EXPORTER.
//
//   - "otlp": OTLP gRPC exporter to the provided endpoint
//   - "" or "none": no tracing (returns noop shutdown)
//
// Returns a ShutdownFunc to flush and close the provider.
func Init(
	ctx context.Context, res *resource.Resource,
) (lifecycle.ShutdownFunc, error) {
	exporterType := os.Getenv(envVar)

	if exporterType == "" || exporterType == "none" {
		noop := func(context.Context) error { return nil }
		return noop, nil
	}

	if exporterType != "otlp" {
		return nil, fmt.Errorf("unsupported %s: %s (supported: otlp, none)", envVar, exporterType)
	}

	// OTLP gRPC exporters use lazy connections - they cannot fail at creation time.
	// Errors occur later during export if the collector is unreachable.
	// Endpoint and TLS are configured via OTEL_EXPORTER_OTLP_ENDPOINT env var.
	exporter, _ := otlptracegrpc.New(ctx)

	spanProcessor := baggagecopy.NewSpanProcessor(nil)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
	)
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	return tp.Shutdown, nil
}
