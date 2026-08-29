package server

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"git.sonicoriginal.software/logger"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Any non-zero pair makes a span context valid, which is all withTrace asks of
// one.
var (
	traceID = trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	spanID  = trace.SpanID{0x01, 0x02, 0x03, 0x04}
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

func TestWithTrace(t *testing.T) {
	t.Run("adds the trace identifiers when the context carries a span", func(t *testing.T) {
		var out bytes.Buffer
		log := slog.New(slog.NewTextHandler(&out, nil))

		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		})

		withTrace(trace.ContextWithSpanContext(t.Context(), spanContext), log).Info("message")

		if !strings.Contains(out.String(), traceID.String()) {
			t.Errorf("output = %q, want it to carry the trace ID", out.String())
		}
		if !strings.Contains(out.String(), spanID.String()) {
			t.Errorf("output = %q, want it to carry the span ID", out.String())
		}
	})

	t.Run("returns the logger unchanged when the context carries no span", func(t *testing.T) {
		log := slog.New(slog.NewTextHandler(io.Discard, nil))

		if withTrace(t.Context(), log) != log {
			t.Error("expected the logger it was given")
		}
	})
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
