# errorlint: missing useErrorsAs rule

## Affected rule
`useErrorsAs` (maps to golangci-lint `errorlint` with `asserts: true`)

## Behavior in golangci-lint
When errorlint is enabled with `asserts: true` (the default), it detects type assertions and type switches on error values that should use `errors.As()` instead. Type assertions like `err.(*MyError)` fail when errors are wrapped because the wrapper type doesn't match the inner error type.

Example flagged code:
```go
myErr, ok := err.(*MyError)        // should use errors.As
switch e := err.(type) { ... }     // should use errors.As
```

## Behavior in vint
There is no implemented vint rule for `useErrorsAs`. The rule is registered in the rules registry but has no Go implementation. Note: vint has `noInvalidErrorsAs` which checks for invalid *arguments* to `errors.As()`, but that is a different check -- it does not detect code that should use `errors.As()` instead of type assertions.

## Gap
The entire type assertion/switch check on error values is missing from vint. Users who rely on errorlint to catch type assertions on errors will lose that coverage after migration.

## Example
```go
package example

import "fmt"

type MyError struct{ Code int }
func (e *MyError) Error() string { return fmt.Sprintf("code %d", e.Code) }

func handle(err error) {
    // golangci-lint errorlint flags this; vint does not
    if myErr, ok := err.(*MyError); ok {
        fmt.Println(myErr.Code)
    }
}
```

## Impact on migration
- Config migration correctly omits this rule (since it doesn't exist), so no broken config is generated.
- `//nolint:errorlint` directives that suppress type assertion findings cannot be converted to specific vint-ignore directives for this check.
- Users lose type assertion error checking coverage after migrating from golangci-lint to vint.
