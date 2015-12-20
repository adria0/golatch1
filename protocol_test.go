package golatch1

import "testing"

func TestAuthorization(t *testing.T) {
	secret := "ppxXX9V2DDpzLATFpwRVNkeZQixaD6Q6VnLqsaBp"
	id := "EB4ZmDeZJJHKMnNFZaZT"
	utc := "2006-01-02 15:04:05"
	method := "POST"
	path := "/quoteoftheday"
	paramsAndHeaders := map[string]string{
		"The":             " best preparation",
		"for tomorrow":    "is doing your",
		"X-11paths-best":  "today.",
		"X-11paths-afaik": "structs\nends\nhere.",
	}
	expected := "11PATHS EB4ZmDeZJJHKMnNFZaZT DVzXHF5IvkeEhKRUCuEVUbfmKd0="
	la := NewLatchApp(id, secret)
	auth, _ := la.authHeader(method, path, paramsAndHeaders, paramsAndHeaders, utc)
	if auth != expected {
		t.Errorf("Test failed: expected '%v' got '%v'", expected, auth)
	}
}
