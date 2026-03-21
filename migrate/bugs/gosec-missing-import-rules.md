# gosec: Missing import-ban and unchecked-error rules

## Affected rules
- `lint/security/noMd5Import` (maps to gosec `G501`)
- `lint/security/noDesImport` (maps to gosec `G502`)
- `lint/security/noRc4Import` (maps to gosec `G503`)
- `lint/security/noSha1Import` (maps to gosec `G505`)
- `lint/security/noMd4Import` (maps to gosec `G506`)
- `lint/security/noRipemd160Import` (maps to gosec `G507`)
- `lint/security/noUncheckedErrorGosec` (maps to gosec `G104`)

## Behavior in golangci-lint
These rules are enabled by default when gosec is active:
- G501-G507 flag imports of deprecated or insecure cryptographic packages (crypto/md5, crypto/des, crypto/rc4, crypto/sha1, golang.org/x/crypto/md4, golang.org/x/crypto/ripemd160, net/http/cgi).
- G104 flags unchecked errors (similar to errcheck but gosec-specific).

## Behavior in vint
These 7 rules are registered in the rules registry but have no vint rule implementations. The migrator silently omits them.

## Gap
Users who rely on these gosec rules will lose coverage after migration. The import-ban rules (G501-G507, except G504 which maps to noCgiImport) and G104 (unchecked errors) will not be enforced.

## Example
```go
import "crypto/md5" // G501 would flag this in gosec, but no vint rule fires
```

## Impact on migration
Low-to-medium impact. G104 (unchecked errors) is partially covered by errcheck's `noUncheckedError` rule if errcheck is also enabled. The import-ban rules for deprecated crypto packages are niche but relevant for security-sensitive projects. Users should review that these checks are covered by other means after migration.
