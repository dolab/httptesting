package httptesting

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golib/assert"
)

// AssertOK tests that the response status code is 200.
func (test *Client) AssertOK() {
	test.AssertStatus(http.StatusOK)
}

// AssertNotFound tests that the response status code is 404.
func (test *Client) AssertNotFound() {
	test.AssertStatus(http.StatusNotFound)
}

// AssertStatus tests that the response status code is equal with the given.
func (test *Client) AssertStatus(status int) {
	assert.EqualValues(test.t, status, test.Response.StatusCode, "Expected response status code "+strconv.Itoa(status)+", but got "+test.Response.Status+".")
}

// AssertContentType tests that the response includes Content-Type header with the given value.
func (test *Client) AssertContentType(contentType string) {
	test.AssertHeader("Content-Type", contentType)
}

// AssertHeader tests that the response includes named header with the given value.
func (test *Client) AssertHeader(name, value string) {
	name = http.CanonicalHeaderKey(name)
	actual := test.Response.Header.Get(name)

	assert.EqualValues(test.t, value, actual, "Expected response header "+name+" with "+value+", but got "+actual+".")
}

// AssertExistHeader tests that the response includes named header.
func (test *Client) AssertExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := test.Response.Header[name]
	if !ok {
		assert.Fail(test.t, "Response header: "+name+" (*required)", "Expected response header includes "+name+".")
	}
}

// AssertNotExistHeader tests that the response does not include named header.
func (test *Client) AssertNotExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := test.Response.Header[name]
	if ok {
		assert.Fail(test.t, "Response header: "+name+" (*not required)", "Expected response header does not include "+name+".")
	}
}

// AssertEmpty tests that the response is empty.
func (test *Client) AssertEmpty() {
	assert.Empty(test.t, string(test.ResponseBody))
}

// AssertNotEmpty tests that the response is not empty.
func (test *Client) AssertNotEmpty() {
	assert.NotEmpty(test.t, string(test.ResponseBody))
}

// AssertContains tests that the response contains the given string.
func (test *Client) AssertContains(s string) {
	assert.Contains(test.t, string(test.ResponseBody), s, "Expected response body contains "+s+".")
}

// AssertNotContains tests that the response does not contain the given string.
func (test *Client) AssertNotContains(s string) {
	assert.NotContains(test.t, string(test.ResponseBody), s, "Expected response body does not contain "+s+".")
}

// AssertMatch tests that the response matches the given regular expression.
func (test *Client) AssertMatch(re string) {
	r := regexp.MustCompile(re)

	if !r.Match(test.ResponseBody) {
		test.t.Errorf("Expected response body to match regexp %s", re)
	}
}

// AssertNotMatch tests that the response does not match the given regular expression.
func (test *Client) AssertNotMatch(re string) {
	r := regexp.MustCompile(re)

	if r.Match(test.ResponseBody) {
		test.t.Errorf("Expected response body does not match regexp %s", re)
	}
}

func (test *Client) AssertContainsJSON(key string, value interface{}) {
	var (
		buf = test.ResponseBody
		err error
	)

	keys := strings.Split(key, ".")
	for _, yek := range keys {
		// is the yek a array subscript?
		n, e := strconv.ParseInt(yek, 10, 32)
		if e != nil {
			buf, _, _, err = jsonparser.Get(buf, yek)
		} else {
			var i int64 = 0
			_, err = jsonparser.ArrayEach(buf, func(arrBuf []byte, arrType jsonparser.ValueType, arrOffset int, arrErr error) {
				if i == n {
					buf = arrBuf
					err = arrErr
				}

				i += 1
			})
		}

		if err != nil {
			test.t.Errorf("Expected response body contains json key %s with %s, but got Errr(%v)", key, value, err)
			break
		}
	}

	actual := string(buf)
	switch value.(type) {
	case []byte:
		expected := string(value.([]byte))
		assert.EqualValues(test.t, expected, actual, "Expected response body contains json key "+key+" with "+expected+", but got "+actual+".")

	case string:
		expected := value.(string)
		assert.EqualValues(test.t, expected, actual, "Expected response body contains json key "+key+" with "+expected+", but got "+actual+".")

	case int8:
		expected := int(value.(int8))
		actualInt, _ := strconv.Atoi(actual)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.Itoa(expected)+", but got "+actual+".")

	case int:
		expected := value.(int)
		actualInt, _ := strconv.Atoi(actual)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.Itoa(expected)+", but got "+actual+".")

	case int16:
		expected := int64(value.(int16))
		actualInt, _ := strconv.ParseInt(actual, 10, 16)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case int32:
		expected := int64(value.(int32))
		actualInt, _ := strconv.ParseInt(actual, 10, 32)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case int64:
		expected := value.(int64)
		actualInt, _ := strconv.ParseInt(actual, 10, 64)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case float32:
		expected := float64(value.(float32))
		actualInt, _ := strconv.ParseFloat(actual, 32)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatFloat(expected, 'f', 5, 32)+", but got "+actual+".")

	case float64:
		expected := value.(float64)
		actualInt, _ := strconv.ParseFloat(actual, 64)
		assert.EqualValues(test.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatFloat(expected, 'f', 5, 64)+", but got "+actual+".")

	case bool:
		expected := value.(bool)
		switch actual {
		case "true", "True", "1", "on":
			assert.True(test.t, expected, "Expected response body contains json key "+key+" with [true|True|1|on], but got "+actual+".")

		default:
			assert.False(test.t, expected, "Expected response body contains json key "+key+" with [false|False|0|off], but got "+actual+".")
		}
	}
}

func (test *Client) AssertNotContainsJSON(key string) {
	var (
		buf = test.ResponseBody
		err error
	)

	keys := strings.Split(key, ".")
	for _, yek := range keys {
		// is the yek a array subscript?
		n, e := strconv.ParseInt(yek, 10, 32)
		if e != nil {
			buf, _, _, err = jsonparser.Get(buf, yek)
		} else {
			var i int64 = 0
			_, err = jsonparser.ArrayEach(buf, func(arrBuf []byte, arrType jsonparser.ValueType, arrOffset int, arrErr error) {
				if i == n {
					buf = arrBuf
					err = arrErr
				}

				i += 1
			})
		}
	}

	if err == nil {
		test.t.Errorf("Expected response body does not contain json key %s, but got %s", key, string(buf))
	}
}
