---
title: noExternalErrorReassign
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noExternalErrorReassign`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noExternalErrorReassign:
    # no additional options
```

## Details

Detects suspicious reassignment of exported error variables from other packages. Package-level error sentinel values (like `io.EOF` or `http.ErrServerClosed`) should be compared against, not reassigned. Reassigning these variables modifies global state, which can cause subtle bugs across the entire program since any code depending on the original error value will break silently.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
io.EOF = nil // reassigning a package-level error sentinel
```

```golang
http.ErrServerClosed = errors.New("custom error")
```

### Valid

```golang
if err == io.EOF {
    // compare, don't reassign
}
```
