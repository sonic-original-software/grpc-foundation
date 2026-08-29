package server

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"git.sonicoriginal.software/logger"

	"google.golang.org/grpc"
)

// unaryResponse is what handlerStub.unary returns, so the interceptor can be
// shown to pass a handler's response back untouched.
const unaryResponse = "response"

// handlerStub records the logger a handler finds in the context it is given.
type handlerStub struct {
	log *slog.Logger
}

func (h *handlerStub) unary(ctx context.Context, _ any) (any, error) {
	h.log = logger.FromContext(ctx)

	return unaryResponse, nil
}

func (h *handlerStub) stream(_ any, ss grpc.ServerStream) error {
	h.log = logger.FromContext(ss.Context())

	return nil
}

// serverStreamStub stands in for the stream gRPC hands the interceptor.
type serverStreamStub struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamStub) Context() context.Context {
	return s.ctx
}

func TestMakeLoggerUnaryInterceptor(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := &handlerStub{}

	interceptor := makeLoggerUnaryInterceptor(log)

	resp, err := interceptor(t.Context(), nil, nil, handler.unary)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != unaryResponse {
		t.Errorf("response = %v, want %q", resp, unaryResponse)
	}
	if handler.log != log {
		t.Error("the handler did not receive the logger")
	}
}

func TestMakeLoggerStreamInterceptor(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := &handlerStub{}
	stream := &serverStreamStub{ctx: t.Context()}

	interceptor := makeLoggerStreamInterceptor(log)

	if err := interceptor(nil, stream, nil, handler.stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if handler.log != log {
		t.Error("the handler did not receive the logger")
	}
}
