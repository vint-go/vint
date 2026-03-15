---
title: noXssTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noXssTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noXssTaint:
    # rule options here
```

## Details

Detects Cross-Site Scripting (XSS) vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from HTTP request sources to HTTP response sinks. It performs interprocedural data flow analysis to detect cases where user-controlled input is written to HTTP responses without proper encoding or sanitization.

XSS attacks allow an attacker to inject malicious scripts into web pages viewed by other users. This can lead to session hijacking, credential theft, defacement, or redirection to malicious sites. There are three types of XSS: reflected (input immediately returned in response), stored (input persisted and later displayed), and DOM-based (client-side script manipulation).

User input must be properly encoded for the output context (HTML, JavaScript, CSS, URL) before being included in HTTP responses. Using Go's `html/template` package provides automatic contextual escaping.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func handler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    // Taint flows from request to response - reflected XSS
    fmt.Fprintf(w, "<h1>Hello, %s!</h1>", name)
}
```

```golang
func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("q")
    // User input reflected in HTML response
    w.Write([]byte("<p>Results for: " + query + "</p>"))
}
```

### Valid

```golang
import "html/template"

func handler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    // Using html/template for automatic escaping
    t := template.Must(template.New("page").Parse("<h1>Hello, {{.}}!</h1>"))
    t.Execute(w, name)
}
```

```golang
import "html"

func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("q")
    // Manually escaping user input
    escaped := html.EscapeString(query)
    w.Write([]byte("<p>Results for: " + escaped + "</p>"))
}
```

```golang
import "encoding/json"

func apiHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    // JSON encoding for API responses
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"name": name})
}
```
