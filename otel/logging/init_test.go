package logging

import (
	"context"
	"reflect"
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

func TestInit_OTLP_LogFormatNone(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "none")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatJSON(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "json")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatText(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "text")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatFlat(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "flat")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatStructured(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "structured")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatEmpty(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	shutdownErr := shutdown(timeoutCtx)
	_ = shutdownErr
}

func TestInit_OTLP_LogFormatCaseInsensitive(t *testing.T) {
	t.Setenv(envVar, "otlp")
	t.Setenv("LOG_FORMAT", "JSON")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
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
	t.Setenv("LOG_FORMAT", "structured")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
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
	t.Setenv("LOG_FORMAT", "structured")

	ctx := t.Context()
	res := testResource()

	logger, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if logger == nil {
		t.Fatal("expected non-nil logger")
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

	_, shutdown, err := Init(ctx, res, "localhost:4317", "test-service")
	if err == nil {
		t.Fatal("expected error for unsupported exporter")
	}

	if shutdown != nil {
		t.Fatal("expected nil shutdown on error")
	}
}

func TestDefaultAttributeLevels(t *testing.T) {
	levels := DefaultAttributeLevels
	expectedCount := 5

	levelCount := len(levels)
	if levelCount != expectedCount {
		t.Fatalf("expected %d levels, got %d", expectedCount, levelCount)
	}

	expectedLevel0 := []string{}
	if !reflect.DeepEqual(levels[0], expectedLevel0) {
		t.Fatalf("level 0: expected %v, got %v", expectedLevel0, levels[0])
	}

	expectedLevel1 := []string{"source", "service", "address", "component", "trace_id"}
	if !reflect.DeepEqual(levels[1], expectedLevel1) {
		t.Fatalf("level 1: expected %v, got %v", expectedLevel1, levels[1])
	}

	expectedLevel2 := []string{"span_id"}
	if !reflect.DeepEqual(levels[2], expectedLevel2) {
		t.Fatalf("level 2: expected %v, got %v", expectedLevel2, levels[2])
	}

	expectedLevel3 := []string{"method"}
	if !reflect.DeepEqual(levels[3], expectedLevel3) {
		t.Fatalf("level 3: expected %v, got %v", expectedLevel3, levels[3])
	}

	expectedLevel4 := []string{}
	if !reflect.DeepEqual(levels[4], expectedLevel4) {
		t.Fatalf("level 4: expected %v, got %v", expectedLevel4, levels[4])
	}
}
