package server

import (
	"context"
	"log/slog"
	"os"
	"syscall"
	"testing"
	"time"
)

// stopperStub records that the handler stopped the server. Closing the channel
// gives a test something to wait on instead of a sleep or the handler's return.
type stopperStub struct {
	stopped chan struct{}
}

func newStopperStub() *stopperStub {
	return &stopperStub{stopped: make(chan struct{})}
}

func (s *stopperStub) GracefulStop() { close(s.stopped) }

// blockingStopperStub stands in for a server whose in-flight requests never
// finish, so the handler has to stop waiting on its own.
type blockingStopperStub struct {
	release chan struct{}
}

func (s *blockingStopperStub) GracefulStop() { <-s.release }

// cancelStub records that the handler cancelled the background context. It is
// safe to read called once the stopper has fired, because the handler cancels
// before it stops the server.
type cancelStub struct {
	cancel context.CancelFunc
	called bool
}

func (c *cancelStub) Cancel() {
	c.called = true

	c.cancel()
}

// await fails the test if signal does not fire. The ceiling turns a handler
// that never gets there into a failure rather than a hang.
func await(t *testing.T, signal <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(5 * time.Second):
		t.Fatal(message)
	}
}

func TestHandleGracefulShutdown(t *testing.T) {
	t.Run("stops the server on context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		stopper := newStopperStub()

		go HandleGracefulShutdown(ctx, cancel, slog.New(slog.DiscardHandler), stopper, 5*time.Second)

		cancel()

		await(t, stopper.stopped, "server was not stopped")
	})

	t.Run("cancels the background context before stopping", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		stopper := newStopperStub()
		tracked := &cancelStub{cancel: cancel}

		go HandleGracefulShutdown(ctx, tracked.Cancel, slog.New(slog.DiscardHandler), stopper, 5*time.Second)

		cancel()

		await(t, stopper.stopped, "server was not stopped")

		if !tracked.called {
			t.Fatal("cancel function should have been called")
		}
	})

	t.Run("returns when in-flight requests do not finish in time", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		stopper := &blockingStopperStub{release: make(chan struct{})}
		defer close(stopper.release)

		returned := make(chan struct{})
		go func() {
			HandleGracefulShutdown(
				ctx, cancel, slog.New(slog.DiscardHandler), stopper, time.Millisecond,
			)
			close(returned)
		}()

		cancel()

		await(t, returned, "handler did not stop waiting for the server")
	})

	t.Run("stops the server on OS signal", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		stopper := newStopperStub()

		go HandleGracefulShutdown(ctx, cancel, slog.New(slog.DiscardHandler), stopper, 5*time.Second)

		// The handler has to reach signal.Notify before the signal is sent, or
		// SIGTERM ends the test process instead. Nothing observable marks that
		// point, so the wait is a sleep.
		time.Sleep(50 * time.Millisecond)

		proc, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatalf("unexpected error finding process: %v", err)
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			t.Fatalf("unexpected error sending signal: %v", err)
		}

		await(t, stopper.stopped, "server was not stopped")
	})
}
