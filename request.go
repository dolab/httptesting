package httptesting

import "net/http"

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
