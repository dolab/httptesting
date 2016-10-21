package httptesting

import "testing"

func Test_AssertContainsJSON(t *testing.T) {
	client := New("https://httptesting.example.com", true)
	client.t = t
	client.ResponseBody = []byte(`{"user":{"name":"httptesting","age":3}}`)

	client.AssertContainsJSON("user.name", "httptesting")
}
