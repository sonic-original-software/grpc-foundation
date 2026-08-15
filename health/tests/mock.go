//revive:disable:package-comments
package tests

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// MockClientConn implements grpc.ClientConnInterface for testing.
type MockClientConn struct {
	grpc.ClientConnInterface
	// Response is written to the unary RPC response on success.
	Response *grpc_health_v1.HealthCheckResponse
	// Err is returned by Invoke when set.
	Err error
	// Ctx is the context received by Invoke.
	Ctx context.Context
	// Method is the full RPC method received by Invoke.
	Method string
	// Request is the health check request received by Invoke.
	Request *grpc_health_v1.HealthCheckRequest
	// Calls tracks how many unary RPCs were invoked.
	Calls int
}

// Invoke records the unary RPC and returns the configured response or error.
func (c *MockClientConn) Invoke(
	ctx context.Context, method string, args, reply any, _ ...grpc.CallOption,
) error {
	c.Calls++
	c.Ctx = ctx
	c.Method = method
	c.Request = args.(*grpc_health_v1.HealthCheckRequest)

	if err := ctx.Err(); err != nil {
		return err
	}
	if c.Err != nil {
		return c.Err
	}
	if c.Response != nil {
		response := reply.(*grpc_health_v1.HealthCheckResponse)
		response.Status = c.Response.Status
	}

	return nil
}
