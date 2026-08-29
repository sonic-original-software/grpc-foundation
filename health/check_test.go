package health

import (
	"context"
	"testing"

	"git.sonicoriginal.software/grpc-testing/mocks/clientconn"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// respondServing answers the health check with SERVING.
func respondServing(_ context.Context, _ string, _, reply any) error {
	reply.(*grpc_health_v1.HealthCheckResponse).Status =
		grpc_health_v1.HealthCheckResponse_SERVING

	return nil
}

func TestCheck(t *testing.T) {
	t.Run("serving", func(t *testing.T) {
		conn := &clientconn.Mock{InvokeFn: respondServing}
		got, _ := Check(t.Context(), conn)
		want := grpc_health_v1.HealthCheckResponse_SERVING

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("unavailable", func(t *testing.T) {
		err := status.Error(codes.Unavailable, "health service unavailable")
		conn := &clientconn.Mock{Err: err}
		_, got := Check(t.Context(), conn)
		want := err

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()

		conn := &clientconn.Mock{}
		_, got := Check(ctx, conn)
		want := context.DeadlineExceeded

		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
}
