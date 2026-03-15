---
title: noNolintNonMachineReadable
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noNolintNonMachineReadable`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noNolintNonMachineReadable:
    # No additional options for this rule.
```

## Details

Reports `//nolint` directives that are not written in a machine-readable format. The machine-readable format requires no space between the comment marker `//` and the `nolint` keyword. Directives like `// nolint` should be rewritten as `//nolint`.

When a space exists between `//` and `nolint`, automated tools and golangci-lint itself may not properly recognize the directive, causing lint suppressions to be silently ignored. This can lead to confusing situations where developers think they have suppressed a lint warning, but it still fires.

The linter provides an automatic fix that removes the leading space to produce the correct machine-readable format.

This check is enabled when the `NeedsMachineOnly` flag is set in the nolintlint configuration (controlled by the golangci-lint setting for machine-readable nolint directives).

Source: https://github.com/golangci/golangci-lint/blob/b345eab149baac70a9613e6c7b893c5fde7c3e54/pkg/golinters/nolintlint/internal/nolintlint.go

## Examples

### Invalid

```golang
package main

// nolint:errcheck
func foo() {
    _ = doSomething()
}
```

```golang
package main

// nolint:gosec // reason for suppression
func bar() {
    unsafeOperation()
}
```

```golang
package main

// nolint
func baz() {
    _ = riskyCall()
}
```

### Valid

```golang
package main

//nolint:errcheck
func foo() {
    _ = doSomething()
}
```

```golang
package main

//nolint:gosec // this is safe in this context
func bar() {
    unsafeOperation()
}
```

```golang
package main

//nolint:errcheck,gosec // multiple linters suppressed for valid reason
func baz() {
    _ = riskyUncheckedCall()
}
```
