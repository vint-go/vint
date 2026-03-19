---
title: noRedundantBuildTag
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantBuildTag`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantBuildTag:
    # rule options here
```

## Details

This rule warns about redundant build tag comments `// +build` when `//go:build` is present. `gofmt` in Go 1.17+ automatically adds the `//go:build` constraint, making the `// +build` comment unnecessary. Keeping both forms adds visual noise and can become a maintenance burden if one is updated without the other.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
//go:build tag
// +build tag

package pkg
```

### Valid

```golang
//go:build tag

package pkg
```

```golang
// +build tag

// Pre-Go 1.17 file without //go:build directive
package pkg
```
