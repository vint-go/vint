---
title: noMalformedBuildTag
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMalformedBuildTag`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMalformedBuildTag:
    # rule options here
```

## Details

Checks build tags for correctness. This analyzer verifies that `//go:build` and `// +build` directives are well-formed and consistent. It detects issues such as:

- Malformed `//go:build` constraints
- Mismatched `//go:build` and `// +build` lines (in Go versions that support both)
- Build tags placed in the wrong location in the file (must appear before the package clause)
- Invalid syntax in build constraint expressions

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/buildtag

## Examples

### Invalid

```golang
//go:build linux AND darwin
// Invalid: should use && or , not AND

package example
```

```golang
package example

//go:build linux
// Invalid: build tag must appear before the package clause
```

### Valid

```golang
//go:build linux && amd64

package example
```

```golang
//go:build !windows

package example
```
