---
title: noUnrecognizedDirective
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnrecognizedDirective`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnrecognizedDirective:
    # no configurable options
```

## Details

Detects Go compiler directives that use an unrecognized directive name. For example, `//go:embod` instead of `//go:embed`, or `//go:linknme` instead of `//go:linkname`.

In Go, compiler directives are special comments that begin with `//go:` followed by a recognized directive name. If the directive name is misspelled or is not part of the known set of compiler directives, the Go compiler silently ignores it without producing any error. This can lead to subtle bugs where the intended compiler behavior (such as embedding files, generating code, or applying optimizations) does not occur.

This rule validates that any `//go:` comment uses a directive name from the following known set:

- `build` -- Build constraints
- `embed` -- File embedding
- `generate` -- Code generation
- `linkname` -- Symbol linking
- `noinline` -- Disable inlining
- `nosplit` -- Disable stack split check
- `noescape` -- Disable escape analysis
- `norace` -- Disable race detection
- `nocheckptr` -- Disable pointer checking
- `nointerface` -- Disable interface method
- `nowritebarrier` -- Disable write barrier
- `nowritebarrierrec` -- Disable recursive write barrier
- `yeswritebarrierrec` -- Enable recursive write barrier
- `systemstack` -- Require system stack
- `uintptrescapes` -- Mark uintptr arguments as escaping
- `uintptrkeepalive` -- Keep uintptr arguments alive
- `notinheap` -- Mark type as not in heap
- `cgo_dynamic_linker` -- Set dynamic linker for cgo
- `cgo_export_dynamic` -- Export symbol dynamically for cgo
- `cgo_export_static` -- Export symbol statically for cgo
- `cgo_import_dynamic` -- Import dynamic symbol for cgo
- `cgo_import_static` -- Import static symbol for cgo
- `cgo_ldflag` -- Set linker flag for cgo
- `cgo_unsafe_args` -- Mark cgo arguments as unsafe
- `debug` -- Debug directive
- `wasmimport` -- WebAssembly import
- `wasmexport` -- WebAssembly export

The diagnostic message produced is: `compiler directive unrecognized: <directive>`.

Source: https://github.com/leighmcculloch/gocheckcompilerdirectives

## Examples

### Invalid

```golang
//go:embod hello.txt
var s string
```

```golang
//go:genrate stringer -type=Pill
package painkiller
```

```golang
//go:linknme localname importpath.name
func localname()
```

```golang
//go:noInline
func Add(a, b int) int {
	return a + b
}
```

```golang
//go:inline
func Multiply(a, b int) int {
	return a * b
}
```

### Valid

```golang
//go:embed hello.txt
var s string
```

```golang
//go:generate stringer -type=Pill
package painkiller
```

```golang
//go:linkname localname importpath.name
func localname()
```

```golang
//go:noinline
func Add(a, b int) int {
	return a + b
}
```

```golang
//go:build linux
```

```golang
//go:nosplit
func criticalFunc() {
	// ...
}
```
