---
title: noSpaceInDirective
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSpaceInDirective`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSpaceInDirective:
    # no configurable options
```

## Details

Detects Go compiler directives that contain a space between the comment slashes (`//`) and the `go:` prefix. For example, `// go:embed` instead of `//go:embed`.

In Go, compiler directives such as `//go:generate`, `//go:embed`, `//go:build`, `//go:linkname`, and others must be written without any space between `//` and the directive name. If a space is present (e.g., `// go:embed`), the Go compiler silently ignores the directive without producing any error. This can lead to subtle and hard-to-diagnose bugs where directives simply do not take effect.

This rule catches such mistakes by scanning all comments in Go source files and flagging any that match the pattern `// go:` (with a space), which indicates the developer likely intended to write a compiler directive but accidentally introduced whitespace.

The diagnostic message produced is: `compiler directive contains space: <directive>`.

Source: https://github.com/leighmcculloch/gocheckcompilerdirectives

## Examples

### Invalid

```golang
// go:generate stringer -type=Pill
package painkiller
```

```golang
package main

// go:embed hello.txt
var s string
```

```golang
package main

// go:noinline
func Add(a, b int) int {
	return a + b
}
```

```golang
package main

// go:build linux
```

```golang
package main

// go:linkname localname importpath.name
func localname()
```

### Valid

```golang
//go:generate stringer -type=Pill
package painkiller
```

```golang
package main

//go:embed hello.txt
var s string
```

```golang
package main

//go:noinline
func Add(a, b int) int {
	return a + b
}
```

```golang
package main

//go:build linux
```

```golang
package main

//go:linkname localname importpath.name
func localname()
```
