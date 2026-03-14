package fixtures

import "net/http"

// Invalid: serving the entire filesystem root
func servingRoot() {
	http.Handle("/", http.FileServer(http.Dir("/"))) // MATCH /serving the entire filesystem root via http.Dir is a security risk/
	http.ListenAndServe(":8080", nil)
}

// Valid: serving a specific, restricted directory
func servingSpecificDir() {
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./public"))))
	http.ListenAndServe(":8080", nil)
}

// Valid: serving a relative path
func servingRelativePath() {
	http.Handle("/assets/", http.FileServer(http.Dir("./assets")))
	http.ListenAndServe(":8080", nil)
}

// Valid: serving an absolute path that is not root
func servingAbsolutePath() {
	http.Handle("/files/", http.FileServer(http.Dir("/var/www/html")))
	http.ListenAndServe(":8080", nil)
}
