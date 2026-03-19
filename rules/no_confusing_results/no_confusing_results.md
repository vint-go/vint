---
title: noConfusingResults
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noConfusingResults`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noConfusingResults:
    # no configuration options
```

## Details

Functions or methods that return multiple, unnamed values of the same type could induce errors. When consecutive return values share the same type and are not named, it becomes easy to accidentally swap them or misinterpret which value is which at the call site. Using named return values makes the intent clearer and the code more maintainable.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// getPos yields the geographical position of this tag.
func (t tag) getPos() (float32, float32) {
	// Two unnamed float32 returns - which is longitude and which is latitude?
}
```

```golang
func getfoo() (int, int, error) {
	// Two consecutive unnamed int returns may be confusing
}
```

### Valid

```golang
// getPos yields the geographical position of this tag.
func (t tag) getPos() (longitude float32, latitude float32) {
	// Named return values make the intent clear
}
```

```golang
func getBar(a, b int) (int, error, int) {
	// No consecutive same-type unnamed returns
}
```

```golang
func namedResults() (a string, b string) {
	// Named results are fine even with same types
	return "nil", "nil"
}
```
