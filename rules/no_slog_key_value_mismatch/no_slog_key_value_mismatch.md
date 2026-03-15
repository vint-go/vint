---
title: noSlogKeyValueMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSlogKeyValueMismatch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSlogKeyValueMismatch:
    # rule options here
```

## Details

Checks for mismatched key-value pairs in `log/slog` calls. The `slog` package uses alternating key-value pairs for structured logging. This analyzer detects:

- Odd number of arguments (missing a value for a key)
- Non-string keys where a string key is expected
- Mismatched pairs that would result in malformed log entries

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/slog

## Examples

### Invalid

```golang
import "log/slog"

func example() {
    // Bad: odd number of arguments, "status" has no value
    slog.Info("request", "method", "GET", "status")
}
```

```golang
import "log/slog"

func example() {
    // Bad: key should be a string, not an integer
    slog.Info("request", 200, "OK")
}
```

### Valid

```golang
import "log/slog"

func example() {
    // Good: properly paired key-value arguments
    slog.Info("request", "method", "GET", "status", 200)
}
```

```golang
import "log/slog"

func example() {
    // Good: using slog.Attr for typed attributes
    slog.Info("request",
        slog.String("method", "GET"),
        slog.Int("status", 200),
    )
}
```
