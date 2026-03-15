---
title: noMultiLineIfBreak
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMultiLineIfBreak`
- This rule is not recommended, meaning it is **not** enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMultiLineIfBreak:
    multi-if: true # Enable the multi-line if statement check. Default: false.
```

## Details

Detects multi-line if conditions where the condition does not start on the same line as the `if` keyword. When a complex condition spans multiple lines, the condition should begin immediately after the `if` keyword on the same line, rather than starting on a new line after the `if`.

This ensures consistent formatting of multi-line if statements and improves readability by making it immediately clear what the condition is about. Without this check, a developer might place the `if` keyword on its own line with the condition starting below, which can make the code harder to scan and visually similar to a block statement.

This check is **not enabled by default**. Set the `multi-if` option to `true` to activate it.

Source: https://github.com/ultraware/whitespace

## Examples

### Invalid

```golang
// The condition starts on a new line after the if keyword.
if
	a == 1 &&
	b == 2 &&
	c == 3 {
	doSomething()
}
```

```golang
// The condition starts on a new line after the if keyword.
if
	user.IsActive() &&
	user.HasPermission("admin") {
	grantAccess(user)
}
```

### Valid

```golang
// The condition starts on the same line as the if keyword.
if a == 1 &&
	b == 2 &&
	c == 3 {
	doSomething()
}
```

```golang
// The condition starts on the same line as the if keyword.
if user.IsActive() &&
	user.HasPermission("admin") {
	grantAccess(user)
}
```

```golang
// Single-line if statement; no issue.
if a == 1 {
	doSomething()
}
```
