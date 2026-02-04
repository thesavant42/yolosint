ing# Race condition in sqlite logic causes app failure

## Problem Statement

The application fails and the docker container boot loops if the sqlite table does not exist (*good*) BUT **it never creates the sqlite table.**

---

## Critical Insights 

- **CRITICAL INSIGHT**: The mount volume in docker compose works fine for the cached tar.gz files; this continues to work. The directory permissions are `777`.
- **CRITICSL INSIGHT**: when I copied over **my sqlite dev table it worked**.

### Summary of Key Insights
- On the dev machine where I had already created a sqlite table during testing, everything works.
- But on the system that was not dev, the **fresh install**, it fails to boot **because there's no sqlite database**. 
    - **Once I copied over my dev instance, the app worked. **

---

## Solution: create sqlite db on installation

- Add a shell sqlite database
    - can't be an empty file, needs to be an empty sqlite db with the schema from [log.db](/cache/log.db)
sqlite.go has a create db if not exists function, but is it ever used?

internal\explore\sqlitedb.go:26-53
```
// Why two tables: One layer has many files. FK enables "find all layers containing file X".
	// Why indexes: idx_files_name is the primary search use case; idx_layers_digest for layer lookups.
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
```
---

## Task:

1. Troubleshoot the root cause of the issue, do not speculate.
    - Cite your references and never guess. Thoerizing is fine if you address it as such preemptively and have a valid experiment for KPIs.
2. propose a solution;
3. I'd rather the code properly create the sqlite db the first time, rather than try to commit an empty db to git

