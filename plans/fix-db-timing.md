# Fix Database Initialization Timing

## Problem
Database initialization happens during `explore.New()` at app startup. The `/cache` volume mount is not ready yet. The app crashes immediately in a bootloop before any output.

## Solution
Move database initialization from startup to lazy initialization on first use, with print statements at each step.

## File: internal/explore/sqlitedb.go

Replace the entire file with:

```go
package explore

import (
    "database/sql"
    "log"
    "os"
    "path/filepath"
    "sync"
    "time"

    _ "modernc.org/sqlite"

    "github.com/thesavant42/yolosint/internal/soci"
)

type TocDB struct {
    db   *sql.DB
    mu   sync.Mutex
    path string
    once sync.Once
    err  error
}

func NewTocDB(path string) *TocDB {
    log.Printf("[DB] NewTocDB called with path=%s", path)
    return &TocDB{path: path}
}

func (t *TocDB) init() error {
    t.once.Do(func() {
        log.Printf("[DB] init: starting lazy initialization")
        log.Printf("[DB] init: path=%s", t.path)
        
        dir := filepath.Dir(t.path)
        log.Printf("[DB] init: checking directory %s", dir)
        
        if info, err := os.Stat(dir); err != nil {
            log.Printf("[DB] init: directory does not exist, err=%v", err)
        } else {
            log.Printf("[DB] init: directory exists, mode=%v", info.Mode())
        }
        
        log.Printf("[DB] init: calling MkdirAll on %s", dir)
        if err := os.MkdirAll(dir, 0755); err != nil {
            log.Printf("[DB] init: MkdirAll failed, err=%v", err)
            t.err = err
            return
        }
        log.Printf("[DB] init: MkdirAll succeeded")
        
        log.Printf("[DB] init: checking if database file exists")
        if info, err := os.Stat(t.path); err != nil {
            log.Printf("[DB] init: database file does not exist, err=%v", err)
        } else {
            log.Printf("[DB] init: database file exists, size=%d", info.Size())
        }
        
        log.Printf("[DB] init: calling sql.Open")
        db, err := sql.Open("sqlite", t.path)
        if err != nil {
            log.Printf("[DB] init: sql.Open failed, err=%v", err)
            t.err = err
            return
        }
        log.Printf("[DB] init: sql.Open succeeded")
        
        schema := `
        CREATE TABLE IF NOT EXISTS layers (
            id INTEGER PRIMARY KEY,
            digest TEXT NOT NULL,
            csize INTEGER,
            usize INTEGER,
            type TEXT,
            media_type TEXT,
            indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS files (
            id INTEGER PRIMARY KEY,
            layer_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            typeflag INTEGER,
            size INTEGER,
            mode INTEGER,
            mod DATETIME,
            offset INTEGER,
            linkname TEXT,
            FOREIGN KEY(layer_id) REFERENCES layers(id)
        );
        CREATE INDEX IF NOT EXISTS idx_files_name ON files(name);
        CREATE INDEX IF NOT EXISTS idx_layers_digest ON layers(digest);
        `
        
        log.Printf("[DB] init: executing schema")
        if _, err := db.Exec(schema); err != nil {
            log.Printf("[DB] init: schema exec failed, err=%v", err)
            db.Close()
            t.err = err
            return
        }
        log.Printf("[DB] init: schema exec succeeded")
        
        log.Printf("[DB] init: checking database file after schema")
        if info, err := os.Stat(t.path); err != nil {
            log.Printf("[DB] init: database file STILL does not exist after schema, err=%v", err)
        } else {
            log.Printf("[DB] init: database file exists after schema, size=%d", info.Size())
        }
        
        t.db = db
        log.Printf("[DB] init: initialization complete")
    })
    return t.err
}

func (t *TocDB) Insert(digest string, toc *soci.TOC) error {
    log.Printf("[DB] Insert: called for digest=%s", digest)
    if err := t.init(); err != nil {
        log.Printf("[DB] Insert: init failed, err=%v", err)
        return err
    }
    log.Printf("[DB] Insert: init succeeded, proceeding with insert")
    
    t.mu.Lock()
    defer t.mu.Unlock()

    tx, err := t.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    res, err := tx.Exec(
        `INSERT INTO layers (digest, csize, usize, type, media_type) VALUES (?, ?, ?, ?, ?)`,
        digest, toc.Csize, toc.Usize, toc.Type, toc.MediaType,
    )
    if err != nil {
        return err
    }

    layerID, err := res.LastInsertId()
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare(`INSERT INTO files (layer_id, name, typeflag, size, mode, mod, offset, linkname) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, f := range toc.Files {
        var modTime *time.Time
        if !f.ModTime.IsZero() {
            modTime = &f.ModTime
        }
        if _, err := stmt.Exec(layerID, f.Name, f.Typeflag, f.Size, f.Mode, modTime, f.Offset, f.Linkname); err != nil {
            return err
        }
    }

    log.Printf("[DB] Insert: committing transaction")
    return tx.Commit()
}

func (t *TocDB) Close() error {
    log.Printf("[DB] Close: called")
    if t.db != nil {
        return t.db.Close()
    }
    return nil
}
```

## File: internal/explore/explore.go

**Line 113-117 - Change from:**
```go
db, err := OpenTocDB()
if err != nil {
    log.Fatalf("failed to open log.db: %v", err)
}
h.tocDB = db
```

**Change to:**
```go
h.tocDB = NewTocDB("/cache/log.db")
```

No other changes needed.
