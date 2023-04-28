package endpointresolver

import (
	"context"
)

// ResolveConf holds the configuration to be used by the resolver
type ResolveConf struct {
	// The endpoint to be used as a basis for the endpoint resolution process
	Endpoint string

	// Ports that should be included in the endpoint resolution
	Ports []int

	// The user agent string to use when attempting endpoint resolution
	UserAgent string

	// Any custom headers that might be needed so that endpoint-resolver's requests come across
	CustomHeaders map[string]string
}

// Resolver provides an interface which facilitates the process to resolve an endpoint.
type Resolver interface {

	// Resolve consumes the resolving config provided and will return back a list of URLs consisting of hostname with
	// open ports found, or an error.
	Resolve(ctx context.Context, conf ResolveConf) (urls []string, err error)
}

// Checker provides methods executing the actual endpoint resolution checks performed by the Resolver
type Checker interface {
	// ExternalDNS initiates the DNS resolution process by using an external DNS provider as the resolver.
	ExternalDNS(ctx context.Context, hostname string, externalDNS []string) error

	// NativeDNS initializes the DNS resolution process by using the cluster-internal DNS resolvers.
	NativeDNS(ctx context.Context, hostname string) (ips []string, err error)

	// Ports consumes a list of IPs discovered as well as ports and returns back a list of open ports accross them.
	// It does that by looping (max 3 attempts) through the IPs discovered and consequently the ports provided, and executes
	// a TCP-dial on each combination.
	Ports(ctx context.Context, ips []string, ports []int) (openPorts []int, err error)

	// HTTP sends an HTTP request to the open ports found on your hostname and returns back a list of URLs. It might be
	// that no URLs are to be returned - in that case an out-of-scope error might be returned, or a HTTP timeout warning. In
	// the case the user agent provided resulted in the request being blocked, then a relevant error is returned.
	HTTP(ctx context.Context, userAgent, hostname string, customHeaders map[string]string, openPorts []int) ([]string, error)
}
