package methods

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// TestExtractMethods_WithDefaultFilter verifies that default filter excludes grpc.* services
func TestExtractMethods_WithDefaultFilter(t *testing.T) {
	server := grpc.NewServer()

	// Register standard gRPC services (no external dependencies)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	reflection.Register(server)

	// Use explicit default filter
	methods := Extract(server, DefaultServiceFilter)

	// Should filter out grpc.* infrastructure
	if len(methods) != 0 {
		t.Errorf("expected no methods (grpc.* should be filtered), got: %v", methods)
	}

	t.Log("✓ Default filter correctly excludes grpc.* infrastructure services")
}

// TestExtractMethods_WithNilFilter verifies that nil filter includes all services
func TestExtractMethods_WithNilFilter(t *testing.T) {
	server := grpc.NewServer()

	// Register standard gRPC services
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	reflection.Register(server)

	// Use nil filter (no filtering)
	methods := Extract(server, nil)

	// Should include all services (grpc.health and grpc.reflection)
	if len(methods) == 0 {
		t.Error("expected methods from grpc.* services with nil filter, got none")
	}

	// Verify we got grpc.* methods
	foundGrpcService := false
	for _, method := range methods {
		if len(method) > 5 && method[:6] == "/grpc." {
			foundGrpcService = true
			break
		}
	}
	if !foundGrpcService {
		t.Errorf("expected grpc.* methods with nil filter, got: %v", methods)
	}

	t.Log("✓ Nil filter correctly includes all services")
}

// TestExtractMethods_WithCustomFilter verifies method construction by accepting all services
func TestExtractMethods_WithCustomFilter(t *testing.T) {
	server := grpc.NewServer()

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Use a filter that accepts everything
	acceptAll := func(_ string) bool { return true }
	methods := Extract(server, acceptAll)

	// Should extract health service methods
	if len(methods) == 0 {
		t.Fatal("expected methods to be extracted, got none")
	}

	expectedMethods := map[string]bool{
		"/grpc.health.v1.Health/Check": true,
		"/grpc.health.v1.Health/Watch": true,
	}

	for _, method := range methods {
		t.Logf("Extracted: %s", method)

		// Verify gRPC format
		if len(method) == 0 || method[0] != '/' {
			t.Errorf("method %q doesn't start with /", method)
		}

		slashCount := 0
		for _, ch := range method {
			if ch == '/' {
				slashCount++
			}
		}
		if slashCount < 2 {
			t.Errorf("method %q has %d slashes, expected at least 2", method, slashCount)
		}

		if method[len(method)-1] == '/' {
			t.Errorf("method %q ends with /", method)
		}

		// Mark as found
		delete(expectedMethods, method)
	}

	t.Log("✓ Method construction logic validated")
	t.Log("✓ Methods follow gRPC format")
}

// TestExtractMethods_EmptyServer verifies behavior with no services registered
func TestExtractMethods_EmptyServer(t *testing.T) {
	server := grpc.NewServer()
	methods := Extract(server, nil)

	if len(methods) != 0 {
		t.Errorf("expected no methods from empty server, got: %v", methods)
	}
}

// TestExtractMethods_NilServer verifies behavior with nil server
func TestExtractMethods_NilServer(t *testing.T) {
	methods := Extract(nil, nil)

	if len(methods) != 0 {
		t.Errorf("expected empty slice from nil server, got: %v", methods)
	}
}

// TestNewPatternFilter_ExcludeOnly verifies exclude-only pattern filtering
func TestNewPatternFilter_ExcludeOnly(t *testing.T) {
	server := grpc.NewServer()

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	reflection.Register(server)

	// Exclude grpc.* services
	filter := NewPatternFilter(nil, []string{"grpc."})
	methods := Extract(server, filter)

	// Should filter out all grpc.* services
	if len(methods) != 0 {
		t.Errorf("expected no methods with grpc.* excluded, got: %v", methods)
	}

	t.Log("✓ Exclude-only pattern filter works correctly")
}

// TestNewPatternFilter_IncludeAndExclude verifies combined include/exclude filtering
func TestNewPatternFilter_IncludeAndExclude(t *testing.T) {
	// Create mock service info to test the filter logic
	filter := NewPatternFilter([]string{"api."}, []string{"api.internal."})

	tests := []struct {
		serviceName string
		expected    bool
	}{
		{"api.UserService", true},           // Included
		{"api.internal.DebugService", false}, // Excluded (blacklist overrides)
		{"other.SomeService", false},         // Not in include list
		{"grpc.health.v1.Health", false},     // Not in include list
	}

	for _, tt := range tests {
		result := filter(tt.serviceName)
		if result != tt.expected {
			t.Errorf("filter(%q) = %v, expected %v", tt.serviceName, result, tt.expected)
		}
	}

	t.Log("✓ Include/exclude pattern filter works correctly")
}

// TestNewPatternFilter_MultipleExcludes verifies multiple exclude patterns
func TestNewPatternFilter_MultipleExcludes(t *testing.T) {
	filter := NewPatternFilter(nil, []string{"grpc.", "info.", "diagnostics."})

	tests := []struct {
		serviceName string
		expected    bool
	}{
		{"grpc.health.v1.Health", false},      // Excluded
		{"info.InfoService", false},           // Excluded
		{"diagnostics.DiagnosticsService", false}, // Excluded
		{"principal.PrincipalService", true},  // Not excluded
		{"token.TokenService", true},          // Not excluded
	}

	for _, tt := range tests {
		result := filter(tt.serviceName)
		if result != tt.expected {
			t.Errorf("filter(%q) = %v, expected %v", tt.serviceName, result, tt.expected)
		}
	}

	t.Log("✓ Multiple exclude patterns work correctly")
}
