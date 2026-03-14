---
title: noLogInjectionTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noLogInjectionTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noLogInjectionTaint:
    # rule options here
```

## Details

Detects log injection vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources to logging sinks. It performs interprocedural data flow analysis to detect cases where user-controlled input is written to log messages without sanitization.

Log injection occurs when an attacker injects specially crafted input that, when logged, can forge log entries, corrupt log files, or exploit log processing tools. By injecting newline characters, an attacker can create fake log entries that may mislead forensic analysis. In some cases, log injection can also exploit log management systems that process log entries as commands or structured data.

User input included in log messages should be sanitized to remove or escape control characters (especially newlines, carriage returns, and null bytes). Structured logging formats (JSON) can also help mitigate log injection by properly encoding field values.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "log"

func handler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    // Taint flows from request to log output
    log.Printf("Login attempt for user: %s", username)
    // Attacker could inject: "admin\n2024-01-01 INFO Login successful for user: admin"
}
```

```golang
import "log/slog"

func handler(w http.ResponseWriter, r *http.Request) {
    input := r.URL.Query().Get("search")
    // Unsanitized input in log message
    slog.Info("Search performed", "query", input)
}
```

### Valid

```golang
import (
    "log"
    "strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    // Sanitize control characters before logging
    clean := strings.NewReplacer("\n", "", "\r", "", "\x00", "").Replace(username)
    log.Printf("Login attempt for user: %s", clean)
}
```

```golang
import (
    "encoding/json"
    "log/slog"
)

func handler(w http.ResponseWriter, r *http.Request) {
    input := r.URL.Query().Get("search")
    // Structured logging with proper encoding
    slog.Info("Search performed",
        slog.String("query", input),
        slog.String("remote_addr", r.RemoteAddr),
    )
}
```
