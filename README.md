# httptesting

[![CircleCI](https://circleci.com/gh/dolab/httptesting/tree/master.svg?style=svg)](https://circleci.com/gh/dolab/httptesting/tree/master)

Golang HTTP testing client for human.

## Installation

```bash
$ go get github.com/dolab/httptesting
```

## Getting Started

```go
package httptesting

import (
	"net/http"
	"testing"

	"github.com/dolab/httptesting"
)

// Testing with httptesting.Request
func Test_Request(t *testing.T) {
	host := "https://example.com"
	client := httptesting.New(host, true)

	request := client.New(t)
	request.Get("/")

	// verify http response status
	if request.AssertOK() {
	    // verify http response header
	    request.AssertExistHeader("Content-Length")

	    // verify http response body
	    request.AssertNotEmpty()
	}
}
```

### Connected with `httptest.Server`

```go
package httptesting

import (
	"net/http"
	"testing"

	"github.com/dolab/httptesting"
)

type mockServer struct {
	method	string
	path	  string
	assertion func(w http.ResponseWriter, r *http.Request)
}

func (mock *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mock.assertion(w, r)
}

func Test_Server(t *testing.T) {
	method := "GET"
	uri := "/server/https"
	server := &mockServer{
		method: method,
		path: uri,
		assertion: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-request-method", r.Method)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("TLS"))
		},
	}

	// return default client connected with httptest.Server
	ts := httptesting.NewServer(server, true)
	defer ts.Close()

	request := ts.New(t)
	request.Get("/server/https")

	// verify http response status
	if request.AssertOK() {
	    // verify http response header
	    request.AssertExistHeader("Content-Length")

	    // verify http response body
	    request.AssertContains("TLS")
	}
}
```

### Advantage Usage

```go
package main

import (
	"testing"
)

func Test_Request(t *testing.T) {
	host := "https://example.com"
	client := httptesting.New(host, true)

	t.Run("GET /api/json", func(t *testing.T) {
		request := client.New(t)
		request.WithHeader("X-Mock-Client", "httptesting")

		// assume server response with following json data:
		// {"user":{"name":"httptesting","age":3},"addresses":[{"name":"china"},{"name":"USA"}]}
		request.GetJSON("/api/json", nil)

		// verify http response status
		if request.AssertOK() {
		    // verify http response header
		    request.AssertHeader("X-Mock-Client", "httptesting")

		    // verify http response body with json format
		    request.AssertContainsJSON("user.name", "httptesting")

		    // for array
		    request.AssertContainsJSON("addresses.1.name", "USA")
		    request.AssertNotContainsJSON("addresses.2.name")

		    // use regexp for custom matcher
		    request.AssertMatch("user.*")
		}
	})

	t.Run("POST /api/json", func(t *testing.T) {
		request := client.New(t)
		request.WithHeader("X-Mock-Client", "httptesting")

		payload := struct {
			Name string `json:"name"`
			Age  int	`json:"age"`
		}{"httptesting", 3}

		// assume server response with following json data:
		// {"data":{"name":"httptesting","age":3},"success":true}
		request.PostJSON("/api/json", payload)

		// verify http response status
		if request.AssertOK() {
		    // verify http response header
		    request.AssertHeader("X-Mock-Client", "httptesting")

		    // verify http response body with json format
		    request.AssertContainsJSON("data.name", "httptesting")
		    request.AssertContainsJSON("data.age", 3)
		    request.AssertContainsJSON("success", true)

		    // use regexp for custom matcher
		    request.AssertNotMatch("user.*")
		}
	})
}
```

## Author

[Spring MC](https://twitter.com/mcspring)

## LICENSE

```
The MIT License (MIT)

Copyright (c) 2016

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
