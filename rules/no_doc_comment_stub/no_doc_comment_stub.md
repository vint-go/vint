---
title: noDocCommentStub
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noDocCommentStub`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noDocCommentStub:
    # no additional options
```

## Details

Detects comments that silence go lint complaints about doc-comment. This checker identifies documentation comment stubs on exported functions and types where developers write minimal placeholder documentation (like `...`, `.`, `xxx`, or `whatever`) to suppress linter warnings rather than providing meaningful documentation.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// ProcessData ...
func ProcessData() {}
```

```golang
// MyType whatever
type MyType struct{}
```

### Valid

```golang
// ProcessData reads input data and transforms it into the output format.
func ProcessData() {}
```

```golang
// MyType represents a domain entity with associated metadata.
type MyType struct{}
```
