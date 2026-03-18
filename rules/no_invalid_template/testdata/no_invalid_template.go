package fixtures

import "text/template"

func invalidTemplates() {
	// Invalid template - unclosed action
	t := template.Must(template.New("test").Parse("{{.Name")) // MATCH /invalid template: template: :1: unclosed action/
	_ = t

	// Invalid template - unclosed action with pipe
	t2 := template.Must(template.New("test2").Parse("{{.Name | ")) // MATCH /invalid template: template: :1: unclosed action/
	_ = t2

	// Invalid template via direct Parse call
	t3, _ := template.New("test3").Parse("{{if}}") // MATCH /invalid template: template: :1: missing value for if/
	_ = t3
}

func validTemplates() {
	// Valid template
	t := template.Must(template.New("test").Parse("{{.Name}}"))
	_ = t

	// Valid template - plain text (no actions)
	t3 := template.Must(template.New("test3").Parse("Hello, world!"))
	_ = t3

	// Valid template - empty string
	t4 := template.Must(template.New("test4").Parse(""))
	_ = t4

	// Valid template with if/else
	t5 := template.Must(template.New("test5").Parse("{{if .Name}}Hello{{else}}Bye{{end}}"))
	_ = t5

	// Non-literal template (should be skipped)
	tmplStr := "{{.Name"
	t6 := template.Must(template.New("test6").Parse(tmplStr))
	_ = t6
}
