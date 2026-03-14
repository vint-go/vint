---
title: noFilesystemRootServing
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noFilesystemRootServing`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noFilesystemRootServing:
    # rule options here
```

## Details

Detects the use of `http.Dir("/")` as a potential directory traversal risk.

Using `http.Dir("/")` as a file server root exposes the entire filesystem through the HTTP server. An attacker could access sensitive system files such as `/etc/passwd`, `/etc/shadow`, configuration files, and other critical system data.

Even when the file server is intended to serve a specific directory, using the filesystem root as the base allows path traversal attacks. File servers should be rooted at the most restrictive directory possible and should use `http.Dir` with a specific, limited directory path.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// Serving the entire filesystem root
func main() {
    http.Handle("/", http.FileServer(http.Dir("/")))
    http.ListenAndServe(":8080", nil)
}
```

### Valid

```golang
import "net/http"

// Serving a specific, restricted directory
func main() {
    http.Handle("/static/", http.StripPrefix("/static/",
        http.FileServer(http.Dir("./public"))))
    http.ListenAndServe(":8080", nil)
}
```

```golang
import (
    "io/fs"
    "net/http"
)

// Using embedded filesystem
//go:embed static
var staticFiles embed.FS

func main() {
    http.Handle("/static/", http.FileServer(http.FS(staticFiles)))
    http.ListenAndServe(":8080", nil)
}
```
