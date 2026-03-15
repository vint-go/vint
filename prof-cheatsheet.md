# pprof cheatsheet — command-line one-liners

## CPU

```bash
# Top functions by self CPU time
go tool pprof -top cpu.prof

# Top functions by cumulative time (self + callees)
go tool pprof -top -cum cpu.prof

# Filter to your module only
go tool pprof -top -cum -trim_path=github.com/strowk/vint cpu.prof

# Show more entries
go tool pprof -top -cum -nodecount=50 cpu.prof

# Flat text call tree
go tool pprof -text -cum cpu.prof

# Source-annotated view of a specific function
go tool pprof -list='file\.lint' cpu.prof
go tool pprof -list='actionID' cpu.prof
go tool pprof -list='TypeCheck' cpu.prof

# Show callers of a function (who calls it, how much)
go tool pprof -peek='ast\.Walk' cpu.prof

# Show call tree rooted at a function
go tool pprof -tree -focus='Package\.lint' cpu.prof

# Flame graph in browser
go tool pprof -http=:8080 cpu.prof
```

## Memory (in-use at profile time)

```bash
# Top by bytes still in use
go tool pprof -top -inuse_space mem.prof

# Top by objects still in use
go tool pprof -top -inuse_objects mem.prof

# Source-annotated allocations in a function
go tool pprof -list='lintPackage' -inuse_space mem.prof
```

## Memory (total allocated, including freed)

```bash
# Top by total bytes allocated
go tool pprof -top -alloc_space mem.prof

# Top by total objects allocated (GC pressure)
go tool pprof -top -alloc_objects mem.prof

# Source-annotated
go tool pprof -list='file\.lint' -alloc_space mem.prof

# Flame graph in browser
go tool pprof -http=:8081 mem.prof
```

## Drill down workflow

```bash
# 1. Find the hot spot
go tool pprof -top -cum -nodecount=10 cpu.prof

# 2. See who calls it
go tool pprof -peek='<hot function>' cpu.prof

# 3. See annotated source
go tool pprof -list='<hot function>' cpu.prof

# 4. See what it calls
go tool pprof -tree -focus='<hot function>' -nodecount=20 cpu.prof
```

## Comparing two profiles

```bash
# Diff CPU profiles (before vs after optimization)
go tool pprof -top -cum -diff_base=before.prof after.prof

# Diff in browser
go tool pprof -http=:8080 -diff_base=before.prof after.prof
```
