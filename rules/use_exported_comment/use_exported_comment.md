---
title: useExportedComment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useExportedComment`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useExportedComment:
    arguments:
      - "checkPrivateReceivers"
      - "sayRepetitiveInsteadOfStutters"
      - "checkPublicInterface"
```

## Details

Exported functions, methods, types, constants, and variables should have comments following Go documentation conventions. This rule warns on undocumented exported symbols and checks that comments are properly formatted (e.g. starting with the symbol name).

More information is available in the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments#doc-comments).

By default, the rule checks that:
- Exported functions and methods have comments starting with the function/method name.
- Exported types have comments starting with the type name (with optional leading article).
- Exported constants and variables have comments starting with the const/var name.
- Package-level names do not stutter with the package name.

Available configuration flags (passed as string arguments):
- `checkPrivateReceivers` -- enables checking public methods of private types.
- `disableStutteringCheck` -- disables checking for method names that stutter with the package name.
- `sayRepetitiveInsteadOfStutters` -- replaces the use of the term "stutters" by "is repetitive" in failure messages.
- `checkPublicInterface` -- enables checking public method definitions in public interface types.
- `disableChecksOnConstants` -- disables all checks on constant declarations.
- `disableChecksOnFunctions` -- disables all checks on function declarations.
- `disableChecksOnMethods` -- disables all checks on method declarations.
- `disableChecksOnTypes` -- disables all checks on type declarations.
- `disableChecksOnVariables` -- disables all checks on variable declarations.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package mypackage

// missing comment on exported function
func ExportedFunc() {}

// badprefix does something.
func BadPrefix() {}
```

### Valid

```golang
package mypackage

// ExportedFunc does something useful.
func ExportedFunc() {}

// MyType represents a meaningful concept.
type MyType struct{}
```
