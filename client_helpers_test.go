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

func Test_Get(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/html", r.Header.Get("Content-Type"))
		assertion.Equal("/get?url-key=url-value", r.RequestURI)
		assertion.Equal("url-value", r.URL.Query().Get("url-key"))

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

func Test_GetJSON(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/json", r.Header.Get("Content-Type"))
		assertion.Equal("/get?url-key=url-value", r.RequestURI)
		assertion.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.GetJSON(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("GET /get OK!")
}

func Test_GetXML(t *testing.T) {
	assertion := assert.New(t)
	method := "GET"
	uri := "/get"
	params := url.Values{"url-key": []string{"url-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))
		assertion.Equal("/get?url-key=url-value", r.RequestURI)
		assertion.Equal("url-value", r.URL.Query().Get("url-key"))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.GetXML(t, uri, params)
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
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("form-value", r.FormValue("form-key"))

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
		assertion.Equal("application/json", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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

func Test_PostForm(t *testing.T) {
	assertion := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("form-value", r.FormValue("form-key"))

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
	client.PostForm(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("POST /post?key OK!")
}

func Test_PostJSON(t *testing.T) {
	assertion := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/json", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.PostJSON(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("POST /post?key OK!")
}

func Test_PostXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	assertion := assert.New(t)
	method := "POST"
	uri := "/post?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.PostXML(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("POST /post?key OK!")
}

func Test_PatchForm(t *testing.T) {
	assertion := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)
		assertion.Equal("form-value", r.FormValue("form-key"))

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
	client.PatchForm(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PATCH /patch?key OK!")
}

func Test_PatchJSON(t *testing.T) {
	assertion := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/json", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.PatchJSON(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PATCH /patch?key OK!")
}

func Test_PatchXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	assertion := assert.New(t)
	method := "PATCH"
	uri := "/patch?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.PatchXML(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("PATCH /patch?key OK!")
}

func Test_DeleteForm(t *testing.T) {
	assertion := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := url.Values{"form-key": []string{"form-value"}}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)
		assertion.Empty(r.FormValue("form-key"))

		_, ok := r.URL.Query()["key"]
		assertion.True(ok)
		assertion.Empty(r.FormValue("key"))

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		assertion.Nil(err)
		assertion.Equal(`form-key=form-value`, string(b))

		w.Header().Set("x-request-method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(method + " " + uri + " OK!"))
	})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := New(ts.URL, false)
	client.DeleteForm(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("DELETE /delete?key OK!")
}

func Test_DeleteJSON(t *testing.T) {
	assertion := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Married bool   `json:"married"`
	}{"testing", 1, false}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("application/json", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.DeleteJSON(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("DELETE /delete?key OK!")
}

func Test_DeleteXML(t *testing.T) {
	type xmlData struct {
		XMLName xml.Name `xml:"Person"`
		Name    string   `xml:"Name"`
		Age     int      `xml:"Age"`
		Married bool     `xml:"Married"`
	}

	assertion := assert.New(t)
	method := "DELETE"
	uri := "/delete?key"
	params := xmlData{
		Name:    "testing",
		Age:     1,
		Married: false,
	}
	server := newMockServer(method, uri, func(w http.ResponseWriter, r *http.Request) {
		assertion.Equal(method, r.Method)
		assertion.Equal("text/xml", r.Header.Get("Content-Type"))
		assertion.Equal(uri, r.RequestURI)

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
	client.DeleteXML(t, uri, params)
	client.AssertOK()
	client.AssertHeader("x-request-method", method)
	client.AssertContains("DELETE /delete?key OK!")
}
