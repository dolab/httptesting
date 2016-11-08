package httptesting

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Get issues a GET request to the given path and stores the result in Response and ResponseBody.
func (client *RequestClient) Get(urlpath string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		client.Invoke("GET", urlpath, contentType)
	} else {
		client.Invoke("GET", urlpath, contentType, params[0])
	}
}

// Head issues a HEAD request to the given path and stores the result in Response and ResponseBody.
func (client *RequestClient) Head(urlpath string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		client.Invoke("HEAD", urlpath, contentType)
	} else {
		client.Invoke("HEAD", urlpath, contentType, params[0])
	}
}

// Options issues an OPTIONS request to the given path and stores the result in Response and ResponseBody.
func (client *RequestClient) Options(urlpath string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		client.Invoke("OPTIONS", urlpath, contentType)
	} else {
		client.Invoke("OPTIONS", urlpath, "text/html", params[0])
	}
}

// Put issues a PUT request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) Put(urlpath, contentType string, data ...interface{}) {
	client.Invoke("PUT", urlpath, contentType, data...)
}

// PutForm issues a PUT request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) PutForm(urlpath string, data interface{}) {
	client.Invoke("PUT", urlpath, "application/x-www-form-urlencoded", data)
}

// PutJSON issues a PUT request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (client *RequestClient) PutJSON(urlpath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("PUT", urlpath, "application/json", b)
}

// PutXML issues a PUT request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (client *RequestClient) PutXML(urlpath string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("PUT", urlpath, "text/xml", b)
}

// Post issues a POST request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) Post(urlpath, contentType string, data ...interface{}) {
	client.Invoke("POST", urlpath, contentType, data...)
}

// PostForm issues a POST request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) PostForm(urlpath string, data interface{}) {
	client.Invoke("POST", urlpath, "application/x-www-form-urlencoded", data)
}

// PostJSON issues a POST request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (client *RequestClient) PostJSON(urlpath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("POST", urlpath, "application/json", b)
}

// PostXML issues a POST request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (client *RequestClient) PostXML(urlpath string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("POST", urlpath, "text/xml", b)
}

// Patch issues a PATCH request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) Patch(urlpath, contentType string, data ...interface{}) {
	client.Invoke("PATCH", urlpath, contentType, data...)
}

// PatchForm issues a PATCH request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) PatchForm(urlpath string, data interface{}) {
	client.Invoke("PATCH", urlpath, "application/x-www-form-urlencoded", data)
}

// PatchJSON issues a PATCH request to the given path with with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (client *RequestClient) PatchJSON(urlpath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("PATCH", urlpath, "application/json", b)
}

// PatchXML issues a PATCH request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (client *RequestClient) PatchXML(urlpath string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("PATCH", urlpath, "text/xml", b)
}

// Delete issues a DELETE request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) Delete(urlpath, contentType string, data ...interface{}) {
	client.Invoke("DELETE", urlpath, contentType, data...)
}

// DeleteForm issues a DELETE request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) DeleteForm(urlpath string, data interface{}) {
	client.Invoke("DELETE", urlpath, "application/x-www-form-urlencoded", data)
}

// DeleteJSON issues a DELETE request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (client *RequestClient) DeleteJSON(urlpath string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("DELETE", urlpath, "application/json", b)
}

// DeleteXML issues a DELETE request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (client *RequestClient) DeleteXML(urlpath string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		client.t.Fatal(err)
	}

	client.Invoke("DELETE", urlpath, "text/xml", b)
}

// Invoke issues a HTTP request to the given path with specified method and content type header, and
// stores the result in Response and ResponseBody.
func (client *RequestClient) Invoke(method, urlpath, contentType string, data ...interface{}) {
	var (
		request *http.Request
		err     error
	)

	if len(data) == 0 {
		request, err = http.NewRequest(method, client.Url(urlpath), nil)
	} else {
		var reader io.Reader

		body := data[0]
		switch body.(type) {
		case io.Reader:
			reader, _ = body.(io.Reader)

		case string:
			s, _ := body.(string)

			reader = bytes.NewBufferString(s)

		case []byte:
			buf, _ := body.([]byte)

			reader = bytes.NewBuffer(buf)

		case url.Values:
			params, _ := body.(url.Values)

			reader = bytes.NewBufferString(params.Encode())

		default:
			b, _ := json.Marshal(body)

			reader = bytes.NewBuffer(b)
			contentType = "application/json"
		}

		switch method {
		case "GET", "HEAD", "OPTIONS": // apply request params to url
			urlStr := client.Url(urlpath)

			data, _ := ioutil.ReadAll(reader)
			if len(data) != 0 {
				urlStr += "?" + string(data)
			}

			request, err = http.NewRequest(method, urlStr, nil)

		default:
			request, err = http.NewRequest(method, client.Url(urlpath), reader)
		}
	}

	if err != nil {
		client.t.Fatalf("%s %s: %#v\n", method, urlpath, err)
	}

	request.Header.Set("Content-Type", contentType)

	// apply custom headers
	for key, values := range client.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	contentLength, err := strconv.ParseInt(request.Header.Get("Content-Length"), 10, 64)
	if err == nil {
		request.ContentLength = contentLength
	}

	client.NewSessionRequest(request)
}
