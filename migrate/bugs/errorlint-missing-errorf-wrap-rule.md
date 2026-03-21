# errorlint: missing errorf wrapping rule

## Affected rule
`noNonWrappingErrorf` (maps to golangci-lint `errorlint` with `errorf: true`)

## Behavior in golangci-lint
When errorlint is enabled with `errorf: true` (the default), it detects `fmt.Errorf` calls that use non-wrapping format verbs (`%v`, `%s`, etc.) instead of `%w` for error arguments. This ensures errors are wrapped properly to preserve the error chain for `errors.Is` and `errors.As`.

Example flagged code:
```go
fmt.Errorf("failed: %v", err)  // should use %w
```

## Behavior in vint
There is no implemented vint rule for `noNonWrappingErrorf`. The rule is registered in the rules registry but has no Go implementation.

## Gap
The entire errorf wrapping check is missing from vint. Users who rely on errorlint to catch non-wrapping `fmt.Errorf` calls will lose that coverage after migration.

## Example
```go
package example

import "fmt"

func wrap(err error) error {
    // golangci-lint errorlint flags this; vint does not
    return fmt.Errorf("operation failed: %v", err)
}
```

## Impact on migration
- Config migration correctly omits this rule (since it doesn't exist), so no broken config is generated.
- `//nolint:errorlint` directives that suppress errorf wrapping findings cannot be converted to specific vint-ignore directives for this check.
- Users lose errorf wrapping coverage after migrating from golangci-lint to vint.
