package httptesting

import "testing"

func Test_AssertContainsJSON(t *testing.T) {
	client := New("https://httptesting.example.com", true)
	client.t = t
	client.ResponseBody = []byte(`{"user":{"name":"httptesting","age":3},"addresses":[{"name":"china"},{"name":"USA"}]}`)

	client.AssertContainsJSON("user.name", "httptesting")
	client.AssertContainsJSON("addresses.1.name", "USA")
	client.AssertNotContainsJSON("addresses.0.post")
	client.AssertNotContainsJSON("addresses.3.name")
}
