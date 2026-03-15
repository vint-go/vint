---
title: noNolintWithoutSpecificLinter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noNolintWithoutSpecificLinter`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noNolintWithoutSpecificLinter:
    # No additional options for this rule.
```

## Details

Reports `//nolint` directives that do not specify which linter(s) are being suppressed. A bare `//nolint` directive suppresses all linters for the annotated line, which is overly broad and can mask real issues. Instead, the directive should explicitly name the linter(s) being suppressed, such as `//nolint:errcheck` or `//nolint:gosec,errcheck`.

Being specific about which linters are suppressed serves several purposes:
- It documents the intent of the suppression, making it clear which specific warning is being silenced.
- It prevents accidentally suppressing unrelated warnings that may be introduced later.
- It makes code reviews easier because reviewers can assess whether the suppression is justified for the specific linter.

This check is enabled when the `NeedsSpecific` flag is set in the nolintlint configuration (controlled by the `require-specific` golangci-lint setting).

Source: https://github.com/golangci/golangci-lint/blob/b345eab149baac70a9613e6c7b893c5fde7c3e54/pkg/golinters/nolintlint/internal/nolintlint.go

## Examples

### Invalid

```golang
package main

//nolint
func foo() {
    _ = doSomething()
}
```

```golang
package main

//nolint // suppress all linters
func bar() {
    unsafeOperation()
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

//nolint:gosec // this input is validated upstream
func bar() {
    unsafeOperation()
}
```

```golang
package main

//nolint:errcheck,gosec // both are acceptable here
func baz() {
    _ = riskyUncheckedCall()
}
```
