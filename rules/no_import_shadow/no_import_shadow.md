---
title: noImportShadow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noImportShadow`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noImportShadow:
    # no additional options
```

## Details

Detects when imported package names are shadowed in assignments. This checker identifies cases where a variable assignment uses the same name as an imported package, which shadows the package import and makes it inaccessible within that scope. This can lead to confusion and subtle bugs.

This rule is distinct from `lint/suspicious/noVariableShadowing`, which detects general variable shadowing via `:=` in inner scopes, and from `lint/suspicious/noBuiltinShadow` / `lint/suspicious/noBuiltinShadowDecl`, which detect shadowing of Go's predeclared identifiers. This rule specifically targets shadowing of imported package names.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
import "path/filepath"

func f() {
    filepath := "/some/path" // shadows the filepath package
    _ = filepath
}
```

### Valid

```golang
import "path/filepath"

func f() {
    fp := "/some/path"
    fullPath := filepath.Join(fp, "file.txt")
    _ = fullPath
}
```
