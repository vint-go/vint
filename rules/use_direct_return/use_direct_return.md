---
title: useDirectReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useDirectReturn`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useDirectReturn:
    # rule options here
```

## Details

Checking if an error is nil to just after return the error or nil is redundant. Instead of writing a redundant `if err != nil { return err }` followed by `return nil`, you can simplify the code by directly returning the result of the function call.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func foo() error {
	if err := bar(); err != nil {
		return err
	}
	return nil
}
```

### Valid

```golang
func foo() error {
	return bar()
}
```
