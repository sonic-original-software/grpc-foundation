package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulStopper stops a gRPC server once its in-flight requests finish
type GracefulStopper interface {
	GracefulStop()
}

// HandleGracefulShutdown waits for shutdown signals and orchestrates cleanup.
// It blocks until either an OS signal (SIGINT/SIGTERM) is received or the context is cancelled.
// The shutdown sequence is:
//  1. Cancel the context to stop background goroutines
//  2. Gracefully stop the gRPC server
//
// OTel provider shutdown is handled by the caller via defer.
func HandleGracefulShutdown(
	ctx context.Context,
	cancel context.CancelFunc,
	log *slog.Logger,
	grpcServer GracefulStopper,
	cleanupTimeout time.Duration,
) {
	shutdownLog := log.With(slog.String("component", "shutdown-handler"))
	shutdownLog.Info("Awaiting shutdown signal")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for either OS signal or context cancellation
	select {
	case <-sigChan:
		shutdownLog.Info("Shutdown signal received")
	case <-ctx.Done():
		shutdownLog.Info("Context cancelled")
	}
	shutdownLog.Info("Initiating graceful shutdown")

	// 1. Cancel context to stop background goroutines (grpcd poller)
	shutdownLog.Info("Cancelling background context")
	cancel()
	shutdownLog.Info("Background context cancelled")

	// 2. Stop gRPC server (stop accepting new requests, finish in-flight)
	shutdownLog.Info("Stopping gRPC server")
	grpcServer.GracefulStop()
	shutdownLog.Info("gRPC server stopped")
}
