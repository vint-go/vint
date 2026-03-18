package fixtures

import "net/url"

func invalidURLs() {
	// Backslashes are not valid in URLs
	u, _ := url.Parse("http:\\\\example.com") // MATCH /invalid URL in url.Parse: backslash is not allowed in URLs/
	_ = u
}

func validURLs() {
	// Correct URL format with forward slashes
	u, _ := url.Parse("http://example.com")
	_ = u

	// HTTPS URL
	u2, _ := url.Parse("https://example.com/path?query=value")
	_ = u2

	// Variable argument (not a literal, so we skip)
	myURL := "http://example.com"
	u3, _ := url.Parse(myURL)
	_ = u3

	// Empty string is valid
	u4, _ := url.Parse("")
	_ = u4

	// Relative URL is valid
	u5, _ := url.Parse("/relative/path")
	_ = u5
}
