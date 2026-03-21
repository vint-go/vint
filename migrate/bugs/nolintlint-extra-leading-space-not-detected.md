# nolintlint: extra leading spaces in nolint directives not detected

## Affected rule
`lint/style/noNolintNonMachineReadable` (maps to golangci-lint `nolintlint`, extractor rule `extra-leading-space`)

## Behavior in golangci-lint
The original nolintlint linter has a separate check for extra leading spaces. It detects directives like `//  nolint:errcheck` or `//   nolint:gosec` (with more than one space between `//` and `nolint`) and flags them as having extra leading whitespace.

## Behavior in vint
The vint rule `noNolintNonMachineReadable` uses `strings.HasPrefix(text, "// nolint")` to detect non-machine-readable directives. This correctly catches `// nolint:errcheck` (single space) but misses `//  nolint:errcheck` (two spaces) or `//   nolint:errcheck` (three spaces), because the prefix check fails when the third character is a space instead of `n`.

## Gap
Directives with more than one space between `//` and `nolint` (e.g., `//  nolint:errcheck`) are not detected. The rules registry marks `noNolintExtraSpace` as `subsumed_by = "noNolintNonMachineReadable"`, indicating that `noNolintNonMachineReadable` should cover this case, but the current implementation does not.

## Example
```go
package example

//  nolint:errcheck
func foo() {
    _ = doSomething()
}
```

In golangci-lint with nolintlint enabled, this would be flagged as having extra leading space. In vint, `noNolintNonMachineReadable` does not fire on this code.

## Impact on migration
Users who had nolint directives with extra spaces (e.g., `//  nolint`) will not see warnings from vint after migration. This is a minor gap since such directives are uncommon, but it means vint's coverage is slightly less strict than the original nolintlint linter for this edge case.
