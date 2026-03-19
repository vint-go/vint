---
title: noForbiddenCallInWgGo
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noForbiddenCallInWgGo`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noForbiddenCallInWgGo:
    # no configuration options
```

## Details

Since Go 1.25, it is possible to create goroutines with the method `sync.WaitGroup.Go`.
The `Go` method calls a function in a new goroutine and adds (`Add`) that task to the WaitGroup.
When the function returns, the task is removed (`Done`) from the WaitGroup.

This rule ensures that functions passed to `wg.Go` don't call `panic`, `log.Panic`, `log.Panicf`,
`log.Panicln`, or `wg.Done`, as specified in the
[documentation of `WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go).

The rule also warns against a common mistake when refactoring legacy code:
accidentally leaving behind a call to `WaitGroup.Done`, which can cause subtle bugs or panics.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
  wg.Done() // wg.Done is called automatically by wg.Go
})

wg.Wait()
```

```golang
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
  panic("something went wrong") // panic is forbidden inside wg.Go
})

wg.Wait()
```

### Valid

```golang
wg := sync.WaitGroup{}

wg.Go(func() {
  doSomething()
})

wg.Wait()
```
