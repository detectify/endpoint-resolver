package applicationscanning

import (
	"context"
	"testing"

	endpointresolver "github.com/detectify/endpoint-resolver"
	"github.com/stretchr/testify/require"
)

func TestResolve_ResolvingDomain(t *testing.T) {
	urls, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "detectify.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(urls))
}

func TestResolve_ResolvingDomainWithOpenPort(t *testing.T) {
	urls, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "detectify.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
		Ports:     []int{443},
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(urls))
}

func TestResolve_ResolvingDomainWithClosedPort(t *testing.T) {
	_, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "detectify.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
		Ports:     []int{8080},
	})
	require.Error(t, err)
	require.Equal(t, endpointresolver.ErrNoOpenPort, err)
}

func TestResolve_NonResolvingDomainName(t *testing.T) {
	_, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "nonexisting.domain",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
	})
	require.Error(t, err)
	require.Equal(t, endpointresolver.ErrThirdPartyDNSResolutionFailure, err)
}

func TestResolve_DomainRedirectingToWWW(t *testing.T) {
	// this domain is not considered to be redirecting because it redirects to the www-version of it
	urls, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "koslib.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(urls))
}

func TestResolve_RedirectingOutsideScope(t *testing.T) {
	urls, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "302.koslib.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
	})
	require.Equal(t, endpointresolver.WarnRedirectedOutOfScope, err)
	require.Equal(t, 1, len(urls))
}

func TestResolve_RedirectingOutsideScopeDueToOtherSchemeNotInScope(t *testing.T) {
	urls, err := NewResolver([]string{"8.8.8.8:53"}).Resolve(context.TODO(), endpointresolver.ResolveConf{
		Endpoint:  "detectify.com",
		UserAgent: "Mozilla/5.0 (compatible; Detectify)",
		Ports:     []int{80},
	})
	require.Equal(t, endpointresolver.WarnRedirectedOutOfScope, err)
	require.Equal(t, 1, len(urls))
}
