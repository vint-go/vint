package fixtures

import (
	"html/template"
	"net/http"
)

// Invalid: Bypassing HTML escaping with user input
func renderPage(w http.ResponseWriter, userInput string) {
	t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
	data := struct {
		Content template.HTML
	}{
		Content: template.HTML(userInput), // MATCH /use of unescaped data in html/template type conversion template.HTML/
	}
	t.Execute(w, data)
}

// Invalid: Bypassing JS escaping with user input
func renderScript(w http.ResponseWriter, callback string) {
	t := template.Must(template.New("page").Parse(`<script>{{.Code}}</script>`))
	data := struct {
		Code template.JS
	}{
		Code: template.JS(callback), // MATCH /use of unescaped data in html/template type conversion template.JS/
	}
	t.Execute(w, data)
}

// Invalid: Bypassing URL escaping with user input
func renderLink(w http.ResponseWriter, href string) {
	t := template.Must(template.New("page").Parse(`<a href="{{.Link}}">click</a>`))
	data := struct {
		Link template.URL
	}{
		Link: template.URL(href), // MATCH /use of unescaped data in html/template type conversion template.URL/
	}
	t.Execute(w, data)
}

// Invalid: Bypassing CSS escaping with user input
func renderStyle(w http.ResponseWriter, style string) {
	t := template.Must(template.New("page").Parse(`<div style="{{.Style}}">content</div>`))
	data := struct {
		Style template.CSS
	}{
		Style: template.CSS(style), // MATCH /use of unescaped data in html/template type conversion template.CSS/
	}
	t.Execute(w, data)
}

// Valid: Using string literals (developer-controlled constants are safe)
func renderStatic(w http.ResponseWriter) {
	t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
	data := struct {
		Content template.HTML
	}{
		Content: template.HTML("<b>safe static content</b>"),
	}
	t.Execute(w, data)
}

// Valid: Using automatic escaping with string type
func renderSafe(w http.ResponseWriter, userInput string) {
	t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
	data := struct {
		Content string
	}{
		Content: userInput,
	}
	t.Execute(w, data)
}
