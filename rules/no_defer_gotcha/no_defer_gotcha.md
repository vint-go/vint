---
title: noDeferGotcha
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeferGotcha`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeferGotcha:
    # By default all checks are enabled. To selectively enable only certain checks,
    # provide a list of check names as the first argument.
    # Available checks: loop, callChain, methodCall, return, recover, immediateRecover
    arguments:
      - ["loop", "callChain", "methodCall", "return", "recover", "immediateRecover"]
```

## Details

This rule warns on common mistakes when using `defer` statements. It currently alerts on the following situations:

| Name              | Description                                                                                                                                                                                     |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| callChain         | Even if deferring call-chains of the form `foo()()` is valid, it does not help code understanding (only the last call is deferred)                                                             |
| loop              | Deferring inside loops can be misleading (deferred functions are not executed at the end of the loop iteration but of the current function) and it could lead to exhausting the execution stack |
| methodCall        | Deferring a call to a method can lead to subtle bugs if the method does not have a pointer receiver                                                                                             |
| recover           | Calling `recover` outside a deferred function has no effect                                                                                                                                     |
| immediateRecover  | Calling `recover` at the time a defer is registered, rather than as part of the deferred callback (e.g. `defer recover()` or equivalent)                                                       |
| return            | Returning values from a deferred function has no effect                                                                                                                                         |

These gotchas are [described here](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01).

By default, all checks are enabled but it is possible to selectively enable them through configuration. Option names are case-insensitive and hyphens are ignored, so `callChain` and `call-chain` are equivalent.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func example() {
    for {
        defer func() {}() // prefer not to defer inside loops
    }

    defer tt.m() // be careful when deferring calls to methods without pointer receiver

    defer func() error {
        return errors.New("error") // return in a defer function has no effect
    }()

    defer recover() // recover must be called inside a deferred function

    recover() // recover must be called inside a deferred function

    defer helper(recover()) // recover is executed immediately, not inside the deferred function

    defer verify2(func() error {
        return nil
    })() // prefer not to defer chains of function calls
}
```

### Valid

```golang
func example() {
    defer cleanup()

    defer func() {
        if r := recover(); r != nil {
            log.Println("recovered:", r)
        }
    }()

    defer file.Close()
}
```
