package endpointresolver

import "errors"

var (
	// ErrInvalidEndpoint is returned when the endpoint is neither an IP nor a domain
	ErrInvalidEndpoint = errors.New("invalid endpoint")

	// ErrThirdPartyDNSResolutionFailure is returned when an error is returned when using the external DNS resolvers
	ErrThirdPartyDNSResolutionFailure = errors.New("third-party DNS resolution failure")

	// ErrNativeDNSResolutionFailure is returned when an error is returned when using the internal DNS resolvers
	ErrNativeDNSResolutionFailure = errors.New("native DNS resolution failure")

	// ErrNoIPForEndpoint is returned when an endpoint did not resolve to at least one IP address
	ErrNoIPForEndpoint = errors.New("no IP for endpoint")

	// ErrInvalidEndpointPort is returned when a port provided was malformed and could not be used
	ErrInvalidEndpointPort = errors.New("invalid endpoint port")

	// ErrNoOpenPort is returned when no open port was found on a given endpoint
	ErrNoOpenPort = errors.New("no open port")

	// ErrNoHTTPConnection is returned when the resolver's process is unable to establish an HTTP connection
	ErrNoHTTPConnection = errors.New("no HTTP connection")

	// ErrBlockedByUserAgent is returned when the HTTP request was blocked due to the user agent string
	ErrBlockedByUserAgent = errors.New("blocked by user-agent")

	// ErrIPV6Unsupported is returned when the endpoint resolver needs to process an IPv6 address, which is
	// currently unsupported
	ErrIPV6Unsupported = errors.New("IPv6 addresses are not supported")

	// WarnHTTPTimeout error is returned when the HTTP response occurred after a time threshold indicating that
	// responses are slower than anticipated
	WarnHTTPTimeout = errors.New("warning: HTTP timeout") //nolint:revive

	// WarnRedirectedOutOfScope error is returned when the endpoint redirected out of scope, meaning another endpoint
	WarnRedirectedOutOfScope = errors.New("warning: redirection occurred outside scope") //nolint:revive
)
