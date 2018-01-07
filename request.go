package httptesting

import (
	"net/http"
	"testing"

	"golang.org/x/net/websocket"
)

// Request defines http client for human
type Request struct {
	*Client

	header http.Header
}

// NewRequest returns a new *Request with *Client
func NewRequest(client *Client) *Request {
	return &Request{
		Client: client,
		header: http.Header{},
	}
}

// WithHeader sets http.Request header by replace
func (r *Request) WithHeader(key, value string) *Request {
	r.header.Set(key, value)

	return r
}

// WithHttpHeader adds http.Request header
func (r *Request) WithHttpHeader(header http.Header) *Request {
	for key, values := range header {
		for _, value := range values {
			r.header.Add(key, value)
		}
	}

	return r
}

// NewRequest issues any request and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: You have to manage session / cookie data manually.
func (r *Request) NewRequest(request *http.Request) {
	r.Client.NewRequest(r.t, request)
}

// NewSessionRequest issues any request with session / cookie and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: Session data will be added to the request cookies for requested host.
func (r *Request) NewSessionRequest(request *http.Request) {
	r.Client.NewSessionRequest(r.t, request)
}

// NewFilterRequest issues any request with TransportFiler and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: It returns error without apply HTTP request when transport filter returned an error.
func (r *Request) NewFilterRequest(request *http.Request, filter TransportFilter) {
	r.Client.NewFilterRequest(r.t, request, filter)
}

// NewMultipartRequest issues a multipart request for the method & fields given and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
func (r *Request) NewMultipartRequest(method, path, filename string, file interface{}, fields ...map[string]string) {
	r.Client.NewMultipartRequest(r.t, method, path, filename, file, fields...)
}

// NewWebsocket creates a websocket connection to the given path and returns the connection
func (r *Request) NewWebsocket(t *testing.T, path string) *websocket.Conn {
	return r.Client.NewWebsocket(r.t, path)
}
