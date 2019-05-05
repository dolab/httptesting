package httptesting

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"
)

// Request defines http client for human usage.
type Request struct {
	*Client

	Response     *http.Response
	ResponseBody []byte

	mux     sync.Mutex
	t       *testing.T
	cookies []*http.Cookie
	header  http.Header
}

// NewRequest returns a new *Request with *Client
func NewRequest(t *testing.T, client *Client) *Request {
	return &Request{
		Client:  client,
		t:       t,
		cookies: []*http.Cookie{},
		header:  http.Header{},
	}
}

// WithHeader sets http header by replace for the request
func (r *Request) WithHeader(key, value string) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.header.Set(key, value)

	return r
}

// WithHttpHeader adds http header for the request
func (r *Request) WithHttpHeader(header http.Header) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()

	for key, values := range header {
		for _, value := range values {
			r.header.Add(key, value)
		}
	}

	return r
}

// WithCookie sets jar for client by replace for the request
func (r *Request) WithCookies(cookies []*http.Cookie) *Request {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.cookies = append(r.cookies, cookies...)

	return r
}

// NewRequest issues any request and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: You have to manage session / cookie data manually.
func (r *Request) NewRequest(request *http.Request, filters ...RequestFilter) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var err error

	r.Response, err = r.NewClient(filters...).Do(request)
	if err != nil {
		r.t.Fatalf("httptesting: NewRequest:%s %s: %v\n", request.Method, request.URL.RequestURI(), err)
	}
	defer r.Response.Body.Close()

	// Read response body if not empty
	r.ResponseBody = []byte{}

	switch r.Response.StatusCode {
	case http.StatusNoContent:
		// ignore

	default:
		r.ResponseBody, err = ioutil.ReadAll(r.Response.Body)
		if err != nil {
			if err != io.EOF {
				r.t.Fatalf("httptesting: NewRequest:%s %s: %v\n", request.Method, request.URL.RequestURI(), err)
			}

			r.t.Logf("httptesting: NewRequest:%s %s: Unexptected response body with io.EOF\n", request.Method, request.URL.RequestURI())
		}
	}
}

// NewSessionRequest issues any request with session / cookie and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
// NOTE: Session data will be added to the request jar for requested host.
func (r *Request) NewSessionRequest(request *http.Request, filters ...RequestFilter) {
	if cookies, err := r.Cookies(); err == nil {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}

	for _, cookie := range r.cookies {
		request.AddCookie(cookie)
	}

	r.NewRequest(request, filters...)
}

// NewMultipartRequest issues a multipart request for the method & fields given and read the response.
// If successful, the caller may examine the Response and ResponseBody properties.
func (r *Request) NewMultipartRequest(method, path, filename string, file interface{}, fields ...map[string]string) {
	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)

	fw, ferr := mw.CreateFormFile("filename", filename)
	if ferr != nil {
		r.t.Fatalf("httptesting: NewMultipartRequest:%s %s: %v\n", method, path, ferr)
	}

	// apply file
	var (
		reader io.Reader
		err    error
	)
	switch f := file.(type) {
	case io.Reader:
		reader = f

	case *os.File:
		reader = f

	case string:
		reader, err = os.Open(f)
		if err != nil {
			r.t.Fatalf("httptesting: NewMultipartRequest:os.Open(%s): %v\n", f, err)
		}

	default:
		r.t.Fatalf("httptesting: NewMultipartRequest:%T<%v>: Unsupported file type\n", file, file)
	}

	if _, err := io.Copy(fw, reader); err != nil {
		r.t.Fatalf("httptesting: NewMultipartRequest:io.Copy(%T, %T): %v\n", fw, file, err)
	}

	// apply fields
	if len(fields) > 0 {
		for key, value := range fields[0] {
			mw.WriteField(key, value)
		}
	}

	// adds the terminating boundary
	mw.Close()

	request, err := http.NewRequest(method, r.Url(path), &buf)
	if err != nil {
		r.t.Fatalf("httptesting: NewMultipartRequest:%s %s: %v\n", method, path, err)
	}
	request.Header.Set("Content-Type", mw.FormDataContentType())

	r.NewRequest(request)
}
