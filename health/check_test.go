package health

import (
	"context"
	"testing"

	"git.sonicoriginal.software/grpc-foundation/health/tests"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestCheck(t *testing.T) {
	t.Run("serving", func(t *testing.T) {
		servingStatus := grpc_health_v1.HealthCheckResponse_SERVING
		conn := &tests.MockClientConn{
			Response: &grpc_health_v1.HealthCheckResponse{Status: servingStatus},
		}
		got, _ := Check(t.Context(), conn)
		want := servingStatus

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("unavailable", func(t *testing.T) {
		err := status.Error(codes.Unavailable, "health service unavailable")
		conn := &tests.MockClientConn{Err: err}
		_, got := Check(t.Context(), conn)
		want := err

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()

		conn := &tests.MockClientConn{}
		_, got := Check(ctx, conn)
		want := context.DeadlineExceeded

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}
