# sloglint: vint rules not yet implemented

## Affected rules
- `lint/sloglint/no-mixed-args` (maps to golangci-lint `sloglint` / `no-mixed-args`)
- `lint/sloglint/kv-only` (maps to golangci-lint `sloglint` / `kv-only`)
- `lint/sloglint/attr-only` (maps to golangci-lint `sloglint` / `attr-only`)
- `lint/sloglint/no-global` (maps to golangci-lint `sloglint` / `no-global`)
- `lint/sloglint/context-only` (maps to golangci-lint `sloglint` / `context`)
- `lint/sloglint/static-msg` (maps to golangci-lint `sloglint` / `static-msg`)
- `lint/sloglint/msg-style` (maps to golangci-lint `sloglint` / `msg-style`)
- `lint/sloglint/no-raw-keys` (maps to golangci-lint `sloglint` / `no-raw-keys`)
- `lint/sloglint/key-naming-case` (maps to golangci-lint `sloglint` / `key-naming-case`)
- `lint/sloglint/forbidden-keys` (maps to golangci-lint `sloglint` / `forbidden-keys`)
- `lint/sloglint/args-on-sep-lines` (maps to golangci-lint `sloglint` / `args-on-sep-lines`)
- `lint/sloglint/allowed-keys` (no golangci-lint equivalent)

## Behavior in golangci-lint
sloglint ensures consistent code style when using log/slog. It provides 11 configurable checks covering:
- Mixed argument style detection (key-value vs attributes)
- Enforcing kv-only or attr-only style
- Disallowing global loggers
- Requiring context-aware slog methods
- Enforcing static log messages
- Message style (lowercased/capitalized)
- Disallowing raw string keys
- Key naming convention enforcement (snake/kebab/camel/pascal)
- Forbidden key lists
- Requiring args on separate lines

## Behavior in vint
All 12 sloglint rules are registered in the rules registry but do not have Go rule implementations in the `rules/` directory. The migrator correctly maps all golangci-lint configuration settings to vint rule paths, but the rules cannot fire during nolint conversion.

## Gap
- Config migration works correctly: enabling `sloglint` in golangci-lint with various settings produces the correct `vint.yaml` with properly mapped rules and options.
- Nolint conversion is degraded: `//nolint:sloglint` directives are replaced with a "no vint rules fired" note instead of proper `// vint-ignore` directives, because the rules have no implementation to detect violations.
- The `allowed-keys` setting exists in sloglint's standalone CLI but is not exposed in golangci-lint's `SlogLintSettings` struct, so it cannot be migrated from golangci-lint config.

## Example
Input:
```go
slog.Info("a user has logged in", "user_id", 42, slog.String("ip_address", "192.0.2.0")) //nolint:sloglint
```

Expected (once rules are implemented):
```go
// vint-ignore lint/sloglint/no-mixed-args: migrated from nolint:sloglint
slog.Info("a user has logged in", "user_id", 42, slog.String("ip_address", "192.0.2.0"))
```

Actual:
```go
// NOTE: nolint removed — no vint rules fired (was: //nolint:sloglint)
slog.Info("a user has logged in", "user_id", 42, slog.String("ip_address", "192.0.2.0"))
```

## Impact on migration
Users migrating from golangci-lint will have their `//nolint:sloglint` directives removed with a note rather than converted to `// vint-ignore` directives. The config migration itself is fully functional. Once the vint rules are implemented, the nolint test should be updated to verify proper conversion.
