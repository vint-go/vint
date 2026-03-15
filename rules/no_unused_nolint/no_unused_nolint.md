---
title: noUnusedNolint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnusedNolint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnusedNolint:
    # No additional options for this rule.
```

## Details

Reports `//nolint` directives that are unnecessary because they do not suppress any actual linter warnings. An unused nolint directive indicates one of several situations:

- The underlying code issue that originally triggered the lint warning has been fixed, but the suppression comment was not removed.
- The nolint directive was added preemptively but no actual warning exists for the specified linter on that line.
- The linter name in the directive is misspelled or does not match any enabled linter.
- The directive suppresses a specific linter, but the warning on that line comes from a different linter.

Unused nolint directives add noise to the codebase and can be confusing to maintainers. They suggest that a lint issue exists when it does not, or that the code was previously problematic in some way. Removing them keeps the codebase clean and makes the remaining nolint directives more meaningful.

The linter provides an automatic fix that removes the unused nolint directive entirely. If the directive suppresses multiple linters and only some are unused, the fix removes only the unused linter names from the directive.

This check is enabled when the `NeedsUnused` flag is set (controlled by the `allow-unused` golangci-lint setting being set to `false`).

Source: https://github.com/golangci/golangci-lint/blob/b345eab149baac70a9613e6c7b893c5fde7c3e54/pkg/golinters/nolintlint/internal/nolintlint.go

## Examples

### Invalid

```golang
package main

//nolint:errcheck // no errcheck warning exists on this line
func foo() int {
    return 42
}
```

```golang
package main

//nolint:gosec // gosec does not flag this code
func bar() string {
    return "safe string"
}
```

```golang
package main

//nolint // no linter warnings exist here at all
func baz() {
    fmt.Println("hello")
}
```

### Valid

```golang
package main

//nolint:errcheck // error is intentionally ignored
func foo() {
    _ = doSomething() // doSomething() returns an error that errcheck would flag
}
```

```golang
package main

// No nolint directive needed when there are no warnings
func bar() int {
    return 42
}
```

```golang
package main

//nolint:gosec // G104: unhandled error is acceptable here
func baz() {
    cmd := exec.Command("ls") // gosec flags exec.Command usage
    cmd.Run()
}
```
