package methods

import (
	"slices"
	"testing"

	"google.golang.org/grpc"
)

// serviceInfoStub stands in for a gRPC server, reporting whatever services a
// test needs Extract to see.
type serviceInfoStub struct {
	services map[string]grpc.ServiceInfo
}

func (s *serviceInfoStub) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.services
}

// acceptAll includes every service it is asked about.
func acceptAll(_ string) bool { return true }

// newInfrastructureAndAPI reports one grpc.* infrastructure service and one
// application service, so a filter has something of each kind to act on.
func newInfrastructureAndAPI() *serviceInfoStub {
	return &serviceInfoStub{
		services: map[string]grpc.ServiceInfo{
			"grpc.health.v1.Health": {
				Methods: []grpc.MethodInfo{{Name: "Check"}, {Name: "Watch"}},
			},
			"api.UserService": {
				Methods: []grpc.MethodInfo{{Name: "GetUser"}},
			},
		},
	}
}

func TestExtract(t *testing.T) {
	// Extract sorts its result, so every want below is in sorted order.
	apiMethod := "/api.UserService/GetUser"
	healthMethods := []string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	}
	allMethods := append([]string{apiMethod}, healthMethods...)

	t.Run("excludes grpc infrastructure with the default filter", func(t *testing.T) {
		got := Extract(newInfrastructureAndAPI(), DefaultServiceFilter)
		want := []string{apiMethod}

		if !slices.Equal(got, want) {
			t.Errorf("methods = %v, want %v", got, want)
		}
	})

	t.Run("includes every service with a nil filter", func(t *testing.T) {
		got := Extract(newInfrastructureAndAPI(), nil)

		if !slices.Equal(got, allMethods) {
			t.Errorf("methods = %v, want %v", got, allMethods)
		}
	})

	t.Run("includes every service the filter accepts", func(t *testing.T) {
		got := Extract(newInfrastructureAndAPI(), acceptAll)

		if !slices.Equal(got, allMethods) {
			t.Errorf("methods = %v, want %v", got, allMethods)
		}
	})

	t.Run("excludes services matching an exclude pattern", func(t *testing.T) {
		filter := NewPatternFilter(nil, []string{"grpc."})
		got := Extract(newInfrastructureAndAPI(), filter)
		want := []string{apiMethod}

		if !slices.Equal(got, want) {
			t.Errorf("methods = %v, want %v", got, want)
		}
	})

	t.Run("returns no methods when no services are registered", func(t *testing.T) {
		got := Extract(&serviceInfoStub{}, nil)

		if len(got) != 0 {
			t.Errorf("methods = %v, want none", got)
		}
	})

	t.Run("returns no methods for a nil provider", func(t *testing.T) {
		got := Extract(nil, nil)

		if len(got) != 0 {
			t.Errorf("methods = %v, want none", got)
		}
	})
}

func TestNewPatternFilter(t *testing.T) {
	t.Run("exclude overrides include", func(t *testing.T) {
		filter := NewPatternFilter([]string{"api."}, []string{"api.internal."})

		tests := []struct {
			serviceName string
			want        bool
		}{
			{"api.UserService", true},
			{"api.internal.DebugService", false},
			{"other.SomeService", false},
			{"grpc.health.v1.Health", false},
		}

		for _, tt := range tests {
			if got := filter(tt.serviceName); got != tt.want {
				t.Errorf("filter(%q) = %v, want %v", tt.serviceName, got, tt.want)
			}
		}
	})

	t.Run("excludes every listed prefix", func(t *testing.T) {
		filter := NewPatternFilter(nil, []string{"grpc.", "info.", "diagnostics."})

		tests := []struct {
			serviceName string
			want        bool
		}{
			{"grpc.health.v1.Health", false},
			{"info.InfoService", false},
			{"diagnostics.DiagnosticsService", false},
			{"principal.PrincipalService", true},
			{"token.TokenService", true},
		}

		for _, tt := range tests {
			if got := filter(tt.serviceName); got != tt.want {
				t.Errorf("filter(%q) = %v, want %v", tt.serviceName, got, tt.want)
			}
		}
	})
}
