---
title: noCgiImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noCgiImport`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noCgiImport:
    # rule options here
```

## Details

Detects the import of the `net/http/cgi` package, which is blocklisted.

The `net/http/cgi` package implements the Common Gateway Interface (CGI) protocol, which has known security issues. CGI has historically been a source of vulnerabilities including header injection, environment variable manipulation, and information disclosure. The package has had specific security vulnerabilities (e.g., CVE-2020-24553 for cross-site scripting).

Modern Go applications should use `net/http` handlers directly or use `net/http/fcgi` (FastCGI) if CGI-like functionality is needed, as it provides better isolation.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http/cgi" // Blocklisted import

func main() {
    handler := &cgi.Handler{
        Path: "/usr/bin/script",
    }
    http.Handle("/cgi-bin/", handler)
    http.ListenAndServe(":8080", nil)
}
```

### Valid

```golang
import "net/http"

func main() {
    http.HandleFunc("/api/", apiHandler)
    http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    // Handle request directly in Go
    w.Write([]byte("Hello"))
}
```
