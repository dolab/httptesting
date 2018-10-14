package httptesting

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/dolab/httptesting/internal"
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
	assertion.Nil(client.Response)
	assertion.Empty(client.ResponseBody)
	assertion.Nil(client.ts)
	assertion.NotNil(client.client)
	assertion.Nil(client.t)
	assertion.NotEmpty(client.host)
	assertion.True(client.https)
	assertion.Equal(host, client.Host())
	assertion.Equal(url, client.Url(""))
	assertion.Equal(ws, client.WebsocketUrl(""))

	request := client.New(t)
	assertion.Nil(request.Response)
	assertion.Empty(request.ResponseBody)
	assertion.Nil(request.ts)
	assertion.NotNil(request.client)
	assertion.NotNil(request.t)
	assertion.Equal(client.host, request.host)
	assertion.True(request.https)
}

func Test_NewServer(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/server"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TLS"))
	})

	client := NewServer(server, true)
	defer client.Close()

	assertion.Nil(client.Response)
	assertion.Empty(client.ResponseBody)
	assertion.NotNil(client.ts)
	assertion.NotNil(client.client)
	assertion.Nil(client.t)
	assertion.NotEmpty(client.host)
	assertion.True(client.https)

	request := client.New(t)
	assertion.Nil(request.Response)
	assertion.Empty(request.ResponseBody)
	assertion.Nil(request.ts)
	assertion.NotNil(request.client)
	assertion.NotNil(request.t)
	assertion.Equal(client.host, request.host)
	assertion.True(request.https)
}

func Test_NewServerWithTLS(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/server/https"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("TLS"))
	})

	cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
	assertion.Nil(err)

	ts := NewServerWithTLS(server, cert)
	defer ts.Close()

	// it should work with internal client
	request := ts.New(t)
	request.Get("/server/https", nil)
	request.AssertOK()
	request.AssertContains("TLS")

	// it should work with custom TLS client
	x509cert, err := x509.ParseCertificate(cert.Certificate[0])
	assertion.Nil(err)

	client := NewWithTLS(ts.Url(""), x509cert)

	request = client.New(t)
	request.Get("/server/https", nil)
	request.AssertOK()
	request.AssertContains("TLS")
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
