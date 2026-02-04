# SQLite Database Creation Fix Plan

## Problem

Fresh Docker install crashes. Copying dev `log.db` fixes it.

## Observed Facts

1. Fresh Docker install fails with boot loop
2. Copying dev `log.db` to cache directory fixes it
3. tar.gz caching works in same directory with 777 permissions
4. `OpenTocDB()` has `CREATE TABLE IF NOT EXISTS` but database file is never created

## Root Cause

The SQLite connection string lacks the `mode=rwc` parameter. By default, some SQLite drivers only open existing files. The `rwc` mode explicitly tells SQLite to **Read, Write, and Create** the database file if it does not exist.

## Solution

Add `?mode=rwc` to the SQLite connection string. This is the proper SQLite way to enable auto-creation - no manual file creation hacks needed.

## Implementation

**File:** `internal/explore/sqlitedb.go`

### Single Change: Modify the sql.Open() call

**Current:**
```go
db, err := sql.Open("sqlite", path)
```

**Change to:**
```go
db, err := sql.Open("sqlite", path+"?mode=rwc")
```

The `mode=rwc` parameter tells SQLite to:
- **r** = read
- **w** = write
- **c** = create if not exists

This is the standard SQLite URI parameter for file creation. No additional imports or file manipulation needed.
