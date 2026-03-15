package fixtures

import (
	"fmt"
	"net/http"
	"text/template"
)

// Invalid: Direct taint from HTTP request to template.Parse
func handler(w http.ResponseWriter, r *http.Request) {
	userTemplate := r.FormValue("template")
	t, _ := template.New("page").Parse(userTemplate) // MATCH /potential server-side template injection: tainted data used in text/template.Parse/
	t.Execute(w, nil)
}

// Invalid: User input embedded in template string via concatenation
func renderMessage(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	tmpl := `<div>` + msg + `</div>`
	t, _ := template.New("msg").Parse(tmpl) // MATCH /potential server-side template injection: tainted data used in text/template.Parse/
	t.Execute(w, nil)
}

// Invalid: User input flows through fmt.Sprintf
func handlerSprintf(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tmpl := fmt.Sprintf(`<h1>Hello, %s!</h1>`, name)
	t, _ := template.New("greet").Parse(tmpl) // MATCH /potential server-side template injection: tainted data used in text/template.Parse/
	t.Execute(w, nil)
}

// Invalid: Taint flows through helper function
func handlerWithHelper(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("content")
	renderTemplate(w, input) // MATCH /potential server-side template injection: tainted data used in text/template.Parse/
}

func renderTemplate(w http.ResponseWriter, content string) {
	t, _ := template.New("content").Parse(content) // MATCH /potential server-side template injection: tainted data used in text/template.Parse/
	t.Execute(w, nil)
}

// Valid: Pre-defined static template with user data as parameters
var pageTemplate = template.Must(template.New("page").Parse(`
	<h1>{{.Title}}</h1>
	<p>{{.Message}}</p>
`))

func safeHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")
	pageTemplate.Execute(w, struct {
		Title   string
		Message string
	}{
		Title:   "Page",
		Message: msg,
	})
}

// Valid: Static template string (no taint)
func safeStaticTemplate(w http.ResponseWriter) {
	t, _ := template.New("static").Parse(`<h1>Hello World</h1>`)
	t.Execute(w, nil)
}
