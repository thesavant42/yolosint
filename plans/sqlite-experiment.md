# SQLite Fix

## Experiment Result

`go run dbexperiment.go` succeeded locally. Created `experiment.db` with schema.

## Git History Analysis

Commit `d0b4ab9` ("works in local container") had:

```go
func OpenTocDB(path string) (*TocDB, error) {
    db, err := sql.Open("sqlite", path)
```

Called from explore.go:

```go
db, err := OpenTocDB(filepath.Join(cd, "log.db"))
```

**No `file:` prefix.** Raw path passed directly.

## Proposed Fix

Revert sqlitedb.go to accept path parameter:

```go
func OpenTocDB(path string) (*TocDB, error) {
    db, err := sql.Open("sqlite", path)
```

And restore the caller in explore.go to pass the path:

```go
db, err := OpenTocDB(filepath.Join(cd, "log.db"))
```

The `file:` prefix was added during debugging and broke it.
