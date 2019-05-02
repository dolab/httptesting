package httptesting

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golib/assert"
)

func TestRequest_Get(t *testing.T) {
	it := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/html", r.Header.Get("Content-Type"))
		it.Equal("/get?url-key=url-value", r.RequestURI)
		it.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.Get(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("GET /get OK!")
}

func TestRequest_GetJSON(t *testing.T) {
	it := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/json", r.Header.Get("Content-Type"))
		it.Equal("/get?url-key=url-value", r.RequestURI)
		it.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.GetJSON(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("GET /get OK!")
}

func TestRequest_GetXML(t *testing.T) {
	it := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/xml", r.Header.Get("Content-Type"))
		it.Equal("/get?url-key=url-value", r.RequestURI)
		it.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.GetXML(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("GET /get OK!")
}

func TestRequest_Head(t *testing.T) {
	it := assert.New(t)
	method := "HEAD"
	uri := "/head?key"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal(uri, r.RequestURI)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.Head(uri)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertEmpty()
}

func TestRequest_Options(t *testing.T) {
	it := assert.New(t)
	method := "OPTIONS"
	uri := "/options?key"
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal(uri, r.RequestURI)

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.Options(uri)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("OPTIONS /options?key OK!")
}

func TestRequest_PutForm(t *testing.T) {
	it := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)
		it.Equal("form-value", r.FormValue("form-key"))

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PutForm(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PUT /put?key OK!")
}

func TestRequest_PutJSON(t *testing.T) {
	it := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/json", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`{"name":"testing","age":1,"married":false}`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PutJSON(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PUT /put?key OK!")
}

func TestRequest_PutXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	it := assert.New(t)
	method := "PUT"
	uri := "/put?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/xml", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`<Person><Name>testing</Name><Age>1</Age><Married>false</Married></Person>`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PutXML(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PUT /put?key OK!")
}

func TestRequest_PostForm(t *testing.T) {
	it := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)
		it.Equal("form-value", r.FormValue("form-key"))

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PostForm(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("POST /post?key OK!")
}

func TestRequest_PostJSON(t *testing.T) {
	it := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/json", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`{"name":"testing","age":1,"married":false}`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PostJSON(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("POST /post?key OK!")
}

func TestRequest_PostXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	it := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/xml", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`<Person><Name>testing</Name><Age>1</Age><Married>false</Married></Person>`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PostXML(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("POST /post?key OK!")
}

func TestRequest_PatchForm(t *testing.T) {
	it := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)
		it.Equal("form-value", r.FormValue("form-key"))

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PatchForm(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PATCH /patch?key OK!")
}

func TestRequest_PatchJSON(t *testing.T) {
	it := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/json", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`{"name":"testing","age":1,"married":false}`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PatchJSON(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PATCH /patch?key OK!")
}

func TestRequest_PatchXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	it := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/xml", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`<Person><Name>testing</Name><Age>1</Age><Married>false</Married></Person>`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.PatchXML(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("PATCH /patch?key OK!")
}

func TestRequest_DeleteForm(t *testing.T) {
	it := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)
		it.Empty(r.FormValue("form-key"))

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`form-key=form-value`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.DeleteForm(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("DELETE /delete?key OK!")
}

func TestRequest_DeleteJSON(t *testing.T) {
	it := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("application/json", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`{"name":"testing","age":1,"married":false}`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.DeleteJSON(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("DELETE /delete?key OK!")
}

func TestRequest_DeleteXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	it := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		it.Equal(method, r.Method)
		it.Equal("text/xml", r.Header.Get("Content-Type"))
		it.Equal(uri, r.RequestURI)

		_, ok := r.URL.Query()["key"]
		it.True(ok)
		it.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		it.Nil(err)
		it.Equal(`<Person><Name>testing</Name><Age>1</Age><Married>false</Married></Person>`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	request := New(ts.URL, false).New(t)
	request.DeleteXML(uri, params)
	request.AssertOK()
	request.AssertHeader("x-request-method", method)
	request.AssertContains("DELETE /delete?key OK!")
}
