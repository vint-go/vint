package fixtures

import "net/http"

func handlerBad(w http.ResponseWriter, r *http.Request) {
	// Non-canonical key - will not match
	ct := r.Header["content-type"] // MATCH /non-canonical header key "content-type", use "Content-Type" instead/
	_ = ct
}

func handlerGood(w http.ResponseWriter, r *http.Request) {
	// Using Get() which handles canonicalization
	ct := r.Header.Get("content-type")
	_ = ct

	// Or using canonical form directly
	ct2 := r.Header["Content-Type"]
	_ = ct2
}

func handlerMoreBad(w http.ResponseWriter, r *http.Request) {
	val := r.Header["accept-encoding"] // MATCH /non-canonical header key "accept-encoding", use "Accept-Encoding" instead/
	_ = val

	val2 := r.Header["x-forwarded-for"] // MATCH /non-canonical header key "x-forwarded-for", use "X-Forwarded-For" instead/
	_ = val2
}

func handlerVariable(w http.ResponseWriter, r *http.Request) {
	// Variable key - cannot check at compile time, no failure
	key := "content-type"
	val := r.Header[key]
	_ = val
}

func handlerCanonical(w http.ResponseWriter, r *http.Request) {
	// All canonical forms - no failures
	ct := r.Header["Content-Type"]
	_ = ct
	ae := r.Header["Accept-Encoding"]
	_ = ae
	xff := r.Header["X-Forwarded-For"]
	_ = xff
}
