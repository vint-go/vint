package fixtures

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
)

// Invalid: Taint flows from request query param via fmt.Fprintf to response - reflected XSS
func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	// Taint flows from request to response - reflected XSS
	fmt.Fprintf(w, "<h1>Hello, %s!</h1>", name) // MATCH /potential XSS: tainted data from user input flows to HTTP response/
}

// Invalid: Taint flows from FormValue via w.Write to response
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	// User input reflected in HTML response
	w.Write([]byte("<p>Results for: " + query + "</p>")) // MATCH /potential XSS: tainted data from user input flows to HTTP response/
}

// Invalid: Taint flows through string concatenation to w.Write
func concatHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	output := "<div>Welcome, " + user + "</div>"
	w.Write([]byte(output)) // MATCH /potential XSS: tainted data from user input flows to HTTP response/
}

// Invalid: Taint flows through fmt.Sprintf to w.Write
func sprintfHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	body := fmt.Sprintf("<p>%s</p>", msg)
	w.Write([]byte(body)) // MATCH /potential XSS: tainted data from user input flows to HTTP response/
}

// Valid: Using html/template for automatic escaping
func safeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	t := template.Must(template.New("page").Parse("<h1>Hello, {{.}}!</h1>"))
	t.Execute(w, name)
}

// Valid: Manually escaping user input with html.EscapeString
func safeEscapeHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	escaped := html.EscapeString(query)
	w.Write([]byte("<p>Results for: " + escaped + "</p>"))
}

// Valid: Writing static content only
func staticHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Static Content</h1>"))
}

// Valid: No HTTP input used in response
func noInputHandler(w http.ResponseWriter, r *http.Request) {
	message := "Hello, World!"
	fmt.Fprintf(w, "<p>%s</p>", message)
}
