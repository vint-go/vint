---
title: noTemplateInjection
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noTemplateInjection`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noTemplateInjection:
    # rule options here
```

## Details

Detects Server-Side Template Injection (SSTI) vulnerabilities via `text/template`.

This rule uses taint analysis to trace the flow of untrusted data from sources to `text/template` template parsing and execution. It detects cases where user-controlled input is used to construct or parse templates dynamically, allowing attackers to inject template directives that are executed on the server.

Unlike `html/template`, the `text/template` package does not perform any output escaping or sanitization. When user input is incorporated into a template string that is then parsed and executed, an attacker can inject Go template syntax to:

- Read internal data by accessing exported fields and methods on the template data context
- Call arbitrary methods on objects passed to the template
- Cause denial of service through resource-intensive template operations

User input should never be used to construct template strings. If dynamic templating is required, use pre-defined templates with user data passed as template parameters.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "text/template"

func handler(w http.ResponseWriter, r *http.Request) {
    userTemplate := r.FormValue("template")
    // User input used as template source - SSTI
    t, _ := template.New("page").Parse(userTemplate)
    t.Execute(w, data)
}
```

```golang
import "text/template"

func renderMessage(w http.ResponseWriter, msg string) {
    // User input embedded in template string
    tmpl := `<div>` + msg + `</div>`
    t, _ := template.New("msg").Parse(tmpl)
    t.Execute(w, nil)
}
```

### Valid

```golang
import "html/template"

// Pre-defined template with user data as parameters
var pageTemplate = template.Must(template.New("page").Parse(`
    <h1>{{.Title}}</h1>
    <p>{{.Message}}</p>
`))

func handler(w http.ResponseWriter, r *http.Request) {
    msg := r.FormValue("message")
    // User input passed as data, not as template source
    pageTemplate.Execute(w, struct {
        Title   string
        Message string
    }{
        Title:   "Page",
        Message: msg,
    })
}
```

```golang
import "html/template"

// Using html/template instead of text/template for HTML output
func handler(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    t := template.Must(template.New("greeting").Parse("<h1>Hello, {{.}}!</h1>"))
    t.Execute(w, name) // html/template auto-escapes the output
}
```
