---
title: useSlicesSort
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useSlicesSort`
- This rule is not enabled by default, it must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useSlicesSort:
    # no configuration options
```

## Details

Since Go 1.21 the `slices` package proposes methods that are faster and easier to use
than their equivalents in the `sort` package.
This rule proposes to replace these legacy idioms with calls to the new methods.

The following replacements are suggested:

| sort method             | slices replacement        |
|-------------------------|---------------------------|
| `sort.Float64s`         | `slices.Sort`             |
| `sort.Ints`             | `slices.Sort`             |
| `sort.Strings`          | `slices.Sort`             |
| `sort.Slice`            | `slices.SortFunc`         |
| `sort.Sort`             | `slices.SortFunc`         |
| `sort.SliceStable`      | `slices.SortStableFunc`   |
| `sort.Stable`           | `slices.SortStableFunc`   |
| `sort.Float64sAreSorted`| `slices.IsSorted`         |
| `sort.IntsAreSorted`    | `slices.IsSorted`         |
| `sort.StringsAreSorted` | `slices.IsSorted`         |
| `sort.IsSorted`         | `slices.IsSortedFunc`     |
| `sort.SliceIsSorted`    | `slices.IsSortedFunc`     |

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
sort.Float64s(temperatures)
sort.Ints(years)
sort.Strings(names)
sort.Slice(cfg.Dependencies, func(i, j int) bool {
	return cfg.Dependencies[i].Name < cfg.Dependencies[j].Name
})
sort.Sort(sortable)
sort.Stable(sortable)
sort.SliceStable(years, func(i, j int) bool { return false })
sort.IsSorted(sortable)
sort.SliceIsSorted(years, func(i, j int) bool { return false })
sort.IntsAreSorted(years)
sort.StringsAreSorted(names)
sort.Float64sAreSorted(temperatures)
```

### Valid

```golang
slices.Sort(temperatures)
slices.Sort(years)
slices.Sort(names)
slices.SortFunc(cfg.Dependencies, func(a, b config.Dependency) int {
	return cmp.Compare(a.Name, b.Name)
})
slices.SortFunc(sortable, sortFunc)
slices.SortStableFunc(sortable, sortFunc)
slices.SortStableFunc(years, func(a, b int) int { return 0 })
slices.IsSortedFunc(sortable, sortFunc)
slices.IsSortedFunc(years, func(a, b int) int { return 0 })
slices.IsSorted(years)
slices.IsSorted(names)
slices.IsSorted(temperatures)
```
