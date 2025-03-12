package httptesting

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

// Client defines request component of httptesting.
//
// NOTE: Client is not safe for concurrency, please use client.New(t) after initialized.
type Client struct {
	mux    sync.RWMutex
	server *httptest.Server
	host   string
	certs  *x509.CertPool
	jar    *cookiejar.Jar
	isTLS  bool
}

// New returns an initialized *Client ready for testing
func New(host string, isTLS bool) *Client {
	// adjust host
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		urlobj, err := url.Parse(host)
		if err == nil {
			isTLS = strings.HasPrefix(host, "https://")

			host = urlobj.Host
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("httptesting: New: %v", err))
	}

	return &Client{
		host:  host,
		jar:   jar,
		isTLS: isTLS,
	}
}

// NewWithTLS returns an initialized *Client with custom certificate.
func NewWithTLS(host string, cert *x509.Certificate) *Client {
	// adjust host
	if strings.HasPrefix(host, "https://") {
		urlobj, err := url.Parse(host)
		if err == nil {
			host = urlobj.Host
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(fmt.Sprintf("httptesting: NewWithTLS: %v", err))
	}

	certs := x509.NewCertPool()
	certs.AddCert(cert)

	return &Client{
		host:  host,
		certs: certs,
		jar:   jar,
		isTLS: true,
	}
}

// Host returns the host and port of the server, e.g. "127.0.0.1:9090"
func (c *Client) Host() string {
	if len(c.host) == 0 {
		return ""
	}

	if c.host[0] == ':' {
		return "127.0.0.1" + c.host
	}

	return c.host
}

// Url returns the abs http/isTLS URL of the resource, e.g. "http://127.0.0.1:9090/status".
// The scheme is set to isTLS if http.ssl is set to true in the configuration.
func (c *Client) Url(urlpath string, params ...url.Values) string {
	if len(params) > 0 {
		if !strings.Contains(urlpath, "?") {
			urlpath += "?"
		}

		urlpath += params[0].Encode()
	}

	scheme := "http://"
	if c.isTLS {
		scheme = "https://"
	}

	return scheme + c.Host() + urlpath
}

// WebsocketUrl returns the abs websocket URL of the resource, e.g. "ws://127.0.0.1:9090/status"
func (c *Client) WebsocketUrl(urlpath string, params ...url.Values) string {
	if len(params) > 0 {
		if !strings.Contains(urlpath, "?") {
			urlpath += "?"
		}

		urlpath += params[0].Encode()
	}

	return "ws://" + c.Host() + urlpath
}

// Cookies returns jar related to the host
func (c *Client) Cookies() ([]*http.Cookie, error) {
	urlobj, err := url.Parse(c.Url("/"))
	if err != nil {
		return nil, err
	}

	return c.jar.Cookies(urlobj), nil
}

// SetCookies sets jar for the host
func (c *Client) SetCookies(cookies []*http.Cookie) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	urlobj, err := url.Parse(c.Url("/"))
	if err != nil {
		return err
	}

	c.jar.SetCookies(urlobj, cookies)
	return nil
}

// NewClient creates a http client with cookie and tls for the Client.
func (c *Client) NewClient(filters ...RequestFilter) *http.Client {
	client := &http.Client{
		Transport: NewFilterTransport(filters, c.certs),
		Jar:       c.jar,
	}

	return client
}

// NewWebsocket creates a websocket connection to the given path and returns the connection
func (c *Client) NewWebsocket(t *testing.T, path string) *websocket.Conn {
	origin := c.WebsocketUrl("/")
	target := c.WebsocketUrl(path)

	ws, err := websocket.Dial(target, "", origin)
	if err != nil {
		t.Fatalf("httptesting: NewWebscoket: connect %s with %v\n", path, err)
	}

	return ws
}

// NewRequest returns a *Request which has more customization!
func (c *Client) NewRequest(t *testing.T) *Request {
	return NewRequest(t, c)
}

// New is alias of NewRequest for shortcut.
func (c *Client) New(t *testing.T) *Request {
	return c.NewRequest(t)
}

// Close tries to
//
//   - close *httptest.Server created by NewServer or NewServerWithTLS
func (c *Client) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.server != nil {
		c.server.Close()
		c.server = nil
	}
}
