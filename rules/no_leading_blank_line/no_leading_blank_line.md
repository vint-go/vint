---
title: noLeadingBlankLine
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noLeadingBlankLine`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noLeadingBlankLine:
    # No additional options. This check is enabled by default.
```

## Details

Detects unnecessary blank lines at the beginning of a function body, if/else block, for loop, switch statement, or case clause. Leading blank lines inside blocks add visual noise without improving readability and create inconsistent code style across a codebase.

The whitespace linter automatically removes these blank lines when a fix is applied, bringing the first statement in the block directly after the opening brace.

This check runs by default and does not need to be explicitly enabled.

Source: https://github.com/ultraware/whitespace

## Examples

### Invalid

```golang
// Unnecessary blank line at the start of a function body.
func processOrder(id int) error {

	order, err := fetchOrder(id)
	if err != nil {
		return err
	}
	return order.Process()
}
```

```golang
// Unnecessary blank line at the start of an if block.
if err != nil {

	log.Println(err)
	return err
}
```

```golang
// Unnecessary blank line at the start of a for loop.
for i := 0; i < 10; i++ {

	fmt.Println(i)
}
```

```golang
// Unnecessary blank line at the start of a switch case clause.
switch status {
case "active":

	activate(user)
case "inactive":
	deactivate(user)
}
```

### Valid

```golang
// No blank line at the start of the function body.
func processOrder(id int) error {
	order, err := fetchOrder(id)
	if err != nil {
		return err
	}
	return order.Process()
}
```

```golang
// No blank line at the start of the if block.
if err != nil {
	log.Println(err)
	return err
}
```

```golang
// No blank line at the start of the for loop.
for i := 0; i < 10; i++ {
	fmt.Println(i)
}
```

```golang
// No blank line at the start of a switch case clause.
switch status {
case "active":
	activate(user)
case "inactive":
	deactivate(user)
}
```
