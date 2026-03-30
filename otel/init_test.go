package otel

import (
	"context"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	t.Run("with all exporters disabled", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "none")

		logger, shutdown, err := Init(t.Context(), "test-service")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		if shutdown == nil {
			t.Fatal("expected non-nil shutdown func")
		}

		shutdownErr := shutdown(t.Context())
		if shutdownErr != nil {
			t.Errorf("shutdown error: %v", shutdownErr)
		}
	})

	t.Run("with empty log exporter defaults to otlp", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")

		logger, shutdown, err := Init(t.Context(), "test-service")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		if shutdown == nil {
			t.Fatal("expected non-nil shutdown func")
		}

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()
		_ = shutdown(ctx)
	})

	t.Run("with unsupported console log exporter", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "console")

		_, _, err := Init(t.Context(), "test-service")
		if err == nil {
			t.Fatal("expected error for unsupported console log exporter")
		}
	})

	t.Run("with unsupported trace exporter", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "invalid")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "none")

		_, _, err := Init(t.Context(), "test-service")
		if err == nil {
			t.Fatal("expected error for unsupported trace exporter")
		}
	})

	t.Run("with unsupported metric exporter", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "invalid")
		t.Setenv("OTEL_LOGS_EXPORTER", "none")

		_, _, err := Init(t.Context(), "test-service")
		if err == nil {
			t.Fatal("expected error for unsupported metric exporter")
		}
	})

	t.Run("with all otlp exporters", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "otlp")
		t.Setenv("OTEL_METRICS_EXPORTER", "otlp")
		t.Setenv("OTEL_LOGS_EXPORTER", "otlp")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")

		logger, shutdown, err := Init(t.Context(), "test-service")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		if shutdown == nil {
			t.Fatal("expected non-nil shutdown func")
		}

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()
		_ = shutdown(ctx)
	})

	t.Run("with empty SERVICE_VERSION defaults to dev", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "none")
		t.Setenv("SERVICE_VERSION", "")

		logger, shutdown, err := Init(t.Context(), "test-service")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		if shutdown == nil {
			t.Fatal("expected non-nil shutdown func")
		}

		shutdownErr := shutdown(t.Context())
		if shutdownErr != nil {
			t.Errorf("shutdown error: %v", shutdownErr)
		}
	})

	t.Run("with custom SERVICE_VERSION from env", func(t *testing.T) {
		t.Setenv("OTEL_TRACES_EXPORTER", "none")
		t.Setenv("OTEL_METRICS_EXPORTER", "none")
		t.Setenv("OTEL_LOGS_EXPORTER", "none")
		t.Setenv("SERVICE_VERSION", "2.5.0")

		logger, shutdown, err := Init(t.Context(), "test-service")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if logger == nil {
			t.Fatal("expected non-nil logger")
		}

		if shutdown == nil {
			t.Fatal("expected non-nil shutdown func")
		}

		shutdownErr := shutdown(t.Context())
		if shutdownErr != nil {
			t.Errorf("shutdown error: %v", shutdownErr)
		}
	})
}
