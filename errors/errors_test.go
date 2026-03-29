package errors

import (
	"context"
	"testing"
	"time"

	"git.sonicoriginal.software/grpc-testing/mocks/tracer"

	"go.opentelemetry.io/otel/attribute"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInvalidArgument(t *testing.T) {
	t.Run("without violations", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		err := InvalidArgument(ctx, "validation failed")
		if err == nil {
			t.Fatal("expected error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("expected gRPC status error")
		}
		if st.Code() != grpccodes.InvalidArgument {
			t.Errorf("expected InvalidArgument, got %v", st.Code())
		}
		if st.Message() != "validation failed" {
			t.Errorf("expected message 'validation failed', got %q", st.Message())
		}

		// End span and check no events recorded
		tt.EndSpan()
		spans := tt.GetSpans()
		if len(spans) > 0 && len(spans[0].Events) > 0 {
			t.Error("expected no span events without violations")
		}
	})

	t.Run("with violations", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		err := InvalidArgument(ctx, "validation failed",
			FieldViolation{Field: "email", Description: "invalid format"},
			FieldViolation{Field: "name", Description: "required"},
		)
		if err == nil {
			t.Fatal("expected error")
		}

		st, _ := status.FromError(err)
		if len(st.Details()) == 0 {
			t.Error("expected error details")
		}

		events := tt.EndSpanAndGetEvents(t)
		if len(events) != 2 {
			t.Errorf("expected 2 span events, got %d", len(events))
		}

		if events[0].Name != "validation_error" {
			t.Errorf("expected event name 'validation_error', got %q", events[0].Name)
		}
		assertAttribute(t, events[0].Attributes, "field", "email")
		assertAttribute(t, events[0].Attributes, "description", "invalid format")
	})
}

func TestAlreadyExists(t *testing.T) {
	tt, ctx := tracer.New(t)
	defer tt.Shutdown(t)

	err := AlreadyExists(ctx, "user already exists", "email", "test@example.com")
	if err == nil {
		t.Fatal("expected error")
	}

	st, _ := status.FromError(err)
	if st.Code() != grpccodes.AlreadyExists {
		t.Errorf("expected AlreadyExists, got %v", st.Code())
	}

	events := tt.EndSpanAndGetEvents(t)
	if len(events) != 1 {
		t.Fatalf("expected 1 span event, got %d", len(events))
	}
	if events[0].Name != "already_exists" {
		t.Errorf("expected event name 'already_exists', got %q", events[0].Name)
	}
	assertAttribute(t, events[0].Attributes, "field", "email")
	assertAttribute(t, events[0].Attributes, "conflict_value", "test@example.com")
}

func TestNotFound(t *testing.T) {
	tt, ctx := tracer.New(t)
	defer tt.Shutdown(t)

	err := NotFound(ctx, "User", "123")
	if err == nil {
		t.Fatal("expected error")
	}

	st, _ := status.FromError(err)
	if st.Code() != grpccodes.NotFound {
		t.Errorf("expected NotFound, got %v", st.Code())
	}
	if st.Message() != "User not found" {
		t.Errorf("expected message 'User not found', got %q", st.Message())
	}

	events := tt.EndSpanAndGetEvents(t)
	if len(events) != 1 {
		t.Fatalf("expected 1 span event, got %d", len(events))
	}
	if events[0].Name != "not_found" {
		t.Errorf("expected event name 'not_found', got %q", events[0].Name)
	}
	assertAttribute(t, events[0].Attributes, "resource_type", "User")
	assertAttribute(t, events[0].Attributes, "resource_id", "123")
}

func TestRateLimited(t *testing.T) {
	t.Run("with custom message", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		err := RateLimited(ctx, "too many requests", 5*time.Second)
		if err == nil {
			t.Fatal("expected error")
		}

		st, _ := status.FromError(err)
		if st.Code() != grpccodes.ResourceExhausted {
			t.Errorf("expected ResourceExhausted, got %v", st.Code())
		}
		if st.Message() != "too many requests" {
			t.Errorf("expected message 'too many requests', got %q", st.Message())
		}

		events := tt.EndSpanAndGetEvents(t)
		if len(events) != 1 {
			t.Fatalf("expected 1 span event, got %d", len(events))
		}
		if events[0].Name != "rate_limited" {
			t.Errorf("expected event name 'rate_limited', got %q", events[0].Name)
		}
		assertAttribute(t, events[0].Attributes, "message", "too many requests")
		assertInt64Attribute(t, events[0].Attributes, "retry_after_ms", 5000)
	})

	t.Run("with empty message", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		err := RateLimited(ctx, "", 1*time.Second)
		if err == nil {
			t.Fatal("expected error")
		}

		st, _ := status.FromError(err)
		if st.Message() != "rate limit exceeded" {
			t.Errorf("expected default message, got %q", st.Message())
		}
	})
}

