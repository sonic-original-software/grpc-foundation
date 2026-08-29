//revive:disable:package-comments
package methods

import (
	"sort"
	"strings"

	"google.golang.org/grpc"
)

// ServiceFilter determines whether a service should be included
type ServiceFilter func(serviceName string) bool

// ServiceInfoProvider reports the services registered on a gRPC server
type ServiceInfoProvider interface {
	GetServiceInfo() map[string]grpc.ServiceInfo
}

// DefaultServiceFilter filters out gRPC infrastructure services (health, reflection, etc.)
func DefaultServiceFilter(serviceName string) bool {
	return !strings.HasPrefix(serviceName, "grpc.")
}

// NewPatternFilter creates a ServiceFilter with both include and exclude lists
// includePrefixes: if non-empty, only services matching these prefixes are included (whitelist)
// excludePrefixes: services matching these prefixes are always excluded (blacklist overrides whitelist)
func NewPatternFilter(includePrefixes, excludePrefixes []string) ServiceFilter {
	return func(serviceName string) bool {
		// If include list provided, service must match at least one include
		if len(includePrefixes) > 0 {
			matched := false
			for _, prefix := range includePrefixes {
				if strings.HasPrefix(serviceName, prefix) {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}

		// Then check excludes (overrides includes)
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(serviceName, prefix) {
				return false
			}
		}

		return true
	}
}

// Extract all registered method names from a gRPC server
// Returns fully qualified method names (e.g., "/package.Service/Method")
// If filter is nil, ALL services are included (no filtering)
func Extract(infoProvider ServiceInfoProvider, filter ServiceFilter) []string {
	if infoProvider == nil {
		return []string{}
	}

	methods := []string{}

	// Get all registered services
	serviceInfo := infoProvider.GetServiceInfo()

	for serviceName, info := range serviceInfo {
		// Apply filter if provided
		if filter != nil && !filter(serviceName) {
			continue
		}

		// Extract all methods for this service
		for _, method := range info.Methods {
			// Build fully qualified method name in gRPC format: /package.Service/Method
			fullyQualified := "/" + serviceName + "/" + method.Name
			methods = append(methods, fullyQualified)
		}
	}

	// Sort for deterministic ordering
	sort.Strings(methods)

	return methods
}
