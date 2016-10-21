package httptesting

import (
	"net/http"
	"testing"

	"golang.org/x/net/websocket"
)

type RequestClient struct {
	*Client

	header http.Header
}

func NewRequestClient(client *Client) *RequestClient {
	return &RequestClient{
		Client: client,
		header: http.Header{},
	}
}

func (client *RequestClient) WithHeader(key, value string) *RequestClient {
	client.header.Add(key, value)

	return client
}

func (client *RequestClient) WithHttpHeader(header http.Header) *RequestClient {
	for key, values := range header {
		for _, value := range values {
			client.header.Add(key, value)
		}
	}

	return client
}

// NewRequest issues any request and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: You have to manage session / cookie data manually.
func (client *RequestClient) NewRequest(request *http.Request) {
	client.Client.NewRequest(client.t, request)
}

// NewSessionRequest issues any request with session / cookie and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: Session data will be added to the request cookies for you.
func (client *RequestClient) NewSessionRequest(request *http.Request) {
	client.Client.NewSessionRequest(client.t, request)
}

// NewFilterRequest issues any request with TransportFiler and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: It returns error without apply HTTP request when transport filter returned an error.
func (client *RequestClient) NewFilterRequest(request *http.Request, filter TransportFilter) {
	client.Client.NewFilterRequest(client.t, request, filter)
}

// NewMultipartRequest issues a multipart request for the method & fields given and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
func (client *RequestClient) NewMultipartRequest(method, path, filename string, file interface{}, fields ...map[string]string) {
	client.Client.NewMultipartRequest(client.t, method, path, filename, file, fields...)
}

// NewWebsocket creates a websocket connection to the given path and returns the connection
func (client *RequestClient) NewWebsocket(t *testing.T, path string) *websocket.Conn {
	return client.Client.NewWebsocket(client.t, path)
}
