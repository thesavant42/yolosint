# SQLite Fix

## Working Code

`dbexperiment.go` line 14:
```go
db, err := sql.Open("sqlite", "file:experiment.db")
```

## Missing Function

`sqlitedb.go` is missing `OpenTocDB`. Add after line 18:

```go
func OpenTocDB(path string) (*TocDB, error) {
	db, err := sql.Open("sqlite", "file:"+path)
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
```