# grpc-foundation

Core infrastructure library for building production-ready gRPC services in Go.
Provides server and client creation, observability, error handling, and common
utilities.

## Overview

This library consolidates the essential building blocks needed for gRPC
microservices, including observability, structured error handling, standardized
logging, and authentication context propagation.

## Installation

```bash
go get git.sonicoriginal.software/grpc-foundation
```

## Packages

### server

Creates gRPC servers with production defaults including connection management,
keepalive settings, and configurable timeouts. Configuration can be provided
through environment variables or directly via configuration structs.

### client

Provides client connection creation with standard options and connection pooling
support.

### interceptors

Collection of gRPC interceptors (middleware) for both unary and streaming RPCs:

- **Observability**: Request logging with distributed tracing context (trace_id,
  span_id)
- **Metrics**: OpenTelemetry metrics recording for request counts, durations,
  and errors
- **Recovery**: Panic recovery with graceful error responses
- **AuthPropagation**: Client-side interceptor for propagating authentication
  context across service boundaries

Interceptors integrate with OpenTelemetry for distributed tracing and metrics
collection.

### errors

Structured error construction following Google's error model with rich error
details. Supports field validation errors, resource conflicts, not found
conditions, rate limiting with retry information, precondition failures, and
internal errors with debugging context.

Errors can be enhanced with help links and localized messages for client
consumption.

### logging

Standardized logger creation for gRPC services with support for multiple output
formats (JSON, flat, structured) controlled via environment variables. Defines
hierarchical logging structure optimized for gRPC service context including
trace IDs, span IDs, and method names.

### metadata

Utilities for working with gRPC request metadata (context key-value pairs).
Handles extraction and injection of authentication context (principal IDs) from
incoming and outgoing request metadata.

### methods

Extracts registered method names from gRPC server instances. Useful for service
discovery, method registration, and filtering infrastructure services from
application services.

### health

Client utilities for performing health checks on remote gRPC services.

## Usage

Import the packages you need from the library. The server and client packages
provide the foundation, while interceptors add observability and reliability
features. Use the errors package for consistent error responses and the metadata
package for authentication propagation.

## Requirements

- Go 1.21 or later
- OpenTelemetry libraries for observability features
- gRPC Go libraries

## License

Apache License 2.0
