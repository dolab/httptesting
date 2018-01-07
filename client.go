package httptesting

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

// Client defines request component of httptesting
type Client struct {
	Response     *http.Response
	ResponseBody []byte

	t      *testing.T
	client *http.Client
	host   string
	https  bool
}

// New returns an initialized Client ready for testing
func New(host string, isHTTPS bool) *Client {
	jar, _ := cookiejar.New(nil)

	// adjust host
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		urlobj, err := url.Parse(host)
		if err == nil {
			host = urlobj.Host
		}
	}

	return &Client{
		client: &http.Client{
			Jar: jar,
		},
		host:  host,
		https: isHTTPS,
	}
}

// Host returns the host and port of the server, e.g. "127.0.0.1:9090"
func (c *Client) Host() string {
	if c.host[0] == ':' {
		return "127.0.0.1" + c.host
	}

	return c.host
}

// Url returns the abs http/https URL of the resource, e.g. "http://127.0.0.1:9090/status".
// The scheme is set to https if http.ssl is set to true in the configuration.
func (c *Client) Url(path string) string {
	if c.https {
		return "https://" + c.Host() + path
	}

	return "http://" + c.Host() + path
}

// WebsocketUrl returns the abs websocket URL of the resource, e.g. "ws://127.0.0.1:9090/status"
func (c *Client) WebsocketUrl(path string) string {
	return "ws://" + c.Host() + path
}

// Cookies returns cookies related to the host
func (c *Client) Cookies() []*http.Cookie {
	urlobj, _ := url.Parse(c.Url("/"))

	return c.client.Jar.Cookies(urlobj)
}

// SetCookies sets cookies for the host
func (c *Client) SetCookies(cookies []*http.Cookie) {
	urlobj, _ := url.Parse(c.Url("/"))

	c.client.Jar.SetCookies(urlobj, cookies)
}

// New returns a Request which has more customlization!
func (c *Client) New(t *testing.T) *Request {
	// copy avoiding data race issue
	nc := *c
	nc.t = t

	return NewRequest(&nc)
}

// NewRequest issues any request and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: You have to manage session / cookie data manually.
func (c *Client) NewRequest(t *testing.T, request *http.Request) {
	c.t = t

	var err error

	c.Response, err = c.client.Do(request)
	if err != nil {
		t.Fatalf("[REQUEST] %s %s: %#v\n", request.Method, request.URL.RequestURI(), err)
	}
	defer c.Response.Body.Close()

	// Read response body if not empty
	c.ResponseBody = []byte{}

	switch c.Response.StatusCode {
	case http.StatusNoContent:
		// ignore

	default:
		c.ResponseBody, err = ioutil.ReadAll(c.Response.Body)
		if err != nil {
			if err != io.EOF {
				t.Fatalf("[RESPONSE] %s %s: %#v\n", request.Method, request.URL.RequestURI(), err)
			}

			// unexpected EOF with content-length
			if c.Response.ContentLength > 0 && int64(len(c.ResponseBody)) != c.Response.ContentLength {
				t.Fatalf("[RESPONSE] %s %s: %#v\n", request.Method, request.URL.RequestURI(), err)
			}

			t.Logf("[RESPONSE] %s %s: Unexptected response body with io.EOF error.", request.Method, request.URL.RequestURI())
		}
	}
}

// NewSessionRequest issues any request with session/cookie and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: Session data will be added to the request cookies for you.
func (c *Client) NewSessionRequest(t *testing.T, request *http.Request) {
	for _, cookie := range c.client.Jar.Cookies(request.URL) {
		request.AddCookie(cookie)
	}

	c.NewRequest(t, request)
}

// NewFilterRequest issues any request with TransportFilter and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: It returns error without apply HTTP request when transport filter returned an error.
func (c *Client) NewFilterRequest(t *testing.T, request *http.Request, filter TransportFilter) {
	c.t = t

	var err error

	client := &http.Client{
		Transport: newTransport(filter),
	}

	c.Response, err = client.Do(request)
	if err != nil {
		t.Fatalf("[FILTERED REQUEST] %s %s: %#v\n", request.Method, request.URL.RequestURI(), err)
	}

	// Read response body
	c.ResponseBody, err = ioutil.ReadAll(c.Response.Body)
	if err != nil {
		t.Fatalf("[FILTERED RESPONSE] %s %s: %#v\n", request.Method, request.URL.RequestURI(), err)
	}
	c.Response.Body.Close()
}

// NewMultipartRequest issues a multipart request for the method & fields given and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
func (c *Client) NewMultipartRequest(t *testing.T, method, path, filename string, file interface{}, fields ...map[string]string) {
	c.t = t

	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)

	fw, ferr := mw.CreateFormFile("filename", filename)
	if ferr != nil {
		t.Fatalf("[MULTIPART REQUEST] %s %s: %#v\n", method, path, ferr)
	}

	// apply file
	var (
		reader io.Reader
		err    error
	)
	switch file.(type) {
	case io.Reader:
		reader, _ = file.(io.Reader)

	case *os.File:
		reader, _ = file.(*os.File)

	case string:
		reader, err = os.Open(file.(string))
		if err != nil {
			t.Fatalf("[MULTIPART REQUEST] os.Open(%v): %#v\n", file, err)
		}

	default:
		t.Fatalf("[MULTIPART REQUEST] unsupported file type: %T\n", file)
	}

	if _, err := io.Copy(fw, reader); err != nil {
		t.Fatalf("[MULTIPART REQUEST] io.Copy(%T, %T): %#v\n", fw, file, err)
	}

	// apply fields
	if len(fields) > 0 {
		for key, value := range fields[0] {
			mw.WriteField(key, value)
		}
	}

	// adds the terminating boundary
	mw.Close()

	request, err := http.NewRequest(method, c.Url(path), &buf)
	if err != nil {
		t.Fatalf("[MULTIPART REQUEST] %s %s: %#v\n", method, path, err)
	}
	request.Header.Set("Content-Type", mw.FormDataContentType())

	c.NewRequest(t, request)
}

// NewWebsocket creates a websocket connection to the given path and returns the connection
func (c *Client) NewWebsocket(t *testing.T, path string) *websocket.Conn {
	origin := c.WebsocketUrl("/")
	target := c.WebsocketUrl(path)

	ws, err := websocket.Dial(target, "", origin)
	if err != nil {
		t.Fatalf("[WS REQUEST] connect %s: %#v\n", path, err)
	}

	return ws
}
