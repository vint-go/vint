# gocritic: two rules have no vint equivalent

## Affected rule
`noDuplicateOption` (maps to gocritic check `dupOption`) and `noNolintWithoutExplanationGocritic` (maps to gocritic check `whyNoLint`)

## Behavior in golangci-lint
- `dupOption`: detects duplicate options passed to function calls (e.g. repeated functional options).
- `whyNoLint`: requires `//nolint` directives to have an explanation comment.

## Behavior in vint
These two rules do not exist in vint's rule registry and are not mapped by the migrator.

## Gap
Users relying on these two gocritic checks will lose that coverage after migration.

## Example
```go
// dupOption check would flag this in gocritic:
http.NewServeMux().Handle("/", handler, WithTimeout(5), WithTimeout(10))

// whyNoLint check would flag this in gocritic:
x := something() //nolint:gocritic  // no explanation provided
```

## Impact on migration
Low impact. `dupOption` is rarely triggered. `whyNoLint` is a meta-linting check that becomes less relevant when migrating away from golangci-lint's nolint directive syntax.
