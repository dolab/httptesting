package httptesting

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

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
		w.Write([]byte(r.Header.Get("X-Mock-Client")))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)

	request := client.New(t)
	request.WithHeader("X-Mock-Client", "httptesting")

	request.Get("/request/client", nil)
	request.AssertOK()
	request.AssertContains("httptesting")
}

func TestRequestWithConcurrency(t *testing.T) {
	it := assert.New(t)

	method := "GET"
	uri := "/request/client"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("/request/client", r.RequestURI)

		time.Sleep(10 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Header.Get("X-Mock-Client")))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)

	var (
		wg sync.WaitGroup

		concurrency = 3
	)

	issuedAt := time.Now()

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			request := client.New(t)
			request.WithHeader("X-Mock-Client", "httptesting")

			request.Get("/request/client", nil)
			request.AssertOK()
			request.AssertContains("httptesting")
		}()
	}
	wg.Wait()

	delta := time.Since(issuedAt)
	it.True(delta < 20*time.Millisecond)
}
