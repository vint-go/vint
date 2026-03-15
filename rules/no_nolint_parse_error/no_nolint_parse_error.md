---
title: noNolintParseError
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNolintParseError`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNolintParseError:
    # No additional options for this rule.
```

## Details

Reports `//nolint` directives that have malformed syntax and cannot be properly parsed. A valid nolint directive must match the pattern `//nolint[:<comma-separated-linters>] [// explanation]`.

Common parse errors include:
- Using spaces instead of commas to separate linter names (e.g., `//nolint:errcheck gosec` instead of `//nolint:errcheck,gosec`).
- Missing the colon before linter names (e.g., `//nolint errcheck`).
- Using invalid characters in linter names.
- Having trailing content that does not follow the expected format.

When the directive cannot be parsed, golangci-lint may silently ignore it, meaning the intended suppression does not take effect. This rule ensures that all nolint directives are syntactically correct so they work as intended.

This check is always enabled regardless of configuration flags, as it catches fundamental syntax issues.

Source: https://github.com/golangci/golangci-lint/blob/b345eab149baac70a9613e6c7b893c5fde7c3e54/pkg/golinters/nolintlint/internal/nolintlint.go

## Examples

### Invalid

```golang
package main

//nolint:errcheck gosec
func foo() {
    _ = doSomething()
}
```

```golang
package main

//nolint: errcheck, ,gosec
func bar() {
    _ = anotherCall()
}
```

```golang
package main

//nolint:errcheck:gosec
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

//nolint:errcheck,gosec
func bar() {
    _ = anotherCall()
}
```

```golang
package main

//nolint:errcheck,gosec // both are safe in this context
func baz() {
    _ = riskyCall()
}
```
