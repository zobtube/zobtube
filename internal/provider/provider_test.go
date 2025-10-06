package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// --- helper to mock http.Transport ---
type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// --- helper to override http.DefaultTransport temporarily ---
func testHTTPTransports(
	mainServer *httptest.Server,
	secondaryServer *httptest.Server,
	mainServerOverrides []string,
	secondaryServerOverrides []string,
	fn func(),
) {
	old := http.DefaultTransport

	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		host := strings.TrimPrefix(mainServer.URL, "http://")
		if strings.Contains(req.URL.Host, host) {
			// catched by main server
			return old.RoundTrip(req)
		}
		for _, override := range mainServerOverrides {
			if strings.Contains(req.URL.Host, override) {
				req.URL.Scheme = "http"
				req.URL.Host = host
				return old.RoundTrip(req)
			}
		}

		if secondaryServer != nil {
			host = strings.TrimPrefix(secondaryServer.URL, "http://")
			if strings.Contains(req.URL.Host, host) {
				// catched by secondary server
				return old.RoundTrip(req)
			}
			for _, override := range secondaryServerOverrides {
				if strings.Contains(req.URL.Host, override) {
					req.URL.Scheme = "http"
					req.URL.Host = host
					return old.RoundTrip(req)
				}
			}
		}

		// uncatched
		return nil, fmt.Errorf("uncatched request towards: %s", req.URL.Host)
	})

	defer func() { http.DefaultTransport = old }()
	fn()
}
