package explore

import (
	"database/sql"
	"sync"
	"time"

	_ "modernc.org/sqlite" // Register SQLite driver with database/sql

	"github.com/thesavant42/yolosint/internal/soci"
)

// TocDB wraps SQLite connection for TOC metadata logging.
// Why struct with mutex: Multiple goroutines call Put() concurrently; SQLite needs serialized writes.
type TocDB struct {
	db *sql.DB
	mu sync.Mutex
}

// db, err := sql.Open("sqlite", "file:experiment.db") works
func OpenTocDB() (*TocDB, error) {
	db, err := sql.Open("sqlite", "file:/cache/log.db")
	if err != nil {
		return nil, err
	}

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
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &TocDB{db: db}, nil
}

func (t *TocDB) Insert(digest string, toc *soci.TOC) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	// Why defer Rollback: Safe cleanup if any insert fails mid-transaction. No-op if Commit succeeds.
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

	// Why Prepare: Reuse statement for batch insert efficiency.
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

	return tx.Commit()
}

func (t *TocDB) Close() error {
	return t.db.Close()
}
