package httptesting

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golib/assert"
)

func TestRequest(t *testing.T) {
	it := assert.New(t)
	method := "GET"
	uri := "/request/client"

	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("/request/client", r.RequestURI)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Header.Get("X-Mock-Testing")))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.WithHeader("X-Mock-Testing", "httptesting")

	request.Get("/request/client", nil)
	request.AssertOK()
	request.AssertContains("httptesting")
}
