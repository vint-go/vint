# gocheckcompilerdirectives: unrecognized directive detection is narrower in vint

## Affected rule
`lint/correctness/noMalformedDirective` (maps to golangci-lint `gocheckcompilerdirectives`)

## Behavior in golangci-lint
The original `gocheckcompilerdirectives` linter flags ALL unrecognized `//go:` directives, regardless of how different they are from known directives. For example, `//go:fakething`, `//go:custom`, or `//go:myapp` would all be flagged.

## Behavior in vint
The `noMalformedDirective` rule only flags unrecognized `//go:` directives if they are within Levenshtein distance 2 of a known directive (i.e., likely misspellings). Completely unrelated directive names are silently ignored.

## Gap
Directives with names that are not close to any known directive (edit distance > 2) are not detected by vint. For example:
- `//go:fakething` -- not flagged (too far from any known directive)
- `//go:custom` -- not flagged
- `//go:myapp` -- not flagged

Meanwhile, misspellings within distance 2 are correctly detected:
- `//go:embod` -- flagged (close to `embed`)
- `//go:genrate` -- flagged (close to `generate`)
- `//go:linknme` -- flagged (close to `linkname`)

## Example
```go
package example

//go:fakething
func Foo() {}

//go:custom
func Bar() {}
```

The original `gocheckcompilerdirectives` would flag both `//go:fakething` and `//go:custom`. The vint `noMalformedDirective` rule flags neither, since both are more than 2 edits away from any known directive.

## Impact on migration
Users who relied on `gocheckcompilerdirectives` to catch completely unrecognized or custom `//go:` directives will lose coverage for those cases after migration. Only misspellings of known directives will continue to be caught. This is a moderate gap -- misspellings are the more common real-world error, but fully unknown directives could also indicate mistakes.
