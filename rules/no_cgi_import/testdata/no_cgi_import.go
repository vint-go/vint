package fixtures

import (
	"net/http"
	"net/http/cgi" // MATCH /import of net/http/cgi package is not allowed due to known security issues/
)

// Invalid: Using the CGI package
func badCgi() {
	handler := &cgi.Handler{
		Path: "/usr/bin/script",
	}
	http.Handle("/cgi-bin/", handler)
	http.ListenAndServe(":8080", nil)
}

// Valid: Using net/http directly
func goodHttp() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	http.ListenAndServe(":8080", nil)
}
