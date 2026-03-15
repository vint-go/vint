---
title: noUnboundedFormParsing
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnboundedFormParsing`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnboundedFormParsing:
    # rule options here
```

## Details

Detects unbounded form parsing in HTTP handlers that can cause memory exhaustion.

This rule uses SSA analysis to identify calls to `http.Request.ParseForm`, `http.Request.ParseMultipartForm`, or `http.Request.FormValue` without properly limiting the maximum size of the request body. An attacker can send extremely large form data or multipart uploads that consume all available server memory, causing a denial-of-service condition.

HTTP handlers should use `http.MaxBytesReader` to wrap the request body before parsing forms, or use `ParseMultipartForm` with an explicit `maxMemory` parameter set to a reasonable limit.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // No size limit on form parsing
    r.ParseMultipartForm(0)
    file, _, err := r.FormFile("upload")
    // ...
}
```

```golang
import "net/http"

func formHandler(w http.ResponseWriter, r *http.Request) {
    // Unbounded form parsing
    r.ParseForm()
    value := r.FormValue("data")
    // ...
}
```

### Valid

```golang
import "net/http"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Limit request body size to 10MB
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
    r.ParseMultipartForm(10 << 20)
    file, _, err := r.FormFile("upload")
    // ...
}
```

```golang
import "net/http"

func formHandler(w http.ResponseWriter, r *http.Request) {
    // Limit request body size
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
    r.ParseForm()
    value := r.FormValue("data")
    // ...
}
```
