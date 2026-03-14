package fixtures

import (
	"fmt"
	"net/http"
)

// Invalid: URL from variable parameter - potential SSRF
func ssrfGetVariable(url string) (*http.Response, error) {
	return http.Get(url) // MATCH /potential SSRF: URL passed to http.Get is not a hardcoded constant/
}

// Invalid: URL constructed from user input
func ssrfGetConcatenated(userID string) (*http.Response, error) {
	url := "https://api.example.com/users/" + userID
	return http.Get(url) // MATCH /potential SSRF: URL passed to http.Get is not a hardcoded constant/
}

// Invalid: URL from fmt.Sprintf
func ssrfNewRequestVariable(host string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/api/resource", host)
	req, _ := http.NewRequest("GET", url, nil) // MATCH /potential SSRF: URL passed to http.NewRequest is not a hardcoded constant/
	return http.DefaultClient.Do(req)
}

// Invalid: Head with variable URL
func ssrfHeadVariable(url string) (*http.Response, error) {
	return http.Head(url) // MATCH /potential SSRF: URL passed to http.Head is not a hardcoded constant/
}

// Invalid: Post with variable URL
func ssrfPostVariable(url string) (*http.Response, error) {
	return http.Post(url, "application/json", nil) // MATCH /potential SSRF: URL passed to http.Post is not a hardcoded constant/
}

// Invalid: PostForm with variable URL
func ssrfPostFormVariable(url string) (*http.Response, error) {
	return http.PostForm(url, nil) // MATCH /potential SSRF: URL passed to http.PostForm is not a hardcoded constant/
}

// Invalid: NewRequestWithContext with variable URL
func ssrfNewRequestWithContextVariable(ctx interface{}, url string) {
	http.NewRequestWithContext(nil, "GET", url, nil) // MATCH /potential SSRF: URL passed to http.NewRequestWithContext is not a hardcoded constant/
}

// Valid: Hardcoded URL constant - safe
func ssrfGetHardcoded() (*http.Response, error) {
	return http.Get("https://config.internal.example.com/settings")
}

// Valid: Hardcoded URL in Post - safe
func ssrfPostHardcoded() (*http.Response, error) {
	return http.Post("https://api.example.com/data", "application/json", nil)
}

// Valid: Hardcoded URL in NewRequest - safe
func ssrfNewRequestHardcoded() (*http.Response, error) {
	req, _ := http.NewRequest("GET", "https://api.example.com/resource", nil)
	return http.DefaultClient.Do(req)
}

// Valid: Hardcoded URL in Head - safe
func ssrfHeadHardcoded() (*http.Response, error) {
	return http.Head("https://api.example.com/health")
}

// Valid: Hardcoded URL in PostForm - safe
func ssrfPostFormHardcoded() (*http.Response, error) {
	return http.PostForm("https://api.example.com/form", nil)
}
