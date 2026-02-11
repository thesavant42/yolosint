
DBMS: modernc sqlite: https://gitlab.com/cznic/sqlite
db: cache/log.db
    - use MCP

- [/internal/explore/sqlitedb.go](/internal/explore/sqlitedb.go)
- [/internal/explore/soci.go](/internal/explore/soci.go)
- [/internal/explore/explore.go](/internal/explore/explore.go)
- [/internal/explore/templates.go](/internal/explore/templates.go)

### Currently logged:
`registry`, `namespace`, `repo`, are all logged correctly.

### Incorrect:
`tag` - column exists empty
`image_ref` - **incorrectly omits the tag**
    `image_ref` = `namespace/repo:tag`