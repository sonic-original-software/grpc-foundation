//revive:disable:package-comments
package health

import (
	"context"
	"time"

	"git.sonicoriginal.software/grpc-foundation/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Check checks if a service is healthy by calling its gRPC health check endpoint
func Check(ctx context.Context, address string) string {
	if address == "" {
		return ""
	}

	conn, err := client.New(
		address,
		nil,
		nil,
		nil,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err.Error()
	}

	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)

	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := healthClient.Check(checkCtx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err.Error()
	}

	return resp.Status.String()
}
