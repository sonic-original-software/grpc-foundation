//revive:disable:package-comments
package errors

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

// FieldViolation represents a validation error for a specific field
type FieldViolation struct {
	Field       string
	Description string
}

// InvalidArgument returns a gRPC error with structured field violations
// Use this for validation errors where specific fields failed validation
func InvalidArgument(ctx context.Context, message string, violations ...FieldViolation) error {
	st := status.New(codes.InvalidArgument, message)

	if len(violations) > 0 {
		fieldViolations := make([]*errdetails.BadRequest_FieldViolation, len(violations))
		for i, v := range violations {
			fieldViolations[i] = &errdetails.BadRequest_FieldViolation{
				Field:       v.Field,
				Description: v.Description,
			}
			// Record span event for each field violation
			trace.SpanFromContext(ctx).AddEvent("validation_error", trace.WithAttributes(
				attribute.String("field", v.Field),
				attribute.String("description", v.Description),
			))
		}

		br := &errdetails.BadRequest{FieldViolations: fieldViolations}
		st, _ = st.WithDetails(br)
	}

	return st.Err()
}

// AlreadyExists returns a gRPC error indicating a resource already exists
// Includes precondition failure details for the conflicting field
func AlreadyExists(ctx context.Context, message string, field string, conflictValue string) error {
	st := status.New(codes.AlreadyExists, message)

	pf := &errdetails.PreconditionFailure{
		Violations: []*errdetails.PreconditionFailure_Violation{
			{
				Type:        "RESOURCE_ALREADY_EXISTS",
				Subject:     field,
				Description: message + ": " + conflictValue,
			},
		},
	}
	st, _ = st.WithDetails(pf)

	trace.SpanFromContext(ctx).AddEvent("already_exists", trace.WithAttributes(
		attribute.String("field", field),
		attribute.String("conflict_value", conflictValue),
	))

	return st.Err()
}

// NotFound returns a gRPC error indicating a resource was not found
// Includes resource info for the missing resource
func NotFound(ctx context.Context, resourceType string, resourceID string) error {
	message := resourceType + " not found"
	st := status.New(codes.NotFound, message)

	ri := &errdetails.ResourceInfo{
		ResourceType: resourceType,
		ResourceName: resourceID,
		Description:  message,
	}
	st, _ = st.WithDetails(ri)

	trace.SpanFromContext(ctx).AddEvent("not_found", trace.WithAttributes(
		attribute.String("resource_type", resourceType),
		attribute.String("resource_id", resourceID),
	))

	return st.Err()
}

// RateLimited returns a gRPC error indicating rate limiting
// Includes retry information with the recommended retry delay
func RateLimited(ctx context.Context, message string, retryAfter time.Duration) error {
	if message == "" {
		message = "rate limit exceeded"
	}
	st := status.New(codes.ResourceExhausted, message)

	ri := &errdetails.RetryInfo{
		RetryDelay: durationpb.New(retryAfter),
	}
	st, _ = st.WithDetails(ri)

	qf := &errdetails.QuotaFailure{
		Violations: []*errdetails.QuotaFailure_Violation{
			{
				Subject:     "rate_limit",
				Description: message,
			},
		},
	}
	st, _ = st.WithDetails(qf)

	trace.SpanFromContext(ctx).AddEvent("rate_limited", trace.WithAttributes(
		attribute.String("message", message),
		attribute.Int64("retry_after_ms", retryAfter.Milliseconds()),
	))

	return st.Err()
}

// PreconditionFailed returns a gRPC error for failed preconditions
// Use when an operation cannot proceed due to the current state of the resource
func PreconditionFailed(ctx context.Context, message string, violationType string, subject string, description string) error {
	st := status.New(codes.FailedPrecondition, message)

	pf := &errdetails.PreconditionFailure{
		Violations: []*errdetails.PreconditionFailure_Violation{
			{
				Type:        violationType,
				Subject:     subject,
				Description: description,
			},
		},
	}
	st, _ = st.WithDetails(pf)

	trace.SpanFromContext(ctx).AddEvent("precondition_failed", trace.WithAttributes(
		attribute.String("violation_type", violationType),
		attribute.String("subject", subject),
		attribute.String("description", description),
	))

	return st.Err()
}

// Internal returns a gRPC internal error with debug info
// Includes error info with reason and domain for categorization
func Internal(ctx context.Context, message string, reason string, domain string) error {
	st := status.New(codes.Internal, message)

	ei := &errdetails.ErrorInfo{
		Reason: reason,
		Domain: domain,
	}
	st, _ = st.WithDetails(ei)

	trace.SpanFromContext(ctx).AddEvent("internal_error", trace.WithAttributes(
		attribute.String("reason", reason),
		attribute.String("domain", domain),
	))

	return st.Err()
}

// WithHelp adds help links to an existing error
// Useful for providing documentation links for common errors
func WithHelp(err error, links ...*errdetails.Help_Link) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	help := &errdetails.Help{Links: links}
	st, _ = st.WithDetails(help)

	return st.Err()
}

// WithLocalizedMessage adds a localized message to an existing error
// Useful for internationalization support
func WithLocalizedMessage(err error, locale string, localizedMessage string) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	lm := &errdetails.LocalizedMessage{
		Locale:  locale,
		Message: localizedMessage,
	}
	st, _ = st.WithDetails(lm)

	return st.Err()
}
