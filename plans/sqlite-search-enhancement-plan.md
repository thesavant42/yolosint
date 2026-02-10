# SQLite Search Enhancement Plan

## Overview

Enhance the SQLite logging to include image metadata (registry, namespace, repository, tag) to enable rich search capabilities.

## Current State Analysis

### Data Flow
1. [`explore.go:150-156`](../internal/explore/explore.go:150) - `logTOC(key string, toc *soci.TOC)` callback receives only the digest key
2. [`soci.go:82-84`](../internal/explore/soci.go:82) and [`soci.go:245-247`](../internal/explore/soci.go:245) - Sets `idx.OnTOC = h.logTOC` with just the key
3. [`sqlitedb.go:67-91`](../internal/explore/sqlitedb.go:67) - Schema stores `digest` in layers table but lacks image context

### Key Insight
In both [`tryNewIndex`](../internal/explore/soci.go:32) and [`createIndex`](../internal/explore/soci.go:228), we have access to `dig name.Digest` which contains:
- `dig.Context().String()` - full repo path like `ghcr.io/owner/repo`
- `dig.Context().RegistryStr()` - registry like `ghcr.io`
- `dig.Context().RepositoryStr()` - repository path like `owner/repo`
- `dig.Identifier()` - the digest hash

---

## Proposed Schema Enhancement

### Current layers table
```sql
CREATE TABLE IF NOT EXISTS layers (
    id INTEGER PRIMARY KEY,
    digest TEXT NOT NULL,
    csize INTEGER,
    usize INTEGER,
    type TEXT,
    media_type TEXT,
    indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Enhanced layers table
```sql
CREATE TABLE IF NOT EXISTS layers (
    id INTEGER PRIMARY KEY,
    digest TEXT NOT NULL,
    csize INTEGER,
    usize INTEGER,
    type TEXT,
    media_type TEXT,
    indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- NEW COLUMNS
    registry TEXT,      -- e.g., ghcr.io, docker.io
    namespace TEXT,     -- e.g., chainguard, library
    repository TEXT,    -- e.g., nginx, go
    tag TEXT,           -- e.g., latest, v1.0.0
    image_ref TEXT      -- full reference: ghcr.io/chainguard/nginx:latest
);

