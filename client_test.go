package httptesting

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/golib/assert"
)

type mockServer struct {
	method    string
	path      string
	assertion func(w http.ResponseWriter, r *http.Request)
}

func (mock *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mock.assertion(w, r)
}

var (
	newMockServer = func(method, path string, assertion func(http.ResponseWriter, *http.Request)) *mockServer {
		return &mockServer{
			method:    method,
			path:      path,
			assertion: assertion,
		}
	}
)

func Test_New(t *testing.T) {
	assertion := assert.New(t)
	host := "www.example.com"
	url := "https://" + host
	ws := "ws://" + host

	client := New(host, true)
	assertion.Equal(host, client.Host())
	assertion.Equal(url, client.Url(""))
	assertion.Equal(ws, client.WebsocketUrl(""))
}

func Test_NewWithRacy(t *testing.T) {
	method := "GET"
	uri := "/request/racy"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
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

			// t.Run("routine@#"+strconv.Itoa(routine), func(subt *testing.T) {
			request := client.New(t)
			request.Get(uri, nil)
			request.AssertOK()
			// })
		}(i)
	}

	wg.Wait()
}

func Test_NewRequest(t *testing.T) {
	assertion := assert.New(t)

	method := "GET"
	uri := "/request"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request", r.RequestURI)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request, _ := http.NewRequest(method, ts.URL+"/request", nil)

	client := New(ts.URL, false)
	client.NewRequest(t, request)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertEmpty()
}

func Test_NewSessionRequest(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/request/session"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request/session", r.RequestURI)

		cookie, err := r.Cookie("client")
		assertion.Nil(err)
		assertion.Equal("httptesting", cookie.Value)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request, _ := http.NewRequest(method, ts.URL+"/request/session", nil)

	client := New(ts.URL, false)
	client.SetCookies([]*http.Cookie{
		{
			Name:   "client",
			Value:  "httptesting",
			MaxAge: 100,
		},
	})
	client.NewSessionRequest(t, request)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertEmpty()
}

func Test_NewFilterRequest(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/request/filter"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request/filter", r.RequestURI)
		assertion.Equal("[FILTERED]", r.Header.Get("X-Filtered-Header"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request, _ := http.NewRequest(method, ts.URL+"/request/filter", nil)

	client := New(ts.URL, false)
	client.NewFilterRequest(t, request, func(r *http.Request) error {
		r.Header.Add("X-Filtered-Header", "[FILTERED]")

		return nil
	})
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertEmpty()
}

func Test_MultipartRequest(t *testing.T) {
	assertion := assert.New(t)
	method := "put"
	uri := "/request/multipart"
	filename := "gopher.png"
	file := filepath.Clean("./fixtures/gopher.png")
	params := map[string]string{"form-key": "form-value"}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request/multipart", r.RequestURI)
		assertion.Equal("form-value", r.FormValue("form-key"))

		freader, fheader, ferr := r.FormFile("filename")
		assertion.Nil(ferr)
		assertion.Equal(filename, fheader.Filename)

		fb, _ := ioutil.ReadAll(freader)
		b, _ := ioutil.ReadFile(file)
		assertion.Equal(b, fb)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.NewMultipartRequest(t, method, uri, "gopher.png", file, params)
	client.AssertStatus(http.StatusNoContent)
	client.AssertHeader("x-request-method", method)
	client.AssertEmpty()
}
