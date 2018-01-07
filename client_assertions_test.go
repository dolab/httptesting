package httptesting

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golib/assert"
)

func Test_Assertions(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/assertions"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/html", r.Header.Get("Content-Type"))
		assertion.Equal("/assertions?url-key=url-value", r.RequestURI)
		assertion.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user":{"name":"httptesting","age":3},"addresses":[{"name":"china"},{"name":"USA"}]}`))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.Get(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertExistHeader("Content-Type")
	client.AssertNotExistHeader("x-unknown-header")
	client.AssertNotEmpty()
	client.AssertContains(`{"name":"china"}`)
	client.AssertContainsJSON("user.name", "httptesting")
	client.AssertContainsJSON("addresses.1.name", "USA")
	client.AssertNotContainsJSON("addresses.0.post")
	client.AssertNotContainsJSON("addresses.3.name")
}
