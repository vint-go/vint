---
title: noNolintWithoutExplanation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noNolintWithoutExplanation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noNolintWithoutExplanation:
    # List of linters excluded from the explanation requirement.
    # For example, some teams may not require explanations for "lll" (line length).
    exclude:
      - "lll"
```

## Details

Reports `//nolint` directives that do not include an explanation for why the lint suppression is necessary. Each nolint directive should be followed by a comment explaining the rationale, in the format `//nolint:linter // explanation here`.

Requiring explanations serves important purposes:
- It forces developers to consciously justify why a lint rule is being suppressed, reducing careless suppressions.
- It provides context for future maintainers who may not understand why a warning was silenced.
- It makes code reviews more productive, as reviewers can evaluate whether the suppression rationale is sound.
- It discourages the practice of silencing warnings "just to make CI pass" without understanding the underlying issue.

The explanation must be a non-empty comment following the directive. A bare trailing `//` with no text after it does not satisfy the requirement.

Certain linters can be excluded from this requirement via the `excludes` configuration. For example, a team might exclude `lll` (line length linter) because the reason for suppression is usually self-evident.

This check is enabled when the `NeedsExplanation` flag is set in the nolintlint configuration (controlled by the `require-explanation` golangci-lint setting).

Source: https://github.com/golangci/golangci-lint/blob/b345eab149baac70a9613e6c7b893c5fde7c3e54/pkg/golinters/nolintlint/internal/nolintlint.go

## Examples

### Invalid

```golang
package main

//nolint:errcheck
func foo() {
    _ = doSomething()
}
```

```golang
package main

//nolint:gosec
func bar() {
    unsafeOperation()
}
```

```golang
package main

//nolint:errcheck //
func baz() {
    _ = riskyCall()
}
```

### Valid

```golang
package main

//nolint:errcheck // error return is intentionally ignored here
func foo() {
    _ = doSomething()
}
```

```golang
package main

//nolint:gosec // input is validated by the caller
func bar() {
    unsafeOperation()
}
```

```golang
package main

//nolint:errcheck,gosec // both checks are inapplicable in test code
func baz() {
    _ = riskyUncheckedCall()
}
```

```golang
package main

// With "lll" excluded from explanation requirements:
//nolint:lll
func longFunctionSignature(parameterOne string, parameterTwo string, parameterThree string) string {
    return parameterOne + parameterTwo + parameterThree
}
```
