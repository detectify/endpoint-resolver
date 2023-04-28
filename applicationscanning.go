package endpointresolver

import (
	"github.com/detectify/endpoint-resolver/applicationscanning"
	"github.com/detectify/endpoint-resolver/opentelemetry"
)

// NewApplicationScanningResolverWithTracing generates and returns a complete Resolver instance with
// OpenTelemetry tracing.
func NewApplicationScanningResolverWithTracing(externalDNS []string) Resolver {
	checker := applicationscanning.Checker{}
	checkerWithTracing := opentelemetry.NewCheckerWithTracing(checker, "checker")
	resolver := applicationscanning.NewResolverWithCheckers(externalDNS, checkerWithTracing)
	return opentelemetry.NewResolverWithTracing(resolver, "resolver")
}
