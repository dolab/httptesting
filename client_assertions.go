package httptesting

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/golib/assert"
)

// AssertOK asserts that the response status code is 200.
func (c *Client) AssertOK() {
	c.AssertStatus(http.StatusOK)
}

// AssertNotFound asserts that the response status code is 404.
func (c *Client) AssertNotFound() {
	c.AssertStatus(http.StatusNotFound)
}

// AssertStatus asserts that the response status code is equal to value.
func (c *Client) AssertStatus(status int) {
	assert.EqualValues(c.t, status, c.Response.StatusCode, "Expected response status code "+strconv.Itoa(status)+", but got "+c.Response.Status+".")
}

// AssertContentType asserts that the response includes Content-Type header with value.
func (c *Client) AssertContentType(contentType string) {
	c.AssertHeader("Content-Type", contentType)
}

// AssertHeader asserts that the response includes named header with value.
func (c *Client) AssertHeader(name, value string) {
	name = http.CanonicalHeaderKey(name)
	actual := c.Response.Header.Get(name)

	assert.EqualValues(c.t, value, actual, "Expected response header "+name+" with "+value+", but got "+actual+".")
}

// AssertExistHeader asserts that the response includes named header.
func (c *Client) AssertExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := c.Response.Header[name]
	if !ok {
		assert.Fail(c.t, "Response header: "+name+" (*required)", "Expected response header includes "+name+".")
	}
}

// AssertNotExistHeader asserts that the response does not include named header.
func (c *Client) AssertNotExistHeader(name string) {
	name = http.CanonicalHeaderKey(name)

	_, ok := c.Response.Header[name]
	if ok {
		assert.Fail(c.t, "Response header: "+name+" (*not required)", "Expected response header does not include "+name+".")
	}
}

// AssertEmpty asserts that the response body is empty.
func (c *Client) AssertEmpty() {
	assert.Empty(c.t, string(c.ResponseBody))
}

// AssertNotEmpty asserts that the response body is not empty.
func (c *Client) AssertNotEmpty() {
	assert.NotEmpty(c.t, string(c.ResponseBody))
}

// AssertContains asserts that the response body contains the string.
func (c *Client) AssertContains(s string) {
	assert.Contains(c.t, string(c.ResponseBody), s, "Expected response body contains "+s+".")
}

// AssertNotContains asserts that the response body does not contain the string.
func (c *Client) AssertNotContains(s string) {
	assert.NotContains(c.t, string(c.ResponseBody), s, "Expected response body does not contain "+s+".")
}

// AssertMatch asserts that the response body matches the regular expression.
func (c *Client) AssertMatch(re string) {
	r := regexp.MustCompile(re)

	if !r.Match(c.ResponseBody) {
		c.t.Errorf("Expected response body to match regexp %s", re)
	}
}

// AssertNotMatch asserts that the response body does not match the regular expression.
func (c *Client) AssertNotMatch(re string) {
	r := regexp.MustCompile(re)

	if r.Match(c.ResponseBody) {
		c.t.Errorf("Expected response body does not match regexp %s", re)
	}
}

// AssertContainsJSON asserts that the response body contains JSON value of the key.
func (c *Client) AssertContainsJSON(key string, value interface{}) {
	var (
		buf  = c.ResponseBody
		data []byte
		err  error
	)

	for _, yek := range strings.Split(key, ".") {
		data, _, _, err = jsonparser.Get(buf, yek)
		if err == nil {
			buf = data

			continue
		}

		// is the yek an array subscript?
		n, e := strconv.ParseInt(yek, 10, 32)
		if e != nil {
			break
		}

		var i int64 = 0
		jsonparser.ArrayEach(buf, func(arrBuf []byte, arrType jsonparser.ValueType, arrOffset int, arrErr error) {
			if i == n {
				buf = arrBuf
				err = arrErr
			}

			i++
		})
		if err != nil {
			break
		}
	}
	if err != nil {
		c.t.Errorf("Expected response body contains json key %s with %s, but got Errr(%v)", key, value, err)
	}

	actual := string(buf)
	switch value.(type) {
	case []byte:
		expected := string(value.([]byte))
		assert.EqualValues(c.t, expected, actual, "Expected response body contains json key "+key+" with "+expected+", but got "+actual+".")

	case string:
		expected := value.(string)
		assert.EqualValues(c.t, expected, actual, "Expected response body contains json key "+key+" with "+expected+", but got "+actual+".")

	case int8:
		expected := int(value.(int8))
		actualInt, _ := strconv.Atoi(actual)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.Itoa(expected)+", but got "+actual+".")

	case int:
		expected := value.(int)
		actualInt, _ := strconv.Atoi(actual)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.Itoa(expected)+", but got "+actual+".")

	case int16:
		expected := int64(value.(int16))
		actualInt, _ := strconv.ParseInt(actual, 10, 16)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case int32:
		expected := int64(value.(int32))
		actualInt, _ := strconv.ParseInt(actual, 10, 32)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case int64:
		expected := value.(int64)
		actualInt, _ := strconv.ParseInt(actual, 10, 64)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatInt(expected, 10)+", but got "+actual+".")

	case float32:
		expected := float64(value.(float32))
		actualInt, _ := strconv.ParseFloat(actual, 32)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatFloat(expected, 'f', 5, 32)+", but got "+actual+".")

	case float64:
		expected := value.(float64)
		actualInt, _ := strconv.ParseFloat(actual, 64)
		assert.EqualValues(c.t, expected, actualInt, "Expected response body contains json key "+key+" with "+strconv.FormatFloat(expected, 'f', 5, 64)+", but got "+actual+".")

	case bool:
		expected := value.(bool)
		switch actual {
		case "true", "True", "1", "on":
			assert.True(c.t, expected, "Expected response body contains json key "+key+" with [true|True|1|on], but got "+actual+".")

		default:
			assert.False(c.t, expected, "Expected response body contains json key "+key+" with [false|False|0|off], but got "+actual+".")
		}
	}
}

// AssertNotContainsJSON asserts that the response body dose not contain JSON value of the key.
func (c *Client) AssertNotContainsJSON(key string) {
	var (
		buf  = c.ResponseBody
		data []byte
		err  error
	)

	for _, yek := range strings.Split(key, ".") {
		data, _, _, err = jsonparser.Get(buf, yek)
		if err == nil {
			buf = data

			continue
		}

		// is the yek an array subscript?
		n, e := strconv.ParseInt(yek, 10, 32)
		if e != nil {
			break
		}

		var i int64 = 0
		jsonparser.ArrayEach(buf, func(arrBuf []byte, arrType jsonparser.ValueType, arrOffset int, arrErr error) {
			if i == n {
				buf = arrBuf
				err = arrErr
			}

			i++
		})
		if err != nil {
			break
		}
	}

	if err == nil {
		c.t.Errorf("Expected response body does not contain json key %s, but got %s", key, string(buf))
	}
}
