package client

import (
	"context"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	t.Run("creates client with minimal options", func(t *testing.T) {
		conn, err := New("localhost:50051", nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
	})

	t.Run("creates client with logger", func(t *testing.T) {
		conn, err := New("localhost:50051", slog.Default(), nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
	})

	t.Run("creates client with custom unary interceptors", func(t *testing.T) {
		called := false
		customInterceptor := func(
			ctx context.Context,
			method string,
			req, reply any,
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			called = true
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		conn, err := New("localhost:50051", nil, []grpc.UnaryClientInterceptor{customInterceptor}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
		// Note: called remains false since we don't make an actual RPC call
		if called {
			t.Error("expected called to be false")
		}
	})

	t.Run("creates client with custom stream interceptors", func(t *testing.T) {
		called := false
		customInterceptor := func(
			ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			streamer grpc.Streamer,
			opts ...grpc.CallOption,
		) (grpc.ClientStream, error) {
			called = true
			return streamer(ctx, desc, cc, method, opts...)
		}

		conn, err := New("localhost:50051", nil, nil, []grpc.StreamClientInterceptor{customInterceptor})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
		// Note: called remains false since we don't make an actual RPC call
		if called {
			t.Error("expected called to be false")
		}
	})

	t.Run("creates client with additional dial options", func(t *testing.T) {
		customOpt := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024))
		conn, err := New("localhost:50051", nil, nil, nil, customOpt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
	})

	t.Run("creates client with all options combined", func(t *testing.T) {
		unaryInterceptor := func(
			ctx context.Context,
			method string,
			req, reply any,
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		streamInterceptor := func(
			ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			streamer grpc.Streamer,
			opts ...grpc.CallOption,
		) (grpc.ClientStream, error) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		customOpt := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(2048))

		conn, err := New(
			"localhost:50051",
			slog.Default(),
			[]grpc.UnaryClientInterceptor{unaryInterceptor},
			[]grpc.StreamClientInterceptor{streamInterceptor},
			customOpt,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Fatal("expected conn to be non-nil")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("unexpected error closing conn: %v", err)
		}
	})
}
