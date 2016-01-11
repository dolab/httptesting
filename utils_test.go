package httptesting

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url_key": []string{"url_value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("/get?url_key=url_value", r.RequestURI)
		assertion.Equal("url_value", r.FormValue("url_key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.Get(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("GET /get OK!")
}

func Test_Head(t *testing.T) {
	assertion := assert.New(t)
	method := "HEAD"
	uri := "/head?key"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal(uri, r.RequestURI)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.Head(t, uri)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertEmpty()
}

func Test_Options(t *testing.T) {
	assertion := assert.New(t)
	method := "OPTIONS"
	uri := "/options?key"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal(uri, r.RequestURI)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.Options(t, uri)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("OPTIONS /options?key OK!")
}

func Test_PutForm(t *testing.T) {
	assertion := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := url.Values{"form_key": []string{"form_value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assertion.Equal("form_value", r.FormValue("form_key"))

		_, ok := r.URL.Query()["key"]
		assertion.True(ok)
		assertion.Empty(r.FormValue("key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.PutForm(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PUT /put?key OK!")
}

func Test_PutJSON(t *testing.T) {
	assertion := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("application/json", r.Header.Get("Content-Type"))

		_, ok := r.URL.Query()["key"]
		assertion.True(ok)
		assertion.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		assertion.Nil(err)
		assertion.Equal(`{"name":"testing","age":1,"married":false}`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.PutJSON(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PUT /put?key OK!")
}

func Test_PutXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	assertion := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))

		_, ok := r.URL.Query()["key"]
		assertion.True(ok)
		assertion.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		assertion.Nil(err)
		assertion.Equal(`<Person><Name>testing</Name><Age>1</Age><Married>false</Married></Person>`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.PutXML(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PUT /put?key OK!")
}
