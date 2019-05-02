package httptesting

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"testing"

	"github.com/dolab/httptesting/internal"
	"github.com/golib/assert"
)

func Test_NewServer(t *testing.T) {
	it := assert.New(t)

	method := "GET"
	uri := "/server"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TLS"))
	})

	ts := NewServer(server, true)
	defer ts.Close()

	it.NotNil(ts.server)
	it.NotEmpty(ts.host)
	it.True(ts.isTLS)

	// it should work with internal client
	request := ts.New(t)
	request.Get("/server/tls", nil)
	request.AssertOK()
	request.AssertContains("TLS")

	// it should work with custom TLS client
	cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
	if it.Nil(err) {
		x509cert, err := x509.ParseCertificate(cert.Certificate[0])
		if it.Nil(err) {
			client := NewWithTLS(ts.Url(""), x509cert)

			request = client.New(t)
			request.Get("/server/tls", nil)
			request.AssertOK()
			request.AssertContains("TLS")
		}
	}
}

func Test_NewServerWithTLS(t *testing.T) {
	it := assert.New(t)

	method := "GET"
	uri := "/server/tls"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TLS"))
	})

	cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
	if it.Nil(err) {
		ts := NewServerWithTLS(server, cert)
		defer ts.Close()

		it.NotNil(ts.server)
		it.NotEmpty(ts.host)
		it.True(ts.isTLS)

		// it should work with internal client
		request := ts.New(t)
		request.Get("/server/tls", nil)
		request.AssertOK()
		request.AssertContains("TLS")

		// it should work with custom TLS client
		x509cert, err := x509.ParseCertificate(cert.Certificate[0])
		if it.Nil(err) {
			client := NewWithTLS(ts.Url(""), x509cert)

			request = client.New(t)
			request.Get("/server/tls", nil)
			request.AssertOK()
			request.AssertContains("TLS")
		}
	}
}
