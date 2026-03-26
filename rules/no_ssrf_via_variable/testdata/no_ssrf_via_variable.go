package fixtures

import (
	"fmt"
	"net/http"
)

// Invalid: URL from variable parameter - potential SSRF
func ssrfGetVariable(url string) (*http.Response, error) {
	return http.Get(url) // MATCH /potential SSRF: URL passed to http.Get is not a hardcoded constant/
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

// Valid: Hardcoded URL constant - safe
func ssrfGetHardcoded() (*http.Response, error) {
	return http.Get("https://config.internal.example.com/settings")
}

// Valid: Hardcoded URL in Post - safe
func ssrfPostHardcoded() (*http.Response, error) {
	return http.Post("https://api.example.com/data", "application/json", nil)
}

// Valid: Hardcoded URL in Head - safe
func ssrfHeadHardcoded() (*http.Response, error) {
	return http.Head("https://api.example.com/health")
}

// Valid: Hardcoded URL in PostForm - safe
func ssrfPostFormHardcoded() (*http.Response, error) {
	return http.PostForm("https://api.example.com/form", nil)
}

// Invalid: URL from variable (assigned from concatenation) - flagged because url is a variable ident
func ssrfGetConcatenated(userID string) (*http.Response, error) {
	url := "https://api.example.com/users/" + userID
	return http.Get(url) // MATCH /potential SSRF: URL passed to http.Get is not a hardcoded constant/
}

// Valid: URL from fmt.Sprintf (call expression) - not flagged by gosec G107
func ssrfGetSprintf(host string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("https://%s/api/resource", host))
}

// Valid: NewRequest is not monitored (not in gosec G107 scope)
func ssrfNewRequestVariable(host string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/api/resource", host)
	req, _ := http.NewRequest("GET", url, nil)
	return http.DefaultClient.Do(req)
}

// Valid: NewRequestWithContext is not monitored (not in gosec G107 scope)
func ssrfNewRequestWithContextVariable(ctx interface{}, url string) {
	http.NewRequestWithContext(nil, "GET", url, nil)
}

// Valid: Constant identifier is safe
const safeURL = "https://api.example.com/safe"

func ssrfGetConstant() (*http.Response, error) {
	return http.Get(safeURL)
}
