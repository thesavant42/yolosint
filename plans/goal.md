# sqlite db refactor to include 


 **Failure to follow any of these steps will result in immediate termination.**

- Read files in order, do not read files I do not specifically mention
- If you get blocked you must pause and ask the user for guidance
- You must communicate with the user via Chat.
- If the user asks you a question, you MUST respond in chat immediately, not proceed without permission.

## Problem Statement

- When I am viewing a layer, I would like to know what layer I am viewing (IDX)
- When I search for strings, I want to know the namespace, owner, and repository from which they originate
- Current SQLite DB structure is not yet robust enough of a schema 

### Goal

Goal: enhance the logging in the sqlite database to include additional metadata attributes to make the content more searchable and with richer details.

- [explore.go](/internal/explore/explore.go) is the integration point for sqlite. If you try to read `sqlite.go` first you will be terminated.
    - Line 149-156: logTOC callback receives key (which is the digest like sha256:abc123.0) and the *soci.TOC
    - Line 113: NewTocDB("/cache/log.db") initializes the SQLite DB
    - The logTOC function is set as a callback in soci.go at lines 83-84 and 246-247
- [sqlite.go](/internal/explore/sqlitedb.go)
    - layers table: id, digest, csize, usize, type, media_type, indexed_at
    - files table: id, layer_id, name, typeflag, size, mode, mod, offset, linkname
    - The Insert method only receives digest and *soci.TOC - no image reference context
- [soci.go](/internal/explore/soci.go)
    - Line 83-84 in `tryNewIndex: idx.Key = key; idx.OnTOC = h.logTOC` - here we have access to `dig name.Digest` which has:
        - `dig.Context().String()` - full repo path like `ghcr.io/owner/repo`
        - `dig.Context().RegistryStr()` - registry like `ghcr.io`
        - `dig.Context().RepositoryStr()` - repo part
        - `dig.Identifier()` - the sha256 digest
    - Line 246-247 in createIndex: same pattern
- [Sqlite - log.db](/cache/log.db) - use mcp


- The problem is clear: when logTOC is called, it only receives the key (digest-based) and toc.
- The rich context from name.Digest is available at the call sites but not passed through.

### Proposed Solution

- Modify logTOC signature to accept additional image context:
    - registry
    - namespace (owner)
    - repository
    - tag
    - full reference string

- Modify the layers table schema to add:
    registry TEXT - e.g., "ghcr.io"
    namespace TEXT - e.g., "owner" or "org"
    repository TEXT - e.g., "repo-name"
    tag TEXT - e.g., "latest" or "v1.0"
    reference TEXT - full reference string

- Modify Insert method in TocDB to accept these new fields alongside the digest and TOC data, then store them in the layers table when recording indexed content.

- Update call sites in soci.go to pass the additional context at tryNewIndex (line 83-84) and createIndex (line 246-247).

- Add search functionality as a follow-up enhancement that leverages the enriched metadata now available in the schema.

---

log.db -> Files table Columns:
    - id                - autoincrement default row id for sqlite
    - layer_id          - the layer IDX - This is a human friendly mapping of the layer sha256 digest
                        - in a 10 layer container, the last layer would be idx[9], for example.
    - name              - the file or directory path
    - typeflag          - indicates if above is a file, directory, symlink, etc
    - size              - size of object
    - mode              - ?
    - mod               - timestamp?
    - offset            - location of file within the .tar.gzip image
    - linkname          - symlink destination

Not Present:
    - Layer's `sha256 digest` is not mapped to the file
    - image's `tag` 
    - `repository` 
    - `namespace`

---

1. **Current state**: The `layers` table stores just the digest (like `sha256:abc123.0`) but loses the full image reference context (registry/namespace/repo:tag). The `logTOC` callback only receives this digest key.

2. **What I see available**: 
    1. When indexing happens (in `tryNewIndex`, `createIndex`), we have access to `name.Digest` which contains `dig.Context().String()` (full repo path like `ghcr.io/owner/repo`) and `dig.Context().RegistryStr()`.


---

## Search Usecases
- Search for files by name pattern and see which registry/repo/image they came from?
- Search by image reference and list all files?
- Search for content from a specific repository
- Search for a specific namespace
- Search broadly for a term across all fields
