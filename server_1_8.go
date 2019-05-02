// +build !go1.9

package httptesting

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"

	"github.com/dolab/httptesting/internal"
)

// NewServer returns an initialized *Testing along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServer(handler http.Handler, isTLS bool) *Testing {
	var (
		ts    *httptest.Server
		certs *x509.CertPool
	)
	if isTLS {
		cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
		if err != nil {
			panic(fmt.Sprintf("httptesting: NewTLSServer: %v", err))
		}

		ts = httptest.NewTLSServer(handler)

		certs = x509.NewCertPool()
		certs.AddCert(x509cert)
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

	return &Testing{
		server: ts,
		host:   urlobj.Host,
		certs:  certs,
		jar:    jar,
		isTLS:  isTLS,
	}
}

// NewServerWithTLS returns an initialized *Testing along with mocked server for testing
// NOTE: You MUST call client.Close() for cleanup after testing.
func NewServerWithTLS(handler http.Handler, cert tls.Certificate) *Testing {
	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(handler)
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	ts.StartTLS()

	certs := x509.NewCertPool()
	certs.AddCert(x509cert)

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err.Error())
	}

	urlobj, err := url.Parse(ts.URL)
	if err != nil {
		panic(err.Error())
	}

	return &Testing{
		server: ts,
		host:   urlobj.Host,
		certs:  certs,
		jar:    jar,
		isTLS:  true,
	}
}
