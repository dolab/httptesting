package httptesting

import (
	"net/http"

	"github.com/golib/assert"
)

// AssertStatus asserts that the response status code is equal to value.
func (r *Request) AssertStatus(status int) bool {
	return assert.EqualValues(r.t, status, r.Response.StatusCode,
		"Expected response status code of %d, but got %d",
		status,
		r.Response.StatusCode,
	)
}

// AssertOK asserts that the response status code is 200.
func (r *Request) AssertOK() bool {
	return r.AssertStatus(http.StatusOK)
}

// AssertForbidden asserts that the response status code is 403.
func (r *Request) AssertForbidden() bool {
	return r.AssertStatus(http.StatusForbidden)
}

// AssertNotFound asserts that the response status code is 404.
func (r *Request) AssertNotFound() bool {
	return r.AssertStatus(http.StatusNotFound)
}

// AssertInternalError asserts that the response status code is 500.
func (r *Request) AssertInternalError() bool {
	return r.AssertStatus(http.StatusInternalServerError)
}

// AssertHeader asserts that the response includes named header with value.
func (r *Request) AssertHeader(name, value string) bool {
	actual := r.Response.Header.Get(name)

	return assert.EqualValues(r.t, value, actual,
		"Expected response header contains %s of %s, but got %s",
		http.CanonicalHeaderKey(name),
		value,
		actual,
	)
}

// AssertContentType asserts that the response includes Content-Type header with value.
func (r *Request) AssertContentType(contentType string) bool {
	return r.AssertHeader("Content-Type", contentType)
}

// AssertExistHeader asserts that the response includes named header.
func (r *Request) AssertExistHeader(name string) bool {
	name = http.CanonicalHeaderKey(name)

	_, ok := r.Response.Header[name]
	if !ok {
		assert.Fail(r.t, "Response header: "+name+" (*required)",
			"Expected response header includes %s",
			name,
		)
	}

	return ok
}

// AssertNotExistHeader asserts that the response does not include named header.
func (r *Request) AssertNotExistHeader(name string) bool {
	name = http.CanonicalHeaderKey(name)

	_, ok := r.Response.Header[name]
	if ok {
		assert.Fail(r.t, "Response header: "+name+" (*not required)",
			"Expected response header does not include %s",
			name,
		)
	}

	return !ok
}

// AssertEmpty asserts that the response body is empty.
func (r *Request) AssertEmpty() bool {
	return assert.Empty(r.t, string(r.ResponseBody))
}

// AssertNotEmpty asserts that the response body is not empty.
func (r *Request) AssertNotEmpty() bool {
	return assert.NotEmpty(r.t, string(r.ResponseBody))
}

// AssertContains asserts that the response body contains the string.
func (r *Request) AssertContains(s string) bool {
	return assert.Contains(r.t, string(r.ResponseBody), s,
		"Expected response body contains %q",
		s,
	)
}

// AssertNotContains asserts that the response body does not contain the string.
func (r *Request) AssertNotContains(s string) bool {
	return assert.NotContains(r.t, string(r.ResponseBody), s,
		"Expected response body does not contain %q",
		s,
	)
}

// AssertMatch asserts that the response body matches the regular expression.
func (r *Request) AssertMatch(re string) bool {
	return assert.Match(r.t, re, r.ResponseBody,
		"Expected response body matches regexp %q",
		re,
	)
}

// AssertNotMatch asserts that the response body does not match the regular expression.
func (r *Request) AssertNotMatch(re string) bool {
	return assert.NotMatch(r.t, re, r.ResponseBody,
		"Expected response body does not match regexp %q",
		re,
	)
}

// AssertContainsJSON asserts that the response body contains JSON value of the key.
func (r *Request) AssertContainsJSON(key string, value interface{}) bool {
	return assert.ContainsJSON(r.t, string(r.ResponseBody), key, value)
}

// AssertNotContainsJSON asserts that the response body dose not contain JSON value of the key.
func (r *Request) AssertNotContainsJSON(key string) bool {
	return assert.NotContainsJSON(r.t, string(r.ResponseBody), key)
}
