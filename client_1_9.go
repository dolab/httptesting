// +build go1.9

package httptesting

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
)

// NewServer returns an initialized *Client along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServer(handler http.Handler, isTLS bool) *Client {
	var ts *httptest.Server
	if isTLS {
		ts = httptest.NewTLSServer(handler)
	} else {
		ts = httptest.NewServer(handler)
	}

	urlobj, _ := url.Parse(ts.URL)

	client := ts.Client()
	client.Jar, _ = cookiejar.New(nil)

	return &Client{
		ts:     ts,
		client: client,
		host:   urlobj.Host,
		https:  isTLS,
	}
}

// NewServerWithTLS returns an initialized *Client along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServerWithTLS(handler http.Handler, cert tls.Certificate) *Client {
	ts := httptest.NewUnstartedServer(handler)
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ts.StartTLS()

	urlobj, _ := url.Parse(ts.URL)

	client := ts.Client()
	client.Jar, _ = cookiejar.New(nil)

	return &Client{
		ts:     ts,
		client: client,
		host:   urlobj.Host,
		https:  true,
	}
}
