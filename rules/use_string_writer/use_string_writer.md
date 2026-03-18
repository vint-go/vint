---
title: useStringWriter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useStringWriter`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useStringWriter:
    # no additional options
```

## Details

Detects `Write`/`WriteString` calls that can be replaced with the `io.StringWriter` interface's `WriteString` method. When writing a string to a writer, using `WriteString` directly is more efficient than converting to `[]byte` first.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
w.Write([]byte("hello world"))
```

### Valid

```golang
io.WriteString(w, "hello world")
```
