# RWMutex.RLock Contention on Package Getters

## Problem

Sync blocking profile from `k8s_all_vintlint_trace.out` shows 136.66s of
`sync.(*RWMutex).RLock` contention during vintlint0 rule execution:

| Method | Blocked time | % of total |
|--------|-------------|------------|
| `Package.Files()` | 76.97s | 19.85% |
| `Package.IsMain()` | 51.16s | 13.19% |
| `Package.TypesInfo()` | 8.34s | 2.18% |

The callers are AST-only rules running concurrently in the AST pool:
- `no_unused_constant.Apply` → `Package.Files()` — 46.12s
- `use_exported_comment.Apply` → `File.IsImportable()` → `Package.IsMain()` — 32.03s
- `no_blank_import.Apply` → `Package.IsMain()` — 11.69s

## Root Cause

`RWMutex.RLock()` internally does `atomic.AddInt32(&rw.readerCount, 1)` — a
write to a shared cache line. With many goroutines on a 32-core machine all
calling RLock concurrently, the atomic increment causes severe cache-line
bouncing even though there are no writers. The lock is protecting data that
is effectively immutable during rule execution.

In vintlint0's `processPackage`:
1. `pkg.AddFile()` populates `files` — before pool dispatch
2. `pkg.ExportedScanSortable()` populates `sortable` — before pool dispatch
3. `goVersion` is set at construction — never changes
4. `main` is computed lazily by `IsMain()` but once cached never changes

After step 2, these fields are all immutable. The RWMutex is unnecessary.

## Fix

Add a `Freeze()` method to `Package` and an `atomic.Bool` field:

```go
type Package struct {
    // ... existing fields ...
    frozen atomic.Bool
}

func (p *Package) Freeze() {
    p.IsMain()          // pre-compute and cache
    p.frozen.Store(true)
}
```

Add a fast path to each hot getter — if frozen, return the field directly
without touching the RWMutex:

```go
func (p *Package) Files() map[string]*File {
    if p.frozen.Load() {
        return p.files
    }
    p.mu.RLock()
    defer p.mu.RUnlock()
    return p.files
}

func (p *Package) IsMain() bool {
    if p.frozen.Load() {
        return p.main == trueValue
    }
    // ... existing double-checked locking ...
}
```

`atomic.Bool.Load()` on a value that never reverts to false is a pure read —
the cache line stays in Shared state across all cores with zero contention.
This is fundamentally different from RLock's atomic increment.

Apply the same pattern to: `Sortable()`, `IsAtLeastGoVersion()`,
`GoVersionString()`.

In vintlint0's `processPackage`, call `pkg.Freeze()` after
`ExportedScanSortable()` but before dispatching to pools:

```go
pkg.ExportedScanSortable()
pkg.Freeze()  // ← all pre-computed fields are now lock-free
```

The old linter path (`lint.Linter`) never calls `Freeze()`, so its behavior
is unchanged.

## Expected Impact

- Eliminates ~128s of RLock contention (Files + IsMain)
- No rule code changes required
- Backward compatible (old linter path unaffected)

## Other Findings from the Same Trace

- `runtime.chanrecv2` — 210.78s: pool workers idle waiting for work.
  This is structural (AST pool finishes faster than TC pool) and not
  directly fixable without rebalancing work across pools.
- `runtime.selectgo` — 9.70s: select-based channel operations.
- `sync.(*WaitGroup).Wait` — 8.15s: inside `go/types.(*Config).Check`.
