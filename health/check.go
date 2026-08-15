//revive:disable:package-comments
package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Check checks if a service is healthy by calling its gRPC health check endpoint
func Check(
	ctx context.Context, conn grpc.ClientConnInterface,
) (grpc_health_v1.HealthCheckResponse_ServingStatus, error) {
	healthClient := grpc_health_v1.NewHealthClient(conn)

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return grpc_health_v1.HealthCheckResponse_UNKNOWN, err
	}

	return resp.Status, nil
}
