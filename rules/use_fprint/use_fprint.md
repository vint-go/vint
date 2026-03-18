---
title: useFprint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useFprint`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useFprint:
    # no additional options
```

## Details

Detects `fmt.Sprint`/`fmt.Sprintf` calls that can be replaced with `fmt.Fprint`/`fmt.Fprintf` variants. When the result of `fmt.Sprint` is immediately written to an `io.Writer`, it is more efficient to use `fmt.Fprint` directly, avoiding the intermediate string allocation.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
w.Write([]byte(fmt.Sprintf("hello %s", name)))
```

### Valid

```golang
fmt.Fprintf(w, "hello %s", name)
```