-- Add indexes for search
CREATE INDEX IF NOT EXISTS idx_layers_registry ON layers(registry);
CREATE INDEX IF NOT EXISTS idx_layers_namespace ON layers(namespace);
CREATE INDEX IF NOT EXISTS idx_layers_repository ON layers(repository);
CREATE INDEX IF NOT EXISTS idx_layers_image_ref ON layers(image_ref);
```

---

## Code Changes

### 1. sqlitedb.go - Add ImageContext struct and update schema

```go
// ImageContext holds the parsed image reference metadata
type ImageContext struct {
    Registry   string // e.g., ghcr.io
    Namespace  string // e.g., chainguard (first path segment)
    Repository string // e.g., nginx (remaining path)
    Tag        string // e.g., latest
    ImageRef   string // full reference
}
```

**Update schema** in `init()` function to include new columns and indexes.

**Update Insert signature:**
```go
func (t *TocDB) Insert(digest string, toc *soci.TOC, imgCtx *ImageContext) error
```

**Update INSERT statement:**
```go
res, err := tx.Exec(
    `INSERT INTO layers (digest, csize, usize, type, media_type, registry, namespace, repository, tag, image_ref) 
     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    digest, toc.Csize, toc.Usize, toc.Type, toc.MediaType,
    imgCtx.Registry, imgCtx.Namespace, imgCtx.Repository, imgCtx.Tag, imgCtx.ImageRef,
)
```

### 2. explore.go - Update logTOC callback signature

**Change from:**
```go
func (h *handler) logTOC(key string, toc *soci.TOC) {
    if h.tocDB != nil {
        if err := h.tocDB.Insert(key, toc); err != nil {
            log.Printf("SQLite insert failed for %s: %v", key, err)
        }
    }
}
```

**To:**
```go
func (h *handler) logTOC(key string, toc *soci.TOC, imgCtx *ImageContext) {
    if h.tocDB != nil {
        if err := h.tocDB.Insert(key, toc, imgCtx); err != nil {
            log.Printf("SQLite insert failed for %s: %v", key, err)
        }
    }
}
```

### 3. soci.go - Pass ImageContext when calling OnTOC

**In tryNewIndex around line 82-84:**
```go
// Extract image context from dig
imgCtx := extractImageContext(dig)

idx.Key = key
idx.OnTOC = func(k string, t *soci.TOC) {
    h.logTOC(k, t, imgCtx)
}
```

**In createIndex around line 245-247:**
```go
// Extract image context - note: createIndex receives prefix not dig
// Need to modify signature or pass context differently
idx.Key = key
idx.OnTOC = func(k string, t *soci.TOC) {
    h.logTOC(k, t, imgCtx)
}
```

**Add helper function:**
```go
func extractImageContext(dig name.Digest) *ImageContext {
    ctx := dig.Context()
    repoStr := ctx.RepositoryStr()
    
    // Parse namespace and repository from path
    // e.g., "chainguard/nginx" -> namespace="chainguard", repo="nginx"
    parts := strings.SplitN(repoStr, "/", 2)
    namespace := ""
    repository := repoStr
    if len(parts) == 2 {
        namespace = parts[0]
        repository = parts[1]
    }
    
    return &ImageContext{
        Registry:   ctx.RegistryStr(),
        Namespace:  namespace,
        Repository: repository,
        Tag:        "", // Tag not available from Digest, only from Tag reference
        ImageRef:   ctx.String(),
    }
}
```

---

## Files to Modify

| File | Line Range | Change Description |
|------|------------|-------------------|
| `internal/explore/sqlitedb.go` | 16-22 | Add ImageContext struct |
| `internal/explore/sqlitedb.go` | 67-91 | Update schema with new columns and indexes |
| `internal/explore/sqlitedb.go` | 115 | Update Insert signature to accept ImageContext |
| `internal/explore/sqlitedb.go` | 132-135 | Update INSERT SQL to include new columns |
| `internal/explore/explore.go` | 149-156 | Update logTOC signature and call |
| `internal/explore/soci.go` | 82-84 | Extract ImageContext and wrap callback |
| `internal/explore/soci.go` | 245-247 | Extract ImageContext and wrap callback |
| `internal/explore/soci.go` | new | Add extractImageContext helper function |

---

## Search Capabilities Enabled

With this enhanced schema, users can:

1. **Search files by name pattern** and see which registry/repo/image they came from
   ```sql
   SELECT f.name, l.registry, l.namespace, l.repository 
   FROM files f JOIN layers l ON f.layer_id = l.id 
   WHERE f.name LIKE '%passwd%';
   ```

2. **Search by image reference** and list all files
   ```sql
   SELECT f.* FROM files f 
   JOIN layers l ON f.layer_id = l.id 
   WHERE l.image_ref LIKE '%nginx%';
   ```

3. **Search for content from a specific repository**
   ```sql
   SELECT * FROM layers WHERE repository = 'nginx';
   ```

4. **Search for a specific namespace**
   ```sql
   SELECT * FROM layers WHERE namespace = 'chainguard';
   ```

5. **Search broadly across all fields**
   ```sql
   SELECT * FROM layers 
   WHERE registry LIKE '%ghcr%' 
   OR namespace LIKE '%chainguard%' 
   OR repository LIKE '%go%';
   ```

---

## Migration Strategy

For existing databases, add ALTER TABLE statements:
```sql
ALTER TABLE layers ADD COLUMN registry TEXT;
ALTER TABLE layers ADD COLUMN namespace TEXT;
ALTER TABLE layers ADD COLUMN repository TEXT;
ALTER TABLE layers ADD COLUMN tag TEXT;
ALTER TABLE layers ADD COLUMN image_ref TEXT;
```

The schema uses `CREATE TABLE IF NOT EXISTS` so new columns need to be handled gracefully for existing installations.
