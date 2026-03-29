package logging

import (
	"log/slog"
	"testing"

	mockslog "git.sonicoriginal.software/grpc-testing/mocks/slog"
)

func TestInterceptorLogger(t *testing.T) {
	t.Run("logs at correct level with message and fields", func(t *testing.T) {
		var captured []slog.Record
		handler := mockslog.NewCaptureHandler(&captured)
		log := slog.New(handler)

		logger := InterceptorLogger(log)
		ctx := t.Context()

		logger.Log(ctx, 0, "debug message", "key", "value")
		logger.Log(ctx, 4, "info message", "count", 42)

		if len(captured) != 2 {
			t.Fatalf("expected 2 records, got %d", len(captured))
		}
		if captured[0].Level != slog.Level(0) {
			t.Errorf("expected level %v, got %v", slog.Level(0), captured[0].Level)
		}
		if captured[0].Message != "debug message" {
			t.Errorf("expected message %q, got %q", "debug message", captured[0].Message)
		}
		if captured[1].Level != slog.Level(4) {
			t.Errorf("expected level %v, got %v", slog.Level(4), captured[1].Level)
		}
		if captured[1].Message != "info message" {
			t.Errorf("expected message %q, got %q", "info message", captured[1].Message)
		}
	})

	t.Run("logs with empty fields", func(t *testing.T) {
		var captured []slog.Record
		handler := mockslog.NewCaptureHandler(&captured)
		log := slog.New(handler)

		logger := InterceptorLogger(log)
		ctx := t.Context()

		logger.Log(ctx, 4, "simple message")

		if len(captured) != 1 {
			t.Fatalf("expected 1 record, got %d", len(captured))
		}
		if captured[0].Message != "simple message" {
			t.Errorf("expected message %q, got %q", "simple message", captured[0].Message)
		}
	})
}
