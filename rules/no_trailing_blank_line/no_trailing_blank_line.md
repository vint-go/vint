---
title: noTrailingBlankLine
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noTrailingBlankLine`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noTrailingBlankLine:
    # No additional options. This check is enabled by default.
```

## Details

Detects unnecessary blank lines at the end of a function body, if/else block, for loop, switch statement, or case clause. Trailing blank lines inside blocks add visual noise, waste vertical screen space, and create inconsistent formatting across a codebase.

The whitespace linter automatically removes these blank lines when a fix is applied, bringing the closing brace directly after the last statement in the block.

This check runs by default and does not need to be explicitly enabled.

Source: https://github.com/ultraware/whitespace

## Examples

### Invalid

```golang
// Unnecessary blank line at the end of a function body.
func greet(name string) string {
	greeting := "Hello, " + name
	return greeting

}
```

```golang
// Unnecessary blank line at the end of an if block.
if user.IsAdmin() {
	grantAccess(user)

}
```

```golang
// Unnecessary blank line at the end of a for loop body.
for _, item := range items {
	process(item)

}
```

```golang
// Unnecessary blank line at the end of a switch block.
switch role {
case "admin":
	return adminDashboard()
default:
	return userDashboard()

}
```

### Valid

```golang
// No blank line at the end of the function body.
func greet(name string) string {
	greeting := "Hello, " + name
	return greeting
}
```

```golang
// No blank line at the end of the if block.
if user.IsAdmin() {
	grantAccess(user)
}
```

```golang
// No blank line at the end of the for loop body.
for _, item := range items {
	process(item)
}
```

```golang
// No blank line at the end of the switch block.
switch role {
case "admin":
	return adminDashboard()
default:
	return userDashboard()
}
```
