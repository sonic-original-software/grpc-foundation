//revive:disable:package-comments
package client

import "google.golang.org/grpc"

// Dependency holds information about a client's connection to a downstream gRPC service
type Dependency interface {
	grpc.ClientConnInterface
	Target() string
}

// Dependencies is a bag of Dependency, keyed by a name
type Dependencies map[string]Dependency
