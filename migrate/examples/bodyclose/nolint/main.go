package example

import "net/http"

func doRequest() {
	resp, _ := http.Get("https://example.com") //nolint:bodyclose
	_ = resp.StatusCode
}
