package applicationscanning

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	endpointresolver "github.com/detectify/endpoint-resolver"
	"github.com/detectify/n5/domain"
)

var (
	defaultPorts = []int{80, 443}
)

const (
	maxRedirects             = 3
	dnsTimeout               = time.Second * 6
	dnsRetriesMaxElapsedTime = time.Minute * 2
	portCheckTimeout         = time.Second * 5
	portRetries              = 3
	httpTimeout              = time.Second * 30
	httpTimeoutLimit         = time.Second * 4
	mozillaUserAgent         = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
	schemeHTTP               = "http"
	schemeHTTPS              = "https"
)

func sendRequest(ctx context.Context, requestURL, userAgent string, customHeaders map[string]string) (*url.URL, error) {
	r, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	if len(customHeaders) > 0 {
		for k, v := range customHeaders {
			r.Header.Add(k, v)
		}
	}
	r.Header.Add("User-Agent", userAgent)

	r = r.WithContext(ctx)

	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateFreelyAsClient,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectedReqHost := via[len(via)-1].URL.Hostname()
			originalReqHost := r.URL.Hostname()

			// break after redirecting out of scope
			if !domain.Contains(originalReqHost, redirectedReqHost) {
				return http.ErrUseLastResponse
			}
			if len(via) > maxRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Timeout: httpTimeout,
	}
	response, err := c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed on request: %w", err)
	}
	return response.Request.URL, nil
}

func anyWithinTimeLimit(urls map[url.URL]time.Duration) bool {
	for _, duration := range urls {
		if duration < httpTimeoutLimit {
			return true
		}
	}
	return false
}

func anyWithinScope(urls map[url.URL]time.Duration, hostname string, openPorts []int) bool {
	for u := range urls {
		h := u.Hostname()
		if (h == hostname || domain.Contains(hostname, h)) && portInScope(u.Scheme, u.Port(), openPorts) {
			return true
		}
	}
	return false
}

func portInScope(scheme, port string, ports []int) bool {
	var portInt int
	if len(port) > 0 {
		portInt, _ = strconv.Atoi(port)
	} else if scheme == schemeHTTP {
		portInt = 80
	} else if scheme == schemeHTTPS {
		portInt = 443
	}
	for _, p := range ports {
		if p == portInt {
			return true
		}
	}
	return false
}

func convertURLs(urls map[url.URL]time.Duration) []string {
	var result []string
	for u := range urls {
		result = append(result, u.String())
	}
	return result
}

func createURLs(hostname string, port int) []string {
	switch port {
	// for default ports we only send specific schemes, and no need to include port in URL
	case 80:
		return []string{fmt.Sprintf("%s://%s/", schemeHTTP, hostname)}
	case 443:
		return []string{fmt.Sprintf("%s://%s/", schemeHTTPS, hostname)}
	default:
		// for any other ports, try both schemes, and explicitly define the port
		return []string{
			fmt.Sprintf("%s://%s:%d/", schemeHTTP, hostname, port),
			fmt.Sprintf("%s://%s:%d/", schemeHTTPS, hostname, port),
		}
	}
}

func fetchPorts(endpointPort string, ports []int) ([]int, error) {
	if len(endpointPort) > 0 {
		port, err := strconv.Atoi(endpointPort)
		if err != nil {
			return nil, endpointresolver.ErrInvalidEndpointPort
		}

		return []int{port}, nil
	}

	if len(ports) > 0 {
		return ports, nil
	}

	return defaultPorts, nil
}
