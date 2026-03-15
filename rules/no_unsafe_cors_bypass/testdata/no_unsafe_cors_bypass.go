package fixtures

import "net/http"

// Invalid: Reflecting the request origin without validation
func reflectOrigin(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin) // MATCH /unsafe CORS bypass: Access-Control-Allow-Origin reflects the request Origin header without validation/
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// Invalid: Using wildcard with credentials
func wildcardOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // MATCH /unsafe CORS bypass: Access-Control-Allow-Origin set to wildcard "*"/
}

// Invalid: Inline r.Header.Get("Origin") directly in Set call
func inlineReflectOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // MATCH /unsafe CORS bypass: Access-Control-Allow-Origin reflects the request Origin header without validation/
}

// Valid: Using an allowlist to validate the origin
var allowedOrigins = map[string]bool{
	"https://app.example.com":   true,
	"https://admin.example.com": true,
}

func safeHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

// Valid: Setting a specific origin
func specificOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://example.com")
}

// Valid: Setting a non-CORS header
func nonCorsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
