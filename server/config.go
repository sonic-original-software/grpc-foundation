// Package server provides gRPC server creation and lifecycle management.
package server

import (
	"os"
	"strconv"
	"time"
)

// Default gRPC server connection settings.
const (
	DefaultAddress = ":50051"
	DefaultVersion = "dev"

	DefaultMaxConnectionIdle     = 5 * time.Minute
	DefaultMaxConnectionAge      = 10 * time.Minute
	DefaultMaxConnectionAgeGrace = 5 * time.Second
	DefaultKeepAliveTime         = 2 * time.Minute
	DefaultKeepAliveTimeout      = 20 * time.Second
	DefaultMaxRecvMsgSize        = 4 * 1024 * 1024
	DefaultMaxSendMsgSize        = 4 * 1024 * 1024
)

// Environment variable names for runtime configuration.
const (
	EnvServerName            = "GRPC_SERVER_NAME"
	EnvServerAddress         = "GRPC_SERVER_ADDRESS"
	EnvServerVersion         = "GRPC_SERVER_VERSION"
	EnvMaxConnectionIdle     = "GRPC_MAX_CONNECTION_IDLE"
	EnvMaxConnectionAge      = "GRPC_MAX_CONNECTION_AGE"
	EnvMaxConnectionAgeGrace = "GRPC_MAX_CONNECTION_AGE_GRACE"
	EnvKeepAliveTime         = "GRPC_KEEPALIVE_TIME"
	EnvKeepAliveTimeout      = "GRPC_KEEPALIVE_TIMEOUT"
	EnvMaxRecvMsgSize        = "GRPC_MAX_RECV_MSG_SIZE"
	EnvMaxSendMsgSize        = "GRPC_MAX_SEND_MSG_SIZE"
)

// Address returns the server address from environment or the default value.
func Address() string {
	return getEnvOrDefault(EnvServerAddress, DefaultAddress)
}

// Version returns the server version from environment or the default value.
func Version() string {
	return getEnvOrDefault(EnvServerVersion, DefaultVersion)
}

// Name returns the server name from environment or the caller's fallback.
// There is no universal default: the fallback is the server's own identity,
// and the environment overrides it per deployment.
func Name(fallback string) string {
	return getEnvOrDefault(EnvServerName, fallback)
}

// getEnvOrDefault returns the environment variable value or default if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvDurationOrDefault parses a duration from environment variable or returns default.
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvIntOrDefault parses an int from environment variable or returns default.
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
