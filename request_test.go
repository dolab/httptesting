package httptesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golib/assert"
)

func Test_RequestClient(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/request/client"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request/client", r.RequestURI)

		fmt.Printf(">>> %#v\n", r.Header)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Header.Get("X-Mock-Client")))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)

	rclient := client.New(t)
	rclient.WithHeader("X-Mock-Client", "httptesting")

	rclient.Get("/request/client", nil)
	rclient.AssertOK()
	rclient.AssertContains("httptesting")
}
