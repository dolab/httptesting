package httptesting

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
)

// Get issues a GET request to the given path with Content-Type: text/html header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) Get(path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		r.Send("GET", path, contentType)
	} else {
		r.Send("GET", path, contentType, params[0])
	}
}

// GetJSON issues a GET request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) GetJSON(path string, params ...url.Values) {
	contentType := "application/json"

	if len(params) == 0 {
		r.Send("GET", path, contentType)
	} else {
		r.Send("GET", path, contentType, params[0])
	}
}

// GetXML issues a GET request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) GetXML(path string, params ...url.Values) {
	contentType := "text/xml"

	if len(params) == 0 {
		r.Send("GET", path, contentType)
	} else {
		r.Send("GET", path, contentType, params[0])
	}
}

// Head issues a HEAD request to the given path with Content-Type: text/html header, and
// stores the result in Response if success.
func (r *Request) Head(path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		r.Send("HEAD", path, contentType)
	} else {
		r.Send("HEAD", path, contentType, params[0])
	}
}

// Options issues an OPTIONS request to the given path Content-Type: text/html header, and
// stores the result in Response if success.
func (r *Request) Options(path string, params ...url.Values) {
	contentType := "text/html"

	if len(params) == 0 {
		r.Send("OPTIONS", path, contentType)
	} else {
		r.Send("OPTIONS", path, contentType, params[0])
	}
}

// Put issues a PUT request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) Put(path, contentType string, data ...interface{}) {
	r.Send("PUT", path, contentType, data...)
}

// PutForm issues a PUT request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) PutForm(path string, data interface{}) {
	r.Put(path, "application/x-www-form-urlencoded", data)
}

// PutJSON issues a PUT request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (r *Request) PutJSON(path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PutJSON] json.Marshal(%T): %v", data, err)
	}

	r.Put(path, "application/json", b)
}

// PutXML issues a PUT request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (r *Request) PutXML(path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PutXML] xml.Marshal(%T): %v", data, err)
	}

	r.Put(path, "text/xml", b)
}

// Post issues a POST request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) Post(path, contentType string, data ...interface{}) {
	r.Send("POST", path, contentType, data...)
}

// PostForm issues a POST request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) PostForm(path string, data interface{}) {
	r.Post(path, "application/x-www-form-urlencoded", data)
}

// PostJSON issues a POST request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (r *Request) PostJSON(path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PostJSON] json.Marshal(%T): %v", data, err)
	}

	r.Post(path, "application/json", b)
}

// PostXML issues a POST request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (r *Request) PostXML(path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PostXML] xml.Marshal(%T): %v", data, err)
	}

	r.Post(path, "text/xml", b)
}

// Patch issues a PATCH request to the given path with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) Patch(path, contentType string, data ...interface{}) {
	r.Send("PATCH", path, contentType, data...)
}

// PatchForm issues a PATCH request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) PatchForm(path string, data interface{}) {
	r.Patch(path, "application/x-www-form-urlencoded", data)
}

// PatchJSON issues a PATCH request to the given path with with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// It will encode data by json.Marshal before making request.
func (r *Request) PatchJSON(path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PatchJSON] json.Marshal(%T): %v", data, err)
	}

	r.Patch(path, "application/json", b)
}

// PatchXML issues a PATCH request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (r *Request) PatchXML(path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		r.t.Fatalf("[PatchXML] xml.Marshal(%T): %v", data, err)
	}

	r.Patch(path, "text/xml", b)
}

// Delete issues a DELETE request to the given path, sending request with specified Content-Type header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) Delete(path, contentType string, data ...interface{}) {
	r.Send("DELETE", path, contentType, data...)
}

// DeleteForm issues a DELETE request to the given path with Content-Type: application/x-www-form-urlencoded header, and
// stores the result in Response and ResponseBody if success.
func (r *Request) DeleteForm(path string, data interface{}) {
	r.Delete(path, "application/x-www-form-urlencoded", data)
}

// DeleteJSON issues a DELETE request to the given path with Content-Type: application/json header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by json.Marshal before making request.
func (r *Request) DeleteJSON(path string, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		r.t.Fatalf("[DeleteJSON] json.Marshal(%T): %v", data, err)
	}

	r.Delete(path, "application/json", b)
}

// DeleteXML issues a DELETE request to the given path with Content-Type: text/xml header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data by xml.Marshal before making request.
func (r *Request) DeleteXML(path string, data interface{}) {
	b, err := xml.Marshal(data)
	if err != nil {
		r.t.Fatalf("[DeleteXML] xml.Marshal(%T): %v", data, err)
	}

	r.Delete(path, "text/xml", b)
}

// Send issues a HTTP request to the given path with specified method and content type header, and
// stores the result in Response and ResponseBody if success.
// NOTE: It will encode data with json.Marshal for unspported types and reset content type to application/json for the request.
func (r *Request) Send(method, path, contentType string, data ...interface{}) {
	request, err := r.Build(method, path, contentType, data...)
	if err != nil {
		r.t.Fatalf("[SEND] %s %s: %v\n", method, path, err)
	}

	// adjust custom headers
	for key, values := range r.header {
		// ignore Content-Type and Content-Length headers
		switch http.CanonicalHeaderKey(key) {
		case "Content-Type", "Content-Length":
			// ignore

		default:
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}
	}

	r.NewSessionRequest(request)
}
