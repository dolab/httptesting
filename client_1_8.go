// +build !go1.9

package httptesting

import (
	"log"
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

	return &Client{
		ts:     ts,
		client: client: &http.Client{
			Jar: jar,
		},,
		host:   urlobj.Host,
		https:  isTLS,
	}
}

// NewServerWithTLS returns an initialized *Client along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServerWithTLS(handler http.Handler, cert tls.Certificate) *Client {
	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(handler)
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ts.StartTLS()

	urlobj, _ := url.Parse(ts.URL)
	jar, _ = cookiejar.New(nil)

	certPool := x509.NewCertPool()
	certPool.AddCert(x509cert)

	return &Client{
		ts:     ts,
		client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			},
		},
		host:   urlobj.Host,
		https:  true,
	}
}
