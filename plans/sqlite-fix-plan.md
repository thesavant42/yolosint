# Fix: Hardcode Cache and SQLite Paths

## The Problem

Container fails because SQLite cannot create `/cache/log.db` when the path is derived from environment variables incorrectly.

## The Fix

### 1. explore.go - SQLite Initialization

**Current broken code at line 111:**
```go
if cd := os.Getenv("CACHE_DIR"); cd != "" {
    db, err := OpenTocDB(filepath.Join(cd, "log.db"))
```

**Replace with:**
```go
if err := os.MkdirAll("/cache", 0755); err != nil {
    log.Printf("failed to create /cache: %v", err)
}
db, err := OpenTocDB("/cache/log.db")
if err != nil {
    log.Printf("failed to open /cache/log.db: %v", err)
} else {
    h.tocDB = db
}
```

### 2. cache.go - buildTocCache()

**Current broken code at line 638:**
```go
if cd := os.Getenv("CACHE_DIR"); cd != "" {
    caches = append(caches, &dirCache{dir: cd})
}
```

**Replace with:**
```go
caches = append(caches, &dirCache{dir: "/cache"})
```

### 3. cache.go - buildIndexCache()

**Current broken code at line 648:**
```go
if cd := os.Getenv("CACHE_DIR"); cd != "" {
```

**Replace with:**
```go
caches = append(caches, &dirCache{dir: "/cache"})
```

## Summary

| Location | Remove | Replace With |
|----------|--------|--------------|
| explore.go:111 | `os.Getenv("CACHE_DIR")` conditional | `os.MkdirAll("/cache", 0755)` + hardcoded `/cache/log.db` |
| cache.go:638 | `os.Getenv("CACHE_DIR")` conditional | Hardcoded `/cache` |
| cache.go:648 | `os.Getenv("CACHE_DIR")` conditional | Hardcoded `/cache` |

The SQLite driver with `mode=rwc` creates the file. We just need `os.MkdirAll` to ensure `/cache` directory exists first.
