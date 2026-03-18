package fixtures

import "net/http"

func invalid() {
	url := "http://example.com"
	req, err := http.NewRequest("GET", url, nil) // MATCH /use http.NoBody instead of nil in http.NewRequest/
	_, _ = req, err
}

func valid() {
	url := "http://example.com"
	req, err := http.NewRequest("GET", url, http.NoBody)
	_, _ = req, err
}
