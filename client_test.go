package httptesting

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"

	"github.com/golib/assert"
)

type mockServer struct {
	method string
	path   string
	it     func(w http.ResponseWriter, r *http.Request)
}

func (mock *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mock.it(w, r)
}

var (
	newMockServer = func(method, path string, it func(http.ResponseWriter, *http.Request)) *mockServer {
		return &mockServer{
			method: method,
			path:   path,
			it:     it,
		}
	}
)

func Test_New(t *testing.T) {
	it := assert.New(t)

	host := "www.example.com"
	absurl := "https://" + host
	ws := "ws://" + host

	client := New(host, true)
	it.Nil(client.server)
	it.NotEmpty(client.host)
	it.True(client.isTLS)
	it.Equal(host, client.Host())
	it.Equal(absurl, client.Url(""))
	it.Equal(ws, client.WebsocketUrl(""))
}

func Test_NewWithRacy(t *testing.T) {
	method := "GET"
	uri := "/request/racy"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Query().Get("routine")))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)

	var (
		wg sync.WaitGroup

		routines = 3
	)
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(routine int) {
			defer wg.Done()

			params := url.Values{}
			params.Add("routine", strconv.Itoa(routine+1))

			request := client.New(t)
			request.Get(uri, params)
			request.AssertOK()
			request.AssertContains(strconv.Itoa(routine + 1))
		}(i)
	}

	wg.Wait()
}

func TestTesting_New(t *testing.T) {
	it := assert.New(t)

	host := "www.example.com"
	client := New(host, true)

	request := client.New(t)
	it.Equal(client, request.Client)
	it.NotNil(request.t)
	it.Nil(request.Response)
	it.Empty(request.ResponseBody)
	it.Nil(request.server)
	it.Equal(client.host, request.host)
	it.True(request.isTLS)
}
