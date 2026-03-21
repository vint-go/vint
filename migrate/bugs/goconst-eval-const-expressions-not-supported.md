# goconst: eval-const-expressions setting not supported

## Affected rule
`lint/style/noRepeatedStrings` (maps to golangci-lint `goconst`)

## Behavior in golangci-lint
When goconst is configured with `eval-const-expressions: true`, it evaluates constant expressions like `Prefix + "suffix"` when analyzing string literals. This allows it to detect cases where a string is effectively repeated through constant concatenation.

Example flagged with `eval-const-expressions: true`:
```go
const Prefix = "api/"
var endpoint1 = Prefix + "users"
var endpoint2 = "api/users" // flagged: matches Prefix + "users"
```

## Behavior in vint
The vint `noRepeatedStrings` rule only examines raw string literals in the source code. It does not evaluate constant expressions or perform string concatenation analysis.

## Gap
The `eval-const-expressions` feature is not available in vint. Users who relied on this golangci-lint feature will lose coverage for detecting string duplication involving constant concatenation.

## Example
```go
package example

const Prefix = "api/"

func routes() {
    _ = Prefix + "users"
    _ = "api/users"   // golangci-lint flags this; vint does not
    _ = "api/users"   // golangci-lint flags this; vint does not
}
```

## Impact on migration
- The `eval-const-expressions` setting is silently dropped during migration since there is no equivalent vint option.
- This is a low-severity gap as `eval-const-expressions` is disabled by default in golangci-lint and few users enable it.
- Users who relied on this feature will lose coverage for constant expression evaluation.
