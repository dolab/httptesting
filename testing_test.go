package httptesting

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
	host := "www.example.com"
	url := "https://" + host
	ws := "ws://" + host
	assertion := assert.New(t)

	client := New(host, true)
	assertion.Equal(host, client.Host())
	assertion.Equal(url, client.Url(""))
	assertion.Equal(ws, client.WebsocketUrl(""))
}

func Test_MultipartRequest(t *testing.T) {
	assertion := assert.New(t)
	method := "put"
	uri := "/put/multipart"
	filename := "gopher.png"
	file := filepath.Clean("./fixtures/gopher.png")
	params := map[string]string{"form_key": "form_value"}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/put/multipart", r.RequestURI)
		assertion.Equal("form_value", r.FormValue("form_key"))

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

func Test_FilterRequest(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/request/filter"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/request/filter", r.RequestURI)
		assertion.Equal("request filter", r.Header.Get("X-Filtered-Header"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request, _ := http.NewRequest(method, ts.URL+"/request/filter", nil)

	client := New(ts.URL, false)
	client.NewFilterRequest(t, request, func(r *http.Request) error {
		r.Header.Add("X-Filtered-Header", "request filter")

		return nil
	})
	client.AssertOK()
	client.AssertEmpty()
}
