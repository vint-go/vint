package fixtures

import "net/http"

// Invalid: ParseMultipartForm without MaxBytesReader
func uploadHandlerUnbounded(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0) // MATCH /form parsing without http.MaxBytesReader allows unbounded memory consumption/
	file, _, err := r.FormFile("upload")
	_ = file
	_ = err
}

// Invalid: ParseForm without MaxBytesReader
func formHandlerUnbounded(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                // MATCH /form parsing without http.MaxBytesReader allows unbounded memory consumption/
	value := r.FormValue("data") // MATCH /form parsing without http.MaxBytesReader allows unbounded memory consumption/
	_ = value
}

// Invalid: FormValue without MaxBytesReader
func formValueUnbounded(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("data") // MATCH /form parsing without http.MaxBytesReader allows unbounded memory consumption/
	_ = value
}

// Valid: ParseMultipartForm with MaxBytesReader
func uploadHandlerBounded(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("upload")
	_ = file
	_ = err
}

// Valid: ParseForm with MaxBytesReader
func formHandlerBounded(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	r.ParseForm()
	value := r.FormValue("data")
	_ = value
}

// Valid: not an HTTP handler function (no *http.Request param)
func notAnHTTPHandler(name string) {
	_ = name
}
