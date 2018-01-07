package httptesting

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

// Get issues a GET request to the given path with Content-Type: text/html header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) Get(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		c.Send(t, "GET", path, contentType)
	} else {
		c.Send(t, "GET", path, contentType, params[0])
	}
}

// GetJSON issues a GET request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) GetJSON(t *testing.T, path string, params ...url.Values) {
	contentType := "application/json"

	if len(params) == 0 {
		c.Send(t, "GET", path, contentType)
	} else {
		c.Send(t, "GET", path, contentType, params[0])
	}
}

// GetXML issues a GET request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) GetXML(t *testing.T, path string, params ...url.Values) {
	contentType := "text/xml"

	if len(params) == 0 {
		c.Send(t, "GET", path, contentType)
	} else {
		c.Send(t, "GET", path, contentType, params[0])
	}
}

// Head issues a HEAD request to the given path with Content-Type: text/html header, and
// stores the result in Response if success.
func (c *Client) Head(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		c.Send(t, "HEAD", path, contentType)
	} else {
		c.Send(t, "HEAD", path, contentType, params[0])
	}
}

// Options issues an OPTIONS request to the given path Content-Type: text/html header, and
// stores the result in Response if success.
func (c *Client) Options(t *testing.T, path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		c.Send(t, "OPTIONS", path, contentType)
	} else {
		c.Send(t, "OPTIONS", path, contentType, params[0])
	}
}

// Put issues a PUT request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) Put(t *testing.T, path, contentType string, data ...interface{}) {
	c.Send(t, "PUT", path, contentType, data...)
}

// PutForm issues a PUT request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) PutForm(t *testing.T, path string, data interface{}) {
	c.Put(t, path, "application/x-www-form-urlencoded", data)
}

// PutJSON issues a PUT request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (c *Client) PutJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("[PutJSON] json.Marshal(%T): %v", data, err)
	}

	c.Put(t, path, "application/json", b)
}

// PutXML issues a PUT request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (c *Client) PutXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatalf("[PutXML] xml.Marshal(%T): %v", data, err)
	}

	c.Put(t, path, "text/xml", b)
}

// Post issues a POST request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) Post(t *testing.T, path, contentType string, data ...interface{}) {
	c.Send(t, "POST", path, contentType, data...)
}

// PostForm issues a POST request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) PostForm(t *testing.T, path string, data interface{}) {
	c.Post(t, path, "application/x-www-form-urlencoded", data)
}

// PostJSON issues a POST request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (c *Client) PostJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("[PostJSON] json.Marshal(%T): %v", data, err)
	}

	c.Post(t, path, "application/json", b)
}

// PostXML issues a POST request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (c *Client) PostXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatalf("[PostXML] xml.Marshal(%T): %v", data, err)
	}

	c.Post(t, path, "text/xml", b)
}

// Patch issues a PATCH request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) Patch(t *testing.T, path, contentType string, data ...interface{}) {
	c.Send(t, "PATCH", path, contentType, data...)
}

// PatchForm issues a PATCH request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) PatchForm(t *testing.T, path string, data interface{}) {
	c.Patch(t, path, "application/x-www-form-urlencoded", data)
}

// PatchJSON issues a PATCH request to the given path with with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// It will encode data by json.Marshal before making request.
func (c *Client) PatchJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("[PatchJSON] json.Marshal(%T): %v", data, err)
	}

	c.Patch(t, path, "application/json", b)
}

// PatchXML issues a PATCH request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (c *Client) PatchXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatalf("[PatchXML] xml.Marshal(%T): %v", data, err)
	}

	c.Patch(t, path, "text/xml", b)
}

// Delete issues a DELETE request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) Delete(t *testing.T, path, contentType string, data ...interface{}) {
	c.Send(t, "DELETE", path, contentType, data...)
}

// DeleteForm issues a DELETE request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (c *Client) DeleteForm(t *testing.T, path string, data interface{}) {
	c.Delete(t, path, "application/x-www-form-urlencoded", data)
}

// DeleteJSON issues a DELETE request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (c *Client) DeleteJSON(t *testing.T, path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("[DeleteJSON] json.Marshal(%T): %v", data, err)
	}

	c.Delete(t, path, "application/json", b)
}

// DeleteXML issues a DELETE request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (c *Client) DeleteXML(t *testing.T, path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		t.Fatalf("[DeleteXML] xml.Marshal(%T): %v", data, err)
	}

	c.Delete(t, path, "text/xml", b)
}

// Send issues a HTTP request to the given path with specified method and content type header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data with json.Marshal for unspported types and reset content type to application/json for the request.
func (c *Client) Send(t *testing.T, method, path, contentType string, data ...interface{}) {
	request, err := c.Build(method, path, contentType, data...)
	if err != nil {
		t.Fatalf("[SEND] %s %s: %v\n", method, path, err)
	}

	c.NewSessionRequest(t, request)
}

func (c *Client) Build(method, path, contentType string, data ...interface{}) (request *http.Request, err error) {
	var (
		buf *bytes.Buffer
	)

	if len(data) == 0 {
		buf = bytes.NewBuffer(nil)

		request, err = http.NewRequest(method, c.Url(path), nil)
	} else {
		body := data[0]

		switch body.(type) {
		case io.Reader:
			buf = bytes.NewBuffer(nil)
			_, err = io.Copy(buf, body.(io.Reader))

		case string:
			buf = bytes.NewBufferString(body.(string))

		case []byte:
			buf = bytes.NewBuffer(body.([]byte))

		case url.Values:
			buf = bytes.NewBufferString(body.(url.Values).Encode())

		default:
			b, _ := json.Marshal(body)

			buf = bytes.NewBuffer(b)
			contentType = "application/json"
		}

		absurl := c.Url(path)
		switch method {
		case "GET", "HEAD", "OPTIONS": // apply request params to url
			if buf.Len() > 0 {
				absurl += "?" + buf.String()

				// clean
				buf.Reset()
			}

			request, err = http.NewRequest(method, absurl, nil)

		default:
			request, err = http.NewRequest(method, absurl, buf)
		}
	}

	if err != nil {
		return
	}

	// hijack request headers
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
	request.ContentLength = int64(buf.Len())
	return
}
