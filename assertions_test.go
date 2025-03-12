package httptesting

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golib/assert"
)

func TestRequest_Assertions(t *testing.T) {
	it := assert.New(t)
	method := "GET"
	uri := "/assertions"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/html", r.Header.Get("Content-Type"))
		it.Equal("/assertions?url-key=url-value", r.RequestURI)
		it.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user":{"name":"httptesting","age":3},"addresses":[{"name":"china"},{"name":"USA"}]}`))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.Get(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertExistHeader("Content-Type")
	request.AssertNotExistHeader("x-unknown-header")
	request.AssertNotEmpty()
	request.AssertContains(`{"name":"china"}`)
	request.AssertContainsJSON("user.name", "httptesting")
	request.AssertContainsJSON("addresses.1.name", "USA")
	request.AssertNotContainsJSON("addresses.0.post")
	request.AssertNotContainsJSON("addresses.3.name")
}
