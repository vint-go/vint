package example

import "net/http"

func doGet() {
	resp, _ := http.Get("https://example.com") //nolint:noctx
	_ = resp
}

func doRequest() {
	req, _ := http.NewRequest("GET", "https://example.com", nil) //nolint:noctx
	_ = req
}
