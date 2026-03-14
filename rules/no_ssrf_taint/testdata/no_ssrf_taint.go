package fixtures

import (
	"fmt"
	"io"
	"net/http"
)

// Invalid: Taint flows from HTTP request through function call to outbound HTTP request
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	// Taint flows from request parameter to outbound HTTP request via function call
	resp := fetchURL(targetURL) // MATCH /potential SSRF: tainted data from user input flows to outbound HTTP request/
	io.Copy(w, resp.Body)
}

func fetchURL(url string) *http.Response {
	resp, _ := http.Get(url)
	return resp
}

// Invalid: Direct taint from HTTP request to http.Get
func directGetHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	resp, _ := http.Get(target) // MATCH /potential SSRF: tainted data from user input flows to outbound HTTP request/
	io.Copy(w, resp.Body)
}

// Invalid: Taint flows from FormValue through fmt.Sprintf to http.Post
func postHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.FormValue("endpoint")
	url := fmt.Sprintf("http://%s/api", endpoint)
	http.Post(url, "application/json", nil) // MATCH /potential SSRF: tainted data from user input flows to outbound HTTP request/
}

// Invalid: Taint flows to http.NewRequest
func newRequestHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	http.NewRequest("GET", target, nil) // MATCH /potential SSRF: tainted data from user input flows to outbound HTTP request/
}

// Valid: Hardcoded URL is safe
func safeHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://api.example.com/data")
	io.Copy(w, resp.Body)
	_ = resp
}

// Valid: URL validated against allowlist
var allowedHosts = map[string]bool{
	"api.example.com": true,
	"cdn.example.com": true,
}

func validatedHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")

	if !allowedHosts[targetURL] {
		http.Error(w, "host not allowed", http.StatusForbidden)
		return
	}

	resp, _ := http.Get("https://api.example.com/safe")
	io.Copy(w, resp.Body)
	_ = resp
}

// Valid: No HTTP input, normal function
func fetchData() {
	resp, _ := http.Get("https://internal.example.com/data")
	_ = resp
}
