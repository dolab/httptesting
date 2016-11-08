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
	"testing"
)

// Get issues a GET request to the given path and stores the result in Response and ResponseBody.
func (test *Client) Get(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		test.Invoke(t, "GET", path, contentType)
	} else {
		test.Invoke(t, "GET", path, contentType, params[0])
	}
}

// Head issues a HEAD request to the given path and stores the result in Response and ResponseBody.
func (test *Client) Head(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		test.Invoke(t, "HEAD", path, contentType)
	} else {
		test.Invoke(t, "HEAD", path, contentType, params[0])
	}
}

// Options issues an OPTIONS request to the given path and stores the result in Response and ResponseBody.
func (test *Client) Options(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		test.Invoke(t, "OPTIONS", path, contentType)
	} else {
		test.Invoke(t, "OPTIONS", path, contentType, params[0])
	}
}

// Put issues a PUT request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (test *Client) Put(t *testing.T, path, contentType string, data ...interface{}) {
	test.Invoke(t, "PUT", path, contentType, data...)
}

// PutForm issues a PUT request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (test *Client) PutForm(t *testing.T, path string, data interface{}) {
	test.Invoke(t, "PUT", path, "application/x-www-form-urlencoded", data)
}

// PutJSON issues a PUT request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (test *Client) PutJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "PUT", path, "application/json", b)
}

// PutXML issues a PUT request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (test *Client) PutXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "PUT", path, "text/xml", b)
}

// Post issues a POST request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (test *Client) Post(t *testing.T, path, contentType string, data ...interface{}) {
	test.Invoke(t, "POST", path, contentType, data...)
}

// PostForm issues a POST request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (test *Client) PostForm(t *testing.T, path string, data interface{}) {
	test.Invoke(t, "POST", path, "application/x-www-form-urlencoded", data)
}

// PostJSON issues a POST request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (test *Client) PostJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "POST", path, "application/json", b)
}

// PostXML issues a POST request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (test *Client) PostXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "POST", path, "text/xml", b)
}

// Patch issues a PATCH request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (test *Client) Patch(t *testing.T, path, contentType string, data ...interface{}) {
	test.Invoke(t, "PATCH", path, contentType, data...)
}

// PatchForm issues a PATCH request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (test *Client) PatchForm(t *testing.T, path string, data interface{}) {
	test.Invoke(t, "PATCH", path, "application/x-www-form-urlencoded", data)
}

// PatchJSON issues a PATCH request to the given path with with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (test *Client) PatchJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "PATCH", path, "application/json", b)
}

// PatchXML issues a PATCH request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (test *Client) PatchXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "PATCH", path, "text/xml", b)
}

// Delete issues a DELETE request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody.
func (test *Client) Delete(t *testing.T, path, contentType string, data ...interface{}) {
	test.Invoke(t, "DELETE", path, contentType, data...)
}

// DeleteForm issues a DELETE request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody.
func (test *Client) DeleteForm(t *testing.T, path string, data interface{}) {
	test.Invoke(t, "DELETE", path, "application/x-www-form-urlencoded", data)
}

// DeleteJSON issues a DELETE request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody.
// It will encode data by json.Marshal before making request.
func (test *Client) DeleteJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "DELETE", path, "application/json", b)
}

// DeleteXML issues a DELETE request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody.
// It will encode data by xml.Marshal before making request.
func (test *Client) DeleteXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	test.Invoke(t, "DELETE", path, "text/xml", b)
}

// Invoke issues a HTTP request to the given path with specified method and content type header, and
// stores the result in Response and ResponseBody.
func (test *Client) Invoke(t *testing.T, method, path, contentType string, data ...interface{}) {
	test.t = t

	var (
		request *http.Request
		err     error
	)

	if len(data) == 0 {
		request, err = http.NewRequest(method, test.Url(path), nil)
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
			urlStr := test.Url(path)

			data, _ := ioutil.ReadAll(reader)
			if len(data) != 0 {
				urlStr += "?" + string(data)
			}

			request, err = http.NewRequest(method, urlStr, nil)

		default:
			request, err = http.NewRequest(method, test.Url(path), reader)
		}
	}

	if err != nil {
		t.Fatalf("%s %s: %#v\n", method, path, err)
	}

	request.Header.Set("Content-Type", contentType)

	contentLength, err := strconv.ParseInt(request.Header.Get("Content-Length"), 10, 64)
	if err == nil {
		request.ContentLength = contentLength
	}

	test.NewSessionRequest(t, request)
}
