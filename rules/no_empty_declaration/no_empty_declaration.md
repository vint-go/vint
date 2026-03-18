---
title: noEmptyDeclaration
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noEmptyDeclaration`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noEmptyDeclaration:
    # no additional options
```

## Details

Detects empty `var`, `const`, `type`, or `import` declaration blocks that have no content. These empty blocks add noise to the code and should be removed. They typically appear as leftover artifacts from refactoring or code generation.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var ()
```

```golang
const ()
```

```golang
type ()
```

### Valid

```golang
var x int
```

```golang
// Remove the empty declaration entirely
```
