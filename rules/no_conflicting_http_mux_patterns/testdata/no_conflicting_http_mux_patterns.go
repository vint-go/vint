package fixtures

import "net/http"

func handleItem(w http.ResponseWriter, r *http.Request)       {}
func handleItemByName(w http.ResponseWriter, r *http.Request)  {}
func createItem(w http.ResponseWriter, r *http.Request)        {}

// Invalid: conflicting wildcard patterns on same mux
func conflictingPatterns() {
	mux := http.NewServeMux()
	mux.HandleFunc("/items/{id}", handleItem)
	mux.HandleFunc("/items/{name}", handleItemByName) // MATCH /pattern "/items/{name}" conflicts with previously registered pattern "/items/{id}"/
}

// Valid: different methods avoid conflict
func differentMethods() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items/{id}", handleItem)
	mux.HandleFunc("POST /items", createItem)
}

// Valid: different literal segments
func differentSegments() {
	mux := http.NewServeMux()
	mux.HandleFunc("/items/{id}", handleItem)
	mux.HandleFunc("/users/{id}", handleItemByName)
}

// Valid: different number of segments
func differentSegmentCount() {
	mux := http.NewServeMux()
	mux.HandleFunc("/items/{id}", handleItem)
	mux.HandleFunc("/items/{id}/details", handleItemByName)
}
