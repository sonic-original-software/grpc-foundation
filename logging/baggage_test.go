package logging

import (
	"testing"

	"go.opentelemetry.io/otel/baggage"
)

func TestBaggageFields(t *testing.T) {
	member, err := baggage.NewMember("key", "value")
	if err != nil {
		t.Fatalf("creating baggage member: %v", err)
	}

	bag, err := baggage.New(member)
	if err != nil {
		t.Fatalf("creating baggage: %v", err)
	}

	ctx := baggage.ContextWithBaggage(t.Context(), bag)
	fields := BaggageFields(ctx)

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields (key + value), got %d", len(fields))
	}

	if fields[0] != "key" {
		t.Errorf("expected key %q, got %q", "key", fields[0])
	}

	if fields[1] != "value" {
		t.Errorf("expected value %q, got %q", "value", fields[1])
	}
}

func TestBaggageFields_Empty(t *testing.T) {
	fields := BaggageFields(t.Context())

	if len(fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(fields))
	}
}

func TestBaggageFields_Multiple(t *testing.T) {
	first, err := baggage.NewMember("a", "1")
	if err != nil {
		t.Fatalf("creating first member: %v", err)
	}

	second, err := baggage.NewMember("b", "2")
	if err != nil {
		t.Fatalf("creating second member: %v", err)
	}

	bag, err := baggage.New(first, second)
	if err != nil {
		t.Fatalf("creating baggage: %v", err)
	}

	ctx := baggage.ContextWithBaggage(t.Context(), bag)
	fields := BaggageFields(ctx)

	if len(fields) != 4 {
		t.Fatalf("expected 4 fields (2 kv pairs), got %d", len(fields))
	}
}
