package fixtures

import (
	"net/http"
	"strconv"
)

// Invalid: Constructing HTTP request with bare LF characters
func buildRequestBareLF(body string) string {
	return "GET / HTTP/1.1\nHost: example.com\nContent-Length: " + strconv.Itoa(len(body)) + "\n\n" + body // MATCH /potential HTTP request smuggling: use CRLF (\r\n) instead of bare LF (\n) in HTTP requests/
}

// Invalid: Setting conflicting headers on same request
func setConflictingHeaders(url string, body *http.Request) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Length", "10")
	req.Header.Set("Transfer-Encoding", "chunked") // MATCH /potential HTTP request smuggling: conflicting Content-Length and Transfer-Encoding headers/
}

// Invalid: Bare LF in HTTP header string
func buildHeaderBareLF() string {
	return "POST /api HTTP/1.1\nHost: target.com\n\n" // MATCH /potential HTTP request smuggling: use CRLF (\r\n) instead of bare LF (\n) in HTTP requests/
}

// Invalid: Conflicting headers in reverse order
func setConflictingHeadersReverse(url string) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Content-Length", "10") // MATCH /potential HTTP request smuggling: conflicting Content-Length and Transfer-Encoding headers/
}

// Valid: Using standard library's HTTP client which handles headers correctly
func makeRequest(url string) (*http.Response, error) {
	return http.Get(url)
}

// Valid: Using NewRequest with proper body handling (no conflicting headers)
func makePost(url string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Valid: Setting only Content-Length header
func setOnlyContentLength(url string) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Length", "10")
}

// Valid: Setting only Transfer-Encoding header
func setOnlyTransferEncoding(url string) {
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Transfer-Encoding", "chunked")
}

// Valid: String with proper CRLF
func buildRequestCRLF(body string) string {
	return "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n" + body
}

// Valid: Regular string with newline (not HTTP context)
func regularString() string {
	return "hello\nworld"
}

// Valid: Setting headers on different request variables
func differentRequests(url string) {
	req1, _ := http.NewRequest("POST", url, nil)
	req2, _ := http.NewRequest("POST", url, nil)
	req1.Header.Set("Content-Length", "10")
	req2.Header.Set("Transfer-Encoding", "chunked")
}
