package tracing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func testResource() *resource.Resource {
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName("test-service")),
	)
	return res
}

func TestInit_OTLP(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")

	ctx := t.Context()
	res := testResource()

	shutdown, err := Init(ctx, res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_None(t *testing.T) {
	t.Setenv(envVar, "none")

	ctx := t.Context()
	res := testResource()

	shutdown, err := Init(ctx, res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	shutdownErr := shutdown(ctx)
	if shutdownErr != nil {
		t.Fatalf("expected nil error from noop shutdown, got %v", shutdownErr)
	}
}

func TestInit_Empty(t *testing.T) {
	t.Setenv(envVar, "")

	ctx := t.Context()
	res := testResource()

	shutdown, err := Init(ctx, res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	shutdownErr := shutdown(ctx)
	if shutdownErr != nil {
		t.Fatalf("expected nil error from noop shutdown, got %v", shutdownErr)
	}
}

func TestInit_Invalid(t *testing.T) {
	t.Setenv(envVar, "invalid")

	ctx := t.Context()
	res := testResource()

	shutdown, err := Init(ctx, res)
	if err == nil {
		t.Fatal("expected error for unsupported exporter")
	}

	if shutdown != nil {
		t.Fatal("expected nil shutdown on error")
	}
}
