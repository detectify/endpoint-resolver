package applicationscanning

import (
	"context"
	endpointresolver "github.com/detectify/endpoint-resolver"
	"github.com/detectify/n5/domain"
	"github.com/detectify/n5/ip"
	"strings"
)

// Resolver implements the Resolver interface and supports external DNS resolvers.
type Resolver struct {
	// externalDNS is a slice of IPs for external DNS resolvers that we can use.
	externalDNS []string
	checker     endpointresolver.Checker
}

// NewResolver generates and returns a Resolver pointer instance including any external DNS resolver configuration.
func NewResolver(externalDNS []string) *Resolver {
	return &Resolver{externalDNS: externalDNS, checker: Checker{}}
}

// NewResolverWithCheckers generates and returns a Resolver pointer instance including any external DNS resolver
// configuration as well as injects a Checker interface implementation.
func NewResolverWithCheckers(externalDNS []string, checker endpointresolver.Checker) *Resolver {
	return &Resolver{externalDNS: externalDNS, checker: checker}
}

// Resolve does a full resolution check by consequently executing open ports, DNS and HTTP checks. Returns back a list
// of valid URLs, or an error.
func (c *Resolver) Resolve(ctx context.Context, conf endpointresolver.ResolveConf) (urls []string, err error) {
	endpointParts := strings.Split(conf.Endpoint, ":")
	hostname := endpointParts[0]
	var portStr string
	if len(endpointParts) >= 2 {
		portStr = endpointParts[1]
	}

	ports, err := fetchPorts(portStr, conf.Ports)
	if err != nil {
		return nil, err
	}

	isDomain := domain.IsDomainName(hostname)

	isIP := ip.IsIP(hostname)

	var ips []string

	// If it is a domain (not an IP), we'll do some DNS checks
	if isDomain {
		err = c.checker.ExternalDNS(ctx, hostname, c.externalDNS)
		if err != nil {
			return nil, err
		}

		ips, err = c.checker.NativeDNS(ctx, hostname)
		if err != nil {
			return nil, err
		}

		if len(ips) == 0 {
			return nil, endpointresolver.ErrNoIPForEndpoint
		}
	}

	if isIP {
		ips = []string{hostname}

		// as application-scanning currently supports only ipv4 addresses, this is why we won't consider ipv6 addresses
		// as valid input.
		if !ip.IsIPv4(hostname) {
			return nil, endpointresolver.ErrIPV6Unsupported
		}
	}

	// it's not a domain, and it's not an IP, so erroring
	if !isDomain && !isIP {
		return nil, endpointresolver.ErrInvalidEndpoint
	}

	openPorts, err := c.checker.Ports(ctx, ips, ports)
	if err != nil {
		return nil, err
	}

	return c.checker.HTTP(ctx, conf.UserAgent, hostname, conf.CustomHeaders, openPorts)
}
