package fixtures

import (
	"net/http"
	"time"
)

var handler http.Handler

// Invalid: Server without ReadHeaderTimeout
func badNoTimeout() {
	server := &http.Server{ // MATCH /http.Server missing ReadHeaderTimeout, which can lead to Slowloris attacks/
		Addr:    ":8080",
		Handler: handler,
	}
	_ = server
}

// Invalid: Only setting ReadTimeout is not sufficient
func badOnlyReadTimeout() {
	server := &http.Server{ // MATCH /http.Server missing ReadHeaderTimeout, which can lead to Slowloris attacks/
		Addr:        ":8080",
		Handler:     handler,
		ReadTimeout: 10 * time.Second,
	}
	_ = server
}

// Valid: Server with ReadHeaderTimeout configured
func goodWithReadHeaderTimeout() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	_ = server
}

// Valid: Server with both timeouts configured
func goodWithBothTimeouts() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	_ = server
}
