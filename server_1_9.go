//go:build go1.9

package httptesting

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
)

// NewServer returns an initialized *Client along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServer(handler http.Handler, isTLS bool) *Client {
	var (
		ts    *httptest.Server
		certs *x509.CertPool
	)
	if isTLS {
		ts = httptest.NewTLSServer(handler)

		if transport, ok := ts.Client().Transport.(*http.Transport); ok {
			certs = transport.TLSClientConfig.RootCAs
		}
	} else {
		ts = httptest.NewServer(handler)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err.Error())
	}

	urlobj, err := url.Parse(ts.URL)
	if err != nil {
		panic(err.Error())
	}

	return &Client{
		server: ts,
		host:   urlobj.Host,
		certs:  certs,
		jar:    jar,
		isTLS:  isTLS,
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err.Error())
	}

	urlobj, err := url.Parse(ts.URL)
	if err != nil {
		panic(err.Error())
	}

	var certs *x509.CertPool
	if transport, ok := ts.Client().Transport.(*http.Transport); ok {
		certs = transport.TLSClientConfig.RootCAs
	}

	return &Client{
		server: ts,
		host:   urlobj.Host,
		certs:  certs,
		jar:    jar,
		isTLS:  true,
	}
}
