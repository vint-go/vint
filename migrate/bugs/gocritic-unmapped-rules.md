# gocritic: one rule has no vint equivalent

## Affected rule
`noNolintWithoutExplanationGocritic` (maps to gocritic check `whyNoLint`)

**Resolved:** `noDuplicateOption` (maps to gocritic check `dupOption`) is now implemented as `lint/suspicious/noDuplicateOption`.

## Behavior in golangci-lint
- `whyNoLint`: requires `//nolint` directives to have an explanation comment.

## Behavior in vint
This rule does not exist in vint's rule registry and is not mapped by the migrator.

## Gap
Users relying on the `whyNoLint` gocritic check will lose that coverage after migration.

## Example
```go
// whyNoLint check would flag this in gocritic:
x := something() //nolint:gocritic  // no explanation provided
```

## Impact on migration
Low impact. `whyNoLint` is a meta-linting check that becomes less relevant when migrating away from golangci-lint's nolint directive syntax.
