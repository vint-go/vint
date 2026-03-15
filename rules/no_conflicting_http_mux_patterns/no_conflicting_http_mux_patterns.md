---
title: noConflictingHttpMuxPatterns
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noConflictingHttpMuxPatterns`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noConflictingHttpMuxPatterns:
    # rule options here
```

## Details

Analyzer for HTTP multiplexer patterns. This analyzer checks for issues related to HTTP request multiplexer usage, such as patterns that may conflict or overlap when using Go 1.22+ enhanced ServeMux routing patterns.

Go 1.22 introduced method-based routing and wildcard patterns in `http.ServeMux`. This analyzer helps detect potential issues with the new routing syntax.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/httpmux

## Examples

### Invalid

```golang
import "net/http"

func setup() {
    mux := http.NewServeMux()
    // Potentially conflicting patterns
    mux.HandleFunc("/items/{id}", handleItem)
    mux.HandleFunc("/items/{name}", handleItemByName) // conflicts with above
}
```

### Valid

```golang
import "net/http"

func setup() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /items/{id}", handleItem)
    mux.HandleFunc("POST /items", createItem)
}
```
