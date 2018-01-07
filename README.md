# httptesting

[![Build Status](https://travis-ci.org/dolab/httptesting.svg?branch=master&style=flat)](https://travis-ci.org/dolab/httptesting)

HTTP testing client of golang for human.

## Installation

```bash
$ go get github.com/dolab/httptesting
```

## Getting Started

```go
package main

import (
    "net/http"
    "testing"

    "github.com/dolab/httptesting"
)

func Test_Client(t *testing.T) {
    host := "https://example.com"
    client := httptesting.New(host, true)

    client.Get("/")

    // verify http response status
    client.AssertOK()

    // verify http response header
    client.AssertExistHeader("Content-Length")

    // verify http response body
    client.AssertNotEmpty()
}

func Test_ClientWithCustomRequest(t *testing.T) {
    r, _ := http.NewRequest("HEAD", "https://example.com", nil)
    r.Header.Add("X-Custom-Header", "custom-header")

    client := httptesting.New("", true)
    client.NewRequest(r)

    // verify http response status
    client.AssertOK()

    // verify http response header
    client.AssertExistHeader("Content-Length")

    // verify http response body
    client.AssertEmpty()
}
```

### Advantage Usage

```go
package main

import (
	"testing"

    "github.com/dolab/httptesting"
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
        request.AssertOK()

        // verify http response header
        request.AssertHeader("X-Mock-Client", "httptesting")

        // verify http response body with json format
        request.AssertContainsJSON("user.name", "httptesting")

        // for array
        request.AssertContainsJSON("addresses.1.name", "USA")
        request.AssertNotContainsJSON("addresses.2.name")

        // use regexp for custom matcher
        request.AssertMatch("user.*")
    })

    t.Run("POST /api/json", func(t *testing.T) {
        request := client.New(t)
        request.WithHeader("X-Mock-Client", "httptesting")

        payload := struct {
            Name string `json:"name"`
            Age  int    `json:"age"`
        }{"httptesting", 3}

        // assume server response with following json data:
        // {"data":{"name":"httptesting","age":3},"success":true}
        request.PostJSON("/api/json", payload)

        // verify http response status
        request.AssertOK()

        // verify http response header
        request.AssertHeader("X-Mock-Client", "httptesting")

        // verify http response body with json format
        request.AssertContainsJSON("data.name", "httptesting")
        request.AssertContainsJSON("data.age", 3)
        request.AssertContainsJSON("success", true)

        // use regexp for custom matcher
        request.AssertNotMatch("user.*")
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
