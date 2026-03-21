# mnd: setting ignored-numbers/ignored-functions/ignored-files replaces defaults instead of extending them

## Affected rules
- `lint/style/noMagicNumberInArgument` (maps to golangci-lint `mnd`)
- `lint/style/noMagicNumberInAssignment`
- `lint/style/noMagicNumberInCase`
- `lint/style/noMagicNumberInCondition`
- `lint/style/noMagicNumberInOperation`
- `lint/style/noMagicNumberInReturn`

## Behavior in golangci-lint
In golangci-lint's mnd linter, the `ignored-numbers` setting is additive. The numbers `0`, `0.0`, `1`, and `1.0` are **always** ignored regardless of the user's configuration. Similarly, `ignored-files` always includes `_test.go`, and `ignored-functions` always includes built-in exclusions like `time.Date`, `strconv.ParseInt`, etc.

For example, with config `ignored-numbers: ["2"]`, the effective ignored set is `{0, 0.0, 1, 1.0, 2}`.

## Behavior in vint
In vint's magic number rules, when `ignored-numbers` is provided via `Configure()`, it completely replaces the default set. The same applies to `ignored-functions` and `ignored-files`.

For example, with config `ignored-numbers: "2"`, the effective ignored set is `{2}` only -- `0`, `0.0`, `1`, and `1.0` would be flagged as magic numbers.

## Gap
When a user migrates from golangci-lint with custom `ignored-numbers`, `ignored-functions`, or `ignored-files`, the vint rules will be stricter than expected because the built-in defaults are lost.

## Example
```yaml
# golangci-lint config
linters:
  settings:
    mnd:
      ignored-numbers:
        - "2"
```

```go
// In golangci-lint: no issue (0 is always ignored)
// In vint with ignored-numbers="2": flagged as magic number
x := 0
```

## Impact on migration
Users who set any of `ignored-numbers`, `ignored-functions`, or `ignored-files` in their golangci-lint config will see more false positives after migration because the default exclusions are lost. The migrator could work around this by prepending the defaults to user-provided values, but this would couple the migrator to the vint rule implementation details. The cleaner fix would be for the vint rules to treat the `ignored-*` options as additive to defaults, matching golangci-lint's behavior.
