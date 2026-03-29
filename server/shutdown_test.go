package server

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestHandleGracefulShutdown(t *testing.T) {
	t.Run("shuts down on context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())

		server := grpc.NewServer()

		var wg sync.WaitGroup

		wg.Go(func() {
			HandleGracefulShutdown(ctx, cancel, slog.Default(), server, 5*time.Second)
		})

		// Trigger shutdown via context cancellation
		cancel()

		// Wait for shutdown with timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for shutdown")
		}
	})

	t.Run("cancel func is called during shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())

		server := grpc.NewServer()

		// Track whether the cancel was called
		cancelCalled := false
		wrappedCancel := func() {
			cancelCalled = true
			cancel()
		}

		var wg sync.WaitGroup

		wg.Go(func() {
			HandleGracefulShutdown(ctx, wrappedCancel, slog.Default(), server, 5*time.Second)
		})

		// Trigger shutdown via the original cancel (not the wrapped one)
		// This simulates external cancellation
		cancel()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			if !cancelCalled {
				t.Fatal("cancel function should have been called")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for shutdown")
		}
	})

	t.Run("shuts down on OS signal", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		server := grpc.NewServer()

		var wg sync.WaitGroup

		wg.Go(func() {
			HandleGracefulShutdown(ctx, cancel, slog.Default(), server, 5*time.Second)
		})

		// Wait for handler to be ready to receive signals
		time.Sleep(50 * time.Millisecond)

		// Send SIGTERM to trigger signal-based shutdown
		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatalf("unexpected error finding process: %v", err)
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			t.Fatalf("unexpected error sending signal: %v", err)
		}

		// Wait for shutdown with timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Success - shutdown triggered by OS signal
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for signal-triggered shutdown")
		}
	})
}
