---
title: noExposedPprof
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noExposedPprof`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noExposedPprof:
    # rule options here
```

## Details

Detects when the profiling endpoint (`net/http/pprof`) is automatically exposed on an HTTP server.

Importing `net/http/pprof` as a side-effect import automatically registers profiling handlers on the default HTTP mux (`http.DefaultServeMux`). If the default mux is used for a public-facing server, profiling data such as memory allocations, goroutine stacks, CPU profiles, and heap dumps are exposed to anyone who can reach the server.

This information can reveal sensitive details about the application's internal state, runtime behavior, and memory contents. Profiling endpoints should only be exposed on internal or debug servers that are not accessible from untrusted networks.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "net/http"
    _ "net/http/pprof" // Automatically exposes profiling endpoints
)

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil) // pprof endpoints are now exposed
}
```

### Valid

```golang
import (
    "net/http"
    "net/http/pprof"
)

func main() {
    // Public server on its own mux
    publicMux := http.NewServeMux()
    publicMux.HandleFunc("/", handler)
    go http.ListenAndServe(":8080", publicMux)

    // Debug server with pprof on separate port, internal only
    debugMux := http.NewServeMux()
    debugMux.HandleFunc("/debug/pprof/", pprof.Index)
    http.ListenAndServe("127.0.0.1:6060", debugMux)
}
```