func TestPreconditionFailed(t *testing.T) {
	tt, ctx := tracer.New(t)
	defer tt.Shutdown(t)

	err := PreconditionFailed(ctx, "cannot delete", "STATE_INVALID", "resource", "must be inactive")
	if err == nil {
		t.Fatal("expected error")
	}

	st, _ := status.FromError(err)
	if st.Code() != grpccodes.FailedPrecondition {
		t.Errorf("expected FailedPrecondition, got %v", st.Code())
	}

	events := tt.EndSpanAndGetEvents(t)
	if len(events) != 1 {
		t.Fatalf("expected 1 span event, got %d", len(events))
	}
	if events[0].Name != "precondition_failed" {
		t.Errorf("expected event name 'precondition_failed', got %q", events[0].Name)
	}
	assertAttribute(t, events[0].Attributes, "violation_type", "STATE_INVALID")
	assertAttribute(t, events[0].Attributes, "subject", "resource")
	assertAttribute(t, events[0].Attributes, "description", "must be inactive")
}

func TestInternal(t *testing.T) {
	tt, ctx := tracer.New(t)
	defer tt.Shutdown(t)

	err := Internal(ctx, "something went wrong", "DATABASE_ERROR", "myapp.storage")
	if err == nil {
		t.Fatal("expected error")
	}

	st, _ := status.FromError(err)
	if st.Code() != grpccodes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}

	events := tt.EndSpanAndGetEvents(t)
	if len(events) != 1 {
		t.Fatalf("expected 1 span event, got %d", len(events))
	}
	if events[0].Name != "internal_error" {
		t.Errorf("expected event name 'internal_error', got %q", events[0].Name)
	}
	assertAttribute(t, events[0].Attributes, "reason", "DATABASE_ERROR")
	assertAttribute(t, events[0].Attributes, "domain", "myapp.storage")
}

func TestWithHelp(t *testing.T) {
	t.Run("with gRPC error", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		originalErr := InvalidArgument(ctx, "validation failed")
		err := WithHelp(originalErr)
		if err == nil {
			t.Fatal("expected error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("expected gRPC status error")
		}
		if st.Code() != grpccodes.InvalidArgument {
			t.Errorf("expected InvalidArgument, got %v", st.Code())
		}
	})

	t.Run("with non-gRPC error", func(t *testing.T) {
		err := WithHelp(context.DeadlineExceeded)
		if err != context.DeadlineExceeded {
			t.Error("expected original error to be returned")
		}
	})
}

func TestWithLocalizedMessage(t *testing.T) {
	t.Run("with gRPC error", func(t *testing.T) {
		tt, ctx := tracer.New(t)
		defer tt.Shutdown(t)

		originalErr := InvalidArgument(ctx, "validation failed")
		err := WithLocalizedMessage(originalErr, "es", "validación fallida")
		if err == nil {
			t.Fatal("expected error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("expected gRPC status error")
		}
		if st.Code() != grpccodes.InvalidArgument {
			t.Errorf("expected InvalidArgument, got %v", st.Code())
		}
	})

	t.Run("with non-gRPC error", func(t *testing.T) {
		err := WithLocalizedMessage(context.DeadlineExceeded, "es", "tiempo agotado")
		if err != context.DeadlineExceeded {
			t.Error("expected original error to be returned")
		}
	})
}

func assertAttribute(t *testing.T, attrs []attribute.KeyValue, key, expected string) {
	t.Helper()
	for _, attr := range attrs {
		if string(attr.Key) == key {
			if attr.Value.AsString() != expected {
				t.Errorf("attribute %q: expected %q, got %q", key, expected, attr.Value.AsString())
			}
			return
		}
	}
	t.Errorf("attribute %q not found", key)
}

func assertInt64Attribute(t *testing.T, attrs []attribute.KeyValue, key string, expected int64) {
	t.Helper()
	for _, attr := range attrs {
		if string(attr.Key) == key {
			if attr.Value.AsInt64() != expected {
				t.Errorf("attribute %q: expected %d, got %d", key, expected, attr.Value.AsInt64())
			}
			return
		}
	}
	t.Errorf("attribute %q not found", key)
}
