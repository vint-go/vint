package fixtures

import (
	"net/http"
	_ "net/http/pprof" // MATCH /blank import of net/http/pprof automatically exposes profiling endpoints on the default mux/
)

func badServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":8080", nil)
}
