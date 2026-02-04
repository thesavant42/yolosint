//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:experiment.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		log.Fatal(err)
	}

	fmt.Println("Created experiment.db with layers and files tables")
}
