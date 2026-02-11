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

// ImageContext holds the parsed image reference metadata for searchable logging
type ImageContext struct {
	Registry   string // e.g., ghcr.io, docker.io
	Namespace  string // e.g., chainguard, library (first path segment)
	Repository string // e.g., nginx, go (remaining path after namespace)
	Tag        string // e.g., latest, v1.0.0 (may be empty for digest refs)
	ImageRef   string // full reference: ghcr.io/chainguard/nginx:latest
}

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
		          indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		          registry TEXT,
		          namespace TEXT,
		          repository TEXT,
		          tag TEXT,
		          image_ref TEXT
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
		      CREATE INDEX IF NOT EXISTS idx_layers_registry ON layers(registry);
		      CREATE INDEX IF NOT EXISTS idx_layers_namespace ON layers(namespace);
		      CREATE INDEX IF NOT EXISTS idx_layers_repository ON layers(repository);
		      CREATE INDEX IF NOT EXISTS idx_layers_image_ref ON layers(image_ref);
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

func (t *TocDB) Insert(digest string, toc *soci.TOC, imgCtx *ImageContext) error {
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

	// Handle nil ImageContext gracefully
	var registry, namespace, repository, tag, imageRef *string
	if imgCtx != nil {
		if imgCtx.Registry != "" {
			registry = &imgCtx.Registry
		}
		if imgCtx.Namespace != "" {
			namespace = &imgCtx.Namespace
		}
		if imgCtx.Repository != "" {
			repository = &imgCtx.Repository
		}
		if imgCtx.Tag != "" {
			tag = &imgCtx.Tag
		}
		if imgCtx.ImageRef != "" {
			imageRef = &imgCtx.ImageRef
		}
	}

	res, err := tx.Exec(
		`INSERT INTO layers (digest, csize, usize, type, media_type, registry, namespace, repository, tag, image_ref) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		digest, toc.Csize, toc.Usize, toc.Type, toc.MediaType, registry, namespace, repository, tag, imageRef,
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
