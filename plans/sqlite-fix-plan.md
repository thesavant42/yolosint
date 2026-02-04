# SQLite Fix Plan

## Root Cause

The modernc.org/sqlite driver creates database files automatically when given a simple path. The `file:` URI format with `?mode=rwc` is causing the failure.

Reference - modernc.org/sqlite usage:
```go
import (
    "database/sql"
    _ "modernc.org/sqlite" // Register CGO-free driver
)

func main() {
    // Opens (and creates if missing) the local file "data.db"
    db, _ := sql.Open("sqlite", "data.db")
    defer db.Close()

    // Initialize schema (only runs if the table doesn't exist)
    db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)`)
}
```

## Fix sqlitedb.go Line 21

Current broken code:
```go
db, err := sql.Open("sqlite", "file:"+path+"?mode=rwc")
```

Change to:
```go
db, err := sql.Open("sqlite", "/cache/log.db")
```

Remove the `path string` parameter from the function signature on line 20.

Remove `"os"` from imports if present.

## Fix explore.go Line 113

Change:
```go
db, err := OpenTocDB()
```

No change needed if already calling without argument.

## Fix explore.go Lines 114-118

Current broken code:
```go
if err != nil {
    log.Printf("failed to open /cache/log.db: %v", err)
} else {
    h.tocDB = db
}
```

Change to:
```go
if err != nil {
    log.Fatalf("failed to open /cache/log.db: %v", err)
}
h.tocDB = db
```

The app must crash if SQLite fails to initialize.
