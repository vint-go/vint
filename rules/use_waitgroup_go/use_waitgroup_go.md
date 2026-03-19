---
title: useWaitgroupGo
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useWaitgroupGo`
- This rule is not enabled by default, it must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useWaitgroupGo:
    # no configuration options
```

## Details

Since Go 1.25 the `sync` package proposes the [`WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go) method.
This method is a shorter and safer replacement for the idiom `wg.Add ... go { ... wg.Done ... }`.
The rule proposes to replace these legacy idioms with calls to the new method.

Limitations: The rule does not rely on type information but on variable names to identify waitgroups.
This means the rule searches for `wg` (the de facto standard name for wait groups);
if the waitgroup variable is named differently than `wg` the rule will skip it.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
wg := sync.WaitGroup{}

wg.Add(1)
go func() {
    defer wg.Done()
    doSomething()
}()
```

```golang
wg := sync.WaitGroup{}

for i, file := range filenames {
    wg.Add(1)
    go func(i int, filename string) {
        parsed[i], errors[i] = parseFile(filename)
        wg.Done()
    }(i, file)
}
wg.Wait()
```

### Valid

```golang
wg := sync.WaitGroup{}

wg.Go(func() {
    doSomething()
})
```

```golang
wg := sync.WaitGroup{}

for i, file := range filenames {
    wg.Go(func() {
        parsed[i], errors[i] = parseFile(filename)
    })
}
wg.Wait()
```
