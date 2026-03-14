package fixtures

import (
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {}

// Invalid: using http.ListenAndServe without timeout support
func badListenAndServe() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil) // MATCH /use of http.ListenAndServe with no support for setting timeouts, use http.Server with timeouts instead/
}

// Invalid: using http.ListenAndServeTLS without timeout support
func badListenAndServeTLS() {
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil) // MATCH /use of http.ListenAndServeTLS with no support for setting timeouts, use http.Server with timeouts instead/
}

// Valid: creating server with explicit timeouts
func goodServerWithTimeouts() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           http.DefaultServeMux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	server.ListenAndServe()
}

// Valid: calling ListenAndServe on a server instance (method, not package-level function)
func goodServerMethodCall() {
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServeTLS("cert.pem", "key.pem", nil)
}
