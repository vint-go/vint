---
title: noUnescapedHtmlTemplate
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnescapedHtmlTemplate`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnescapedHtmlTemplate:
    # rule options here
```

## Details

Detects the use of unescaped data in HTML templates, which can lead to cross-site scripting (XSS) vulnerabilities.

This rule identifies cases where `template.HTML`, `template.JS`, `template.URL`, or `template.CSS` type conversions are used with variable data. These types bypass Go's `html/template` package's automatic contextual escaping, marking the content as safe for direct inclusion in HTML output. When user-controlled data is cast to these types, it can contain malicious JavaScript, CSS, or URL payloads that will be rendered without sanitization.

The `html/template` package provides automatic escaping by default. Using these unescaped types should only be done with trusted, pre-sanitized content.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "html/template"

func renderPage(w http.ResponseWriter, userInput string) {
    t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
    // Bypassing HTML escaping with user input
    data := struct {
        Content template.HTML
    }{
        Content: template.HTML(userInput), // XSS vulnerability
    }
    t.Execute(w, data)
}
```

```golang
import "html/template"

func renderScript(w http.ResponseWriter, callback string) {
    t := template.Must(template.New("page").Parse(`<script>{{.Code}}</script>`))
    data := struct {
        Code template.JS
    }{
        Code: template.JS(callback), // JavaScript injection
    }
    t.Execute(w, data)
}
```

### Valid

```golang
import "html/template"

func renderPage(w http.ResponseWriter, userInput string) {
    // html/template automatically escapes string data
    t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
    data := struct {
        Content string
    }{
        Content: userInput, // Automatically escaped
    }
    t.Execute(w, data)
}
```

```golang
import (
    "html/template"
    "github.com/microcosm-cc/bluemonday"
)

func renderPage(w http.ResponseWriter, userHTML string) {
    // Sanitize HTML before marking as safe
    p := bluemonday.UGCPolicy()
    sanitized := p.Sanitize(userHTML)
    t := template.Must(template.New("page").Parse(`<div>{{.Content}}</div>`))
    data := struct {
        Content template.HTML
    }{
        Content: template.HTML(sanitized), // Sanitized before casting
    }
    t.Execute(w, data)
}
```
