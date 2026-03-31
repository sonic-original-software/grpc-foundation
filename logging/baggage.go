package logging

import (
	"context"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/baggage"
)

// BaggageFields extracts OTel baggage members from the context as logging fields.
func BaggageFields(ctx context.Context) grpclogging.Fields {
	members := baggage.FromContext(ctx).Members()
	kvPairs := len(members) * 2
	fields := make(grpclogging.Fields, 0, kvPairs)
	for _, m := range members {
		fields = append(fields, m.Key(), m.Value())
	}
	return fields
}
