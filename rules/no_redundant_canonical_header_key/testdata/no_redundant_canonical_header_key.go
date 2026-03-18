package fixtures

import "net/http"

func getHeader(h http.Header) string {
	return h.Get(http.CanonicalHeaderKey("content-type")) // MATCH /redundant call to http.CanonicalHeaderKey in method call on http.Header/
}

func setHeader(h http.Header) {
	h.Set(http.CanonicalHeaderKey("content-type"), "text/html") // MATCH /redundant call to http.CanonicalHeaderKey in method call on http.Header/
}

func addHeader(h http.Header) {
	h.Add(http.CanonicalHeaderKey("accept"), "application/json") // MATCH /redundant call to http.CanonicalHeaderKey in method call on http.Header/
}

func delHeader(h http.Header) {
	h.Del(http.CanonicalHeaderKey("x-custom")) // MATCH /redundant call to http.CanonicalHeaderKey in method call on http.Header/
}

func valuesHeader(h http.Header) []string {
	return h.Values(http.CanonicalHeaderKey("accept-encoding")) // MATCH /redundant call to http.CanonicalHeaderKey in method call on http.Header/
}

// Valid cases - no redundant call

func getHeaderDirect(h http.Header) string {
	return h.Get("content-type")
}

func setHeaderDirect(h http.Header) {
	h.Set("content-type", "text/html")
}

func addHeaderDirect(h http.Header) {
	h.Add("accept", "application/json")
}

func delHeaderDirect(h http.Header) {
	h.Del("x-custom")
}

func valuesHeaderDirect(h http.Header) []string {
	return h.Values("accept-encoding")
}

// Using CanonicalHeaderKey outside of Header methods is fine
func canonicalKey() string {
	return http.CanonicalHeaderKey("content-type")
}
