---
title: noDuplicateBuildConstraint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDuplicateBuildConstraint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDuplicateBuildConstraint:
    # rule options here
```

## Details

Multiple, identical build constraints in the same file.

Having duplicate build constraints (e.g., two identical `//go:build` lines) is redundant and likely a copy-paste mistake.

Source: https://staticcheck.dev/docs/checks/#SA4019

## Examples

### Invalid

```golang
//go:build linux
//go:build linux

package main
```

### Valid

```golang
//go:build linux

package main
```
