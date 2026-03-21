# unused: parameters-are-used setting not supported

## Affected rule
No vint rule exists for unused function parameters (maps to golangci-lint `unused` with `parameters-are-used: false`)

## Behavior in golangci-lint
The `unused` linter has a `parameters-are-used` setting (default: `true`). When set to `false`, it flags function parameters that are declared but never read within the function body.

## Behavior in vint
There is no vint rule that checks for unused function parameters. The `noUnusedVariable` rule only checks package-level variables, not function parameters.

## Gap
The `parameters-are-used` setting from golangci-lint cannot be migrated. When a user has `parameters-are-used: false` in their golangci-lint config, the migrator silently drops this setting since there is no corresponding vint rule to configure.

## Example
```go
package example

// With parameters-are-used: false, golangci-lint would flag 'name' as unused.
func greet(name string) string {
    return "hello"
}
```

## Impact on migration
Users who relied on `parameters-are-used: false` to detect unused function parameters will lose that check after migration. This is a relatively uncommon configuration since the default is `true` (which means parameters are considered used and not checked).
