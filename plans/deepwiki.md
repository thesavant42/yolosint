# DeepWiki Q&A with Code Context for Repository: modernc-org/sqlite
## Q1
how do I create a database if it doesnt exist?
### Answer
In `modernc.org/sqlite`, simply opening a connection creates the database file if it does not exist. To safely create schema, use `CREATE TABLE ... IF NOT EXISTS` statements.

**Quick example**
```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file:test.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create a table only if it doesn’t exist
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
    if err != nil {
        panic(err)
    }
}
```

**Why this works**
- `sql.Open("sqlite", dsn)` creates the file if missing; see the basic open pattern in `doc.go` sqlite:59-72 .
- Tests use `t.TempDir()` to open a new database file and then `db.Exec("create table ...")` to create schema sqlite:774-793 .
- SQLite’s `IF NOT EXISTS` clause ensures the command is a no-op when the table already exists sqlite:510-529 .

**Notes**
- The DSN can include pragmas, e.g., `?_pragma=journal_mode(WAL)` sqlite:59-72 .
- `CREATE TABLE IF NOT EXISTS` does not suppress errors if an index with the same name exists sqlite:531-540 .

Wiki pages you might want to explore:
- [Overview (modernc-org/sqlite)](/wiki/modernc-org/sqlite#1)
### Citations
**File:** doc.go (L59-72)
```go
// To access a Sqlite database do something like
//
//	import (
//		"database/sql"
//
//		_ "modernc.org/sqlite"
//	)
//
//	...
//
//
//	db, err := sql.Open("sqlite", dsnURI)
//
//	...
```
**File:** all_test.go (L774-793)
```go
	path := filepath.Join(t.TempDir(), "sqlite-test-")
	db, err := sql.Open(driverName, path)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	if strings.HasPrefix(conns, "single") {
		db.SetMaxOpenConns(1)
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.Exec("create table t(i)"); err != nil {
		t.Fatal(err)
	}
```
**File:** testdata/tcl/e_createtable.test (L510-529)
```text
# EVIDENCE-OF: R-33917-24086 However, if the "IF NOT EXISTS" clause is
# specified as part of the CREATE TABLE statement and a table or view of
# the same name already exists, the CREATE TABLE command simply has no
# effect (and no error message is returned).
#
drop_all_tables
do_execsql_test e_createtable-1.8.0 {
  CREATE TABLE t1(x, y);
  CREATE INDEX i1 ON t1(x);
  CREATE VIEW  v1 AS SELECT * FROM t1;
  CREATE TABLE auxa.tbl1(x, y);
  CREATE INDEX auxa.idx1 ON tbl1(x);
  CREATE VIEW auxa.view1 AS SELECT * FROM tbl1;
} {}
do_createtable_tests 1.8 {
  1    "CREATE TABLE IF NOT EXISTS t1(a, b)"          {}
  2    "CREATE TABLE IF NOT EXISTS auxa.tbl1(a, b)"   {}
  3    "CREATE TABLE IF NOT EXISTS v1(a, b)"          {}
  4    "CREATE TABLE IF NOT EXISTS auxa.view1(a, b)"  {}
}
```
**File:** testdata/tcl/e_createtable.test (L531-540)
```text
# EVIDENCE-OF: R-16465-40078 An error is still returned if the table
# cannot be created because of an existing index, even if the "IF NOT
# EXISTS" clause is specified.
#
do_createtable_tests 1.9 -error { %s } {
  1    "CREATE TABLE IF NOT EXISTS i1(a, b)"   
       {{there is already an index named i1}}
  2    "CREATE TABLE IF NOT EXISTS auxa.idx1(a, b)"   
       {{there is already an index named idx1}}
}
```
## Q2
what if the path is in the /cache/ directory? How do I create `/cache/log.db`
### Answer
To create `/cache/log.db` with `modernc.org/sqlite`, pass the absolute path in the DSN; the driver will create the file if it doesn’t exist. Ensure the process has write permission to `/cache/`.

**Example**
```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file:/cache/log.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY, message TEXT)")
    if err != nil {
        panic(err)
    }
}
```

**Why this works**
- `newConn` calls `openV2` with `SQLITE_OPEN_CREATE`, which creates the file if missing sqlite:36-63 .
- The DSN `file:/cache/log.db` is an absolute path; SQLite resolves it directly sqlite:57-63 .
- Tests use `filepath.Join(tempDir, "test.db")` to open databases in arbitrary directories, confirming the driver respects the provided path sqlite:2273-2285 .

**Notes**
- If `/cache/` doesn’t exist or the process lacks permissions, `openV2` will fail with “unable to open database file” sqlite:96-101 .
- You can add pragmas to the DSN, e.g., `file:/cache/log.db?_pragma=journal_mode(WAL)` sqlite:59-72 .

Wiki pages you might want to explore:
- [Core Driver (modernc-org/sqlite)](/wiki/modernc-org/sqlite#2)
### Citations
**File:** conn.go (L36-63)
```go
func newConn(dsn string) (*conn, error) {
	var query, vfsName string

	// Parse the query parameters from the dsn and them from the dsn if not prefixed by file:
	// https://github.com/mattn/go-sqlite3/blob/3392062c729d77820afc1f5cae3427f0de39e954/sqlite3.go#L1046
	// https://github.com/mattn/go-sqlite3/blob/3392062c729d77820afc1f5cae3427f0de39e954/sqlite3.go#L1383
	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		query = dsn[pos+1:]
		var err error
		vfsName, err = getVFSName(query)
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(dsn, "file:") {
			dsn = dsn[:pos]
		}
	}

	c := &conn{tls: libc.NewTLS()}
	db, err := c.openV2(
		dsn,
		vfsName,
		sqlite3.SQLITE_OPEN_READWRITE|sqlite3.SQLITE_OPEN_CREATE|
			sqlite3.SQLITE_OPEN_FULLMUTEX|
			sqlite3.SQLITE_OPEN_URI,
	)
```
**File:** all_test.go (L2273-2285)
```go
// https://gitlab.com/cznic/sqlite/-/issues/66
func TestIssue66(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		os.RemoveAll(tempDir)
	}()

	fn := filepath.Join(tempDir, "testissue66.db")
	db, err := sql.Open(driverName, fn)
```
**File:** testdata/tcl/oserror.test (L96-101)
```text
do_test 1.3.1 {
  set ::log [list]
  list [catch { sqlite3 dbh /x/y/z/test.db } msg] $msg
} {1 {unable to open database file}}

do_re_test 1.3.2 { lindex $::log 0 } {^os_unix.c:\d+: \(\d+\) open\(.*test.db\) - }
```