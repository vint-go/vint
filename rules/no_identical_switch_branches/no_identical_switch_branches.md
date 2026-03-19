---
title: noIdenticalSwitchBranches
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIdenticalSwitchBranches`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIdenticalSwitchBranches:
    # no options
```

## Details

A `switch` with identical branches makes maintenance harder and might be a source of bugs. Duplicated branches should be consolidated in one case clause. Untagged switches are skipped because the order of case evaluation might be important.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
switch a {
case 1:
    foo()
case 2:
    bar()
case 3:
    foo() // identical to case 1
default:
    return newError("blah")
}
```

```golang
switch a {
case 1:
    foo()
case 3:
    foo() // identical to case 1
default:
}
```

### Valid

```golang
switch a {
case 1:
    foo()
case 2:
    bar()
case 3:
    baz()
default:
    return newError("blah")
}
```

```golang
// Fallthrough branches are not flagged
switch a {
case 1:
    foo()
    fallthrough
case 2:
    fallthrough
case 3:
    foo()
}
```

```golang
// Untagged switches are not checked
switch {
case a > b:
    foo()
default:
    foo()
}
```
