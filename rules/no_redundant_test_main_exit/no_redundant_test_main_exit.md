---
title: noRedundantTestMainExit
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantTestMainExit`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantTestMainExit:
    # no configuration options
```

## Details

This rule warns about redundant `os.Exit` or `syscall.Exit` calls in `TestMain` functions.
Starting from Go 1.15, the test runner automatically handles program termination, making
explicit exit calls in `TestMain` unnecessary.

The rule only applies to test files (`_test.go`) and only when the module's Go version is 1.15 or higher.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func TestMain(m *testing.M) {
	setup()
	i := m.Run()
	teardown()
	os.Exit(i) // redundant: the test runner handles this automatically as of Go 1.15
}
```

```golang
func TestMain(m *testing.M) {
	i := m.Run()
	syscall.Exit(i) // redundant: the test runner handles this automatically as of Go 1.15
}
```

### Valid

```golang
func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}
```

```golang
func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
}
```
