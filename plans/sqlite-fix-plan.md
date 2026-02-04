# SQLite Fix Plan

## File: internal/explore/sqlitedb.go

### Current Code (Lines 20-24)

```go
func OpenTocDB(path string) (*TocDB, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?mode=rwc")
	if err != nil {
		return nil, err
	}
```

### Replacement Code (Lines 20-24)

```go
func OpenTocDB() (*TocDB, error) {
	db, err := sql.Open("sqlite", "/cache/log.db")
	if err != nil {
		return nil, err
	}
```

### Changes

1. Line 20: Remove `path string` parameter from function signature
2. Line 21: Replace `"file:"+path+"?mode=rwc"` with `"/cache/log.db"`

## File: internal/explore/explore.go

### Current Code (Around Line 113)

```go
db, err := OpenTocDB("/cache/log.db")
```

### Replacement Code

```go
db, err := OpenTocDB()
```

### Changes

1. Remove the argument from the OpenTocDB call
