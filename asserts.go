package httptesting

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/stretchr/testify/assert"
)

// AssertOK tests that the response status code is 200.
func (test *HTTPTesting) AssertOK() {
	test.AssertStatus(http.StatusOK)
}

// AssertNotFound tests that the response status code is 404.
func (test *HTTPTesting) AssertNotFound() {
	test.AssertStatus(http.StatusNotFound)
}

// AssertStatus tests that the response status code is equal with the given.
func (test *HTTPTesting) AssertStatus(status int) {
	assert.EqualValues(test.t, status, test.Response.StatusCode, "Expected response status code "+strconv.Itoa(status)+", but got "+test.Response.Status+".")
}

// AssertContentType tests that the response includes Content-Type header with the given value.
func (test *HTTPTesting) AssertContentType(contentType string) {
	test.AssertHeader("Content-Type", contentType)
}

// AssertHeader tests that the response includes named header with the given value.
func (test *HTTPTesting) AssertHeader(name, value string) {
	name = http.CanonicalHeaderKey(name)
	actual := test.Response.Header.Get(name)

	assert.EqualValues(test.t, value, actual, "Expected response header "+name+" with "+value+", but got "+actual+".")
}

// AssertExistHeader tests that the response includes named header.
func (test *HTTPTesting) AssertExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := test.Response.Header[name]
	if !ok {
		assert.Fail(test.t, "Response header: "+name+" (*required)", "Expected response header includes "+name+".")
	}
}

// AssertNotExistHeader tests that the response does not include named header.
func (test *HTTPTesting) AssertNotExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := test.Response.Header[name]
	if ok {
		assert.Fail(test.t, "Response header: "+name+" (*not required)", "Expected response header does not include "+name+".")
	}
}

// AssertEmpty tests that the response is empty.
func (test *HTTPTesting) AssertEmpty() {
	assert.Empty(test.t, string(test.ResponseBody))
}

// AssertNotEmpty tests that the response is not empty.
func (test *HTTPTesting) AssertNotEmpty() {
	assert.NotEmpty(test.t, string(test.ResponseBody))
}

// AssertContains tests that the response contains the given string.
func (test *HTTPTesting) AssertContains(s string) {
	assert.Contains(test.t, string(test.ResponseBody), s, "Expected response body contains "+s+".")
}

// AssertNotContains tests that the response does not contain the given string.
func (test *HTTPTesting) AssertNotContains(s string) {
	assert.NotContains(test.t, string(test.ResponseBody), s, "Expected response body does not contain "+s+".")
}

// AssertMatch tests that the response matches the given regular expression.
func (test *HTTPTesting) AssertMatch(re string) {
	r := regexp.MustCompile(re)

	if !r.Match(test.ResponseBody) {
		test.t.Errorf("Expected response body to match regexp %s", re)
	}
}

// AssertNotMatch tests that the response does not match the given regular expression.
func (test *HTTPTesting) AssertNotMatch(re string) {
	r := regexp.MustCompile(re)

	if r.Match(test.ResponseBody) {
		test.t.Errorf("Expected response body does not match regexp %s", re)
	}
}
