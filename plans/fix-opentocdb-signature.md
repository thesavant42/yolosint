# Fix OpenTocDB Signature Mismatch

## Problem
Docker build fails with:
```
internal/explore/explore.go:113:13: not enough arguments in call to OpenTocDB
  have ()
  want (string)
```

## Root Cause
- `OpenTocDB(path string)` in `sqlitedb.go:21` expects a string parameter
- `OpenTocDB()` is called without arguments at `explore.go:113`
- The `path` parameter is **unused** - function hardcodes `"file:/cache/log.db"` regardless

## Fix
Remove the unused `path` parameter from the function signature:

**File:** `internal/explore/sqlitedb.go` line 21

Change:
```go
func OpenTocDB(path string) (*TocDB, error) {
```

To:
```go
func OpenTocDB() (*TocDB, error) {
```

## Tasks
- [x] Identify compilation error
- [ ] Remove unused `path` parameter from `OpenTocDB` in `sqlitedb.go:21`
- [ ] Verify Docker build succeeds
