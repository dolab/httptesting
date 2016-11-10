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

type Client struct {
	Client       *http.Client
	Response     *http.Response
	ResponseBody []byte

	t     *testing.T
	host  string
	https bool
}

// NewClient returns an initialized Client ready for using
func New(host string, isHttps bool) *Client {
	jar, _ := cookiejar.New(nil)

	// adjust host
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		u, err := url.Parse(host)
		if err == nil {
			host = u.Host
		}
	}

	return &Client{
		Client: &http.Client{Jar: jar},
		host:   host,
		https:  isHttps,
	}
}

// Host returns the host and port of the server, e.g. "127.0.0.1:9090"
func (test *Client) Host() string {
	if test.host[0] == ':' {
		return "127.0.0.1" + test.host
	}

	return test.host
}

// Url returns the abs http/https URL of the resource, e.g. "http://127.0.0.1:9090/status".
// The scheme is set to https if http.ssl is set to true in the configuration.
func (test *Client) Url(path string) string {
	if test.https {
		return "https://" + test.Host() + path
	}

	return "http://" + test.Host() + path
}

// WebsocketUrl returns the abs websocket URL of the resource, e.g. "ws://127.0.0.1:9090/status"
func (test *Client) WebsocketUrl(path string) string {
	return "ws://" + test.Host() + path
}

// Cookies returns cookies related with the host
func (test *Client) Cookies() []*http.Cookie {
	u, _ := url.Parse(test.Url("/"))

	return test.Client.Jar.Cookies(u)
}

// SetCookie sets cookies with the host
func (test *Client) SetCookies(cookies []*http.Cookie) {
	u, _ := url.Parse(test.Url("/"))

	test.Client.Jar.SetCookies(u, cookies)
}

// New returns a RequestClient which has more customlization!
func (test *Client) New(t *testing.T) *RequestClient {
	client := NewRequestClient(test)
	client.t = t

	return client
}

// NewRequest issues any request and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: You have to manage session / cookie data manually.
func (test *Client) NewRequest(t *testing.T, request *http.Request) {
	test.t = t

	var err error

	test.Response, err = test.Client.Do(request)
	if err != nil {
		t.Fatalf("[REQUEST] %s %s: %#v\n", request.Method, request.URL.Path, err.Error())
	}

	// Read response body if not empty
	if test.Response.ContentLength > 0 {
		test.ResponseBody, err = ioutil.ReadAll(test.Response.Body)
		if err != nil {
			t.Fatalf("[RESPONSE] %s %s: %#v\n", request.Method, request.URL.Path, err)
		}
	}
	test.Response.Body.Close()
}

// NewSessionRequest issues any request with session / cookie and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: Session data will be added to the request cookies for you.
func (test *Client) NewSessionRequest(t *testing.T, request *http.Request) {
	for _, cookie := range test.Client.Jar.Cookies(request.URL) {
		request.AddCookie(cookie)
	}

	test.NewRequest(t, request)
}

// NewFilterRequest issues any request with TransportFiler and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: It returns error without apply HTTP request when transport filter returned an error.
func (test *Client) NewFilterRequest(t *testing.T, request *http.Request, filter TransportFilter) {
	test.t = t

	var err error

	client := &http.Client{
		Transport: newTransport(filter),
	}

	test.Response, err = client.Do(request)
	if err != nil {
		t.Fatalf("[REQUEST] %s %s: %#v\n", request.Method, request.URL.Path, err.Error())
	}

	// Read response body
	test.ResponseBody, err = ioutil.ReadAll(test.Response.Body)
	if err != nil {
		t.Fatalf("[RESPONSE] %s %s: %#v\n", request.Method, request.URL.Path, err)
	}
	test.Response.Body.Close()
}

// NewMultipartRequest issues a multipart request for the method & fields given and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
func (test *Client) NewMultipartRequest(t *testing.T, method, path, filename string, file interface{}, fields ...map[string]string) {
	test.t = t

	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)

	fw, ferr := mw.CreateFormFile("filename", filename)
	if ferr != nil {
		t.Fatalf("%s %s: %#v\n", method, path, ferr)
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
		filepath, _ := file.(string)

		reader, err = os.Open(filepath)
		if err != nil {
			t.Fatalf("%s %s: %#v\n", method, path, err)
		}

	}

	if _, err := io.Copy(fw, reader); err != nil {
		t.Fatalf("%s %s: %#v\n", method, path, err)
	}

	// apply fields
	if len(fields) > 0 {
		for key, value := range fields[0] {
			mw.WriteField(key, value)
		}
	}

	// adds the terminating boundary
	mw.Close()

	request, err := http.NewRequest(method, test.Url(path), &buf)
	if err != nil {
		t.Fatalf("%s %s: %#v\n", method, path, err)
	}
	request.Header.Set("Content-Type", mw.FormDataContentType())

	test.NewRequest(t, request)
}

// NewWebsocket creates a websocket connection to the given path and returns the connection
func (test *Client) NewWebsocket(t *testing.T, path string) *websocket.Conn {
	origin := test.WebsocketUrl("/")
	target := test.WebsocketUrl(path)

	ws, err := websocket.Dial(target, "", origin)
	if err != nil {
		t.Fatalf("WS %s: %#v\n", path, err)
	}

	return ws
}
