package applicationscanning

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	endpointresolver "github.com/detectify/endpoint-resolver"
	"github.com/detectify/n5/ip"
	"github.com/miekg/dns"
	"net"
	"net/url"
	"time"
)

type Checker struct {
}

// ExternalDNS initiates the DNS resolution process by using an external DNS provider as the resolver.
func (c Checker) ExternalDNS(ctx context.Context, hostname string, externalDNS []string) error {
	client := dns.Client{}
	client.Timeout = dnsTimeout

	req := new(dns.Msg)
	req.Id = dns.Id()
	req.RecursionDesired = true
	req.Question = make([]dns.Question, 1)
	req.Question[0] = dns.Question{Name: dns.Fqdn(hostname), Qtype: dns.TypeA, Qclass: dns.ClassINET}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = dnsRetriesMaxElapsedTime

	err := backoff.Retry(func() error {
		for _, r := range externalDNS {
			res, _, err := client.ExchangeContext(ctx, req, r)

			switch {
			case err == context.Canceled:
				return err
			case err != nil:
				// We got no results, try with next resolver
				continue
			case res == nil || res.Answer == nil:
				// We got an error or no valid response, try with next resolver
				continue
			case res.Rcode == dns.RcodeRefused || res.Rcode == dns.RcodeServerFailure:
				// We got results, but they were bad, try with next resolver
				continue
			}
			return nil
		}
		return endpointresolver.ErrThirdPartyDNSResolutionFailure
	}, b)
	if err != nil {
		return endpointresolver.ErrThirdPartyDNSResolutionFailure
	}

	return nil
}

// NativeDNS initializes the DNS resolution process by using the cluster-internal DNS resolvers.
func (c Checker) NativeDNS(ctx context.Context, hostname string) (ips []string, err error) {
	allIps, err := net.DefaultResolver.LookupHost(ctx, hostname)
	switch err {
	case nil:
		// since we can process only IPv4 addresses, we filter only for them
		for _, ipAdd := range allIps {
			if ip.IsIPv4(ipAdd) {
				ips = append(ips, ipAdd)
			}
		}
		return ips, nil
	case context.Canceled:
		return nil, err
	default:
		return nil, endpointresolver.ErrNativeDNSResolutionFailure
	}
}

// Ports consumes a list of IPs discovered as well as ports and returns back a list of open ports accross them.
// It does that by looping (max 3 attempts) through the IPs discovered and consequently the ports provided, and executes
// a TCP-dial on each combination.
func (c Checker) Ports(ctx context.Context, ips []string, ports []int) (openPorts []int, err error) {
	openPortMap := make(map[int]interface{}, 0)
	dialer := net.Dialer{
		Timeout: portCheckTimeout,
	}

	for i := 0; i < portRetries; i++ {
		for _, ipAddress := range ips {
			for _, port := range ports {
				if _, ok := openPortMap[port]; ok {
					// ports that were found open already do not need to be revisited
					continue
				}

				conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ipAddress, port))
				switch err {
				case nil:
					_ = conn.Close()
					openPortMap[port] = struct{}{}
				case context.Canceled:
					return nil, err
				default:
				}
			}
		}
	}

	if len(openPortMap) == 0 {
		return nil, endpointresolver.ErrNoOpenPort
	}

	for port := range openPortMap {
		openPorts = append(openPorts, port)
	}

	return openPorts, nil
}

// HTTP sends an HTTP request to the open ports found on your hostname and returns back a list of URLs. It might be
// that no URLs are to be returned - in that case an out-of-scope error might be returned, or a HTTP timeout warning. In
// the case the user agent provided resulted in the request being blocked, then a relevant error is returned.
func (c Checker) HTTP(ctx context.Context, userAgent, hostname string, customHeaders map[string]string, openPorts []int) ([]string, error) {
	urls := make(map[url.URL]time.Duration)

	for _, port := range openPorts {
		for _, candidateURL := range createURLs(hostname, port) {
			start := time.Now()
			responseURL, err := sendRequest(ctx, candidateURL, userAgent, customHeaders)
			switch err {
			case nil:
				duration := time.Since(start)
				urls[*responseURL] = duration
			case context.Canceled:
				return nil, err
			default:
			}
		}
	}
	if len(urls) > 0 {
		if !anyWithinScope(urls, hostname, openPorts) {
			return convertURLs(urls), endpointresolver.WarnRedirectedOutOfScope
		}
		if !anyWithinTimeLimit(urls) {
			return convertURLs(urls), endpointresolver.WarnHTTPTimeout
		}
		return convertURLs(urls), nil
	}

	for _, port := range openPorts {
		for _, requestURL := range createURLs(hostname, port) {
			_, err := sendRequest(ctx, requestURL, mozillaUserAgent, customHeaders)
			switch err {
			case nil:
				return nil, endpointresolver.ErrBlockedByUserAgent
			case context.Canceled:
				return nil, err
			default:
			}
		}
	}

	return nil, endpointresolver.ErrNoHTTPConnection
}
