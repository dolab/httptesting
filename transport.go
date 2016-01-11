package httptesting

import (
	"net"
	"net/http"
	"time"
)

// TransportFilter is a callback for http request injection.
type TransportFilter func(r *http.Request) error

type transport struct {
	callee TransportFilter
}

func newTransport(filter TransportFilter) *transport {
	return &transport{
		callee: filter,
	}
}

func (caller *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, address string) (net.Conn, error) {
			// invoke callee
			err := caller.callee(r)
			if err != nil {
				return nil, err
			}

			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}

			conn, err := dialer.Dial(network, address)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return transport.RoundTrip(r)
}
