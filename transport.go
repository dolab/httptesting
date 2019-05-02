package httptesting

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"time"
)

// RequestFilter is a callback for http request injection.
type RequestFilter func(r *http.Request) error

// FilterTransport defines a custom http.Transport with filters and certs.
type FilterTransport struct {
	filters []RequestFilter
	certs   []*x509.CertPool
}

func NewFilterTransport(filters []RequestFilter, certs ...*x509.CertPool) *FilterTransport {
	return &FilterTransport{
		filters: filters,
		certs:   certs,
	}
}

func (transport *FilterTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			// invoke filters
			for _, filter := range transport.filters {
				err := filter(r)
				if err != nil {
					return nil, err
				}
			}

			dialer := &net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}

			conn, err := dialer.Dial(network, address)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		ResponseHeaderTimeout: 3 * time.Second,
	}

	if len(transport.certs) > 0 {
		tr.TLSClientConfig = &tls.Config{
			RootCAs: transport.certs[0],
		}
		tr.TLSHandshakeTimeout = 5 * time.Second
	}

	return tr.RoundTrip(r)
}
