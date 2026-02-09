# Plan: Column Reorder and Docker Pull Removal

## Summary

Two changes requested:
1. Reorder columns in the layers table: move file size column AFTER digests and BEFORE download links
2. Remove the docker pull box entirely (it was broken and showing sha256 digest instead of a usable reference)

---

## Issue 1: Reorder Layers Table Columns

### Current Column Order
```
| Index | Size | Digest | Download |
```

### Requested Column Order
```
| Index | Digest | Size | Download |
```

### Files to Modify

**File: [`internal/explore/render.go`](internal/explore/render.go:408)**

#### Change 1: Update header row (line 408)

Current:
```go
w.Print(`<table><tr><td colspan="2"><strong>[LIST VIEW] </strong></td><td colspan="2"><strong>[LAYERS VIEW] </strong> [<a href="/layers/` + image + `/">COMBINED LAYERS VIEW</a>]</td><td></td></tr>`)
```

Change to:
```go
w.Print(`<table><tr><td></td><td><strong>[LAYERS VIEW] </strong> [<a href="/layers/` + image + `/">COMBINED LAYERS VIEW</a>]</td><td><strong>[LIST VIEW]</strong></td><td></td></tr>`)
```

This puts:
- Column 1: Index (empty header)
- Column 2: LAYERS VIEW label (over digest column)
- Column 3: LIST VIEW label (over size column)
- Column 4: Download (empty header)

#### Change 2: Update data row (lines 436-440)

Current order in Printf:
```go
w.Printf(`<tr><td>%d</td><td><a href="/size/...">%s</a></td><td><a href="/%s...">%s</a></td><td>...</td></tr>`,
    i+1,                                    // index
    ..., humanize.IBytes(uint64(size)),     // size
    ..., html.EscapeString(digest),         // digest
    ...)                                    // download
```

Change to:
```go
w.Printf(`<tr><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d">%s</a></td><td><a href="/size/%s@%s?mt=%s&size=%d">%s</a></td><td>...</td></tr>`,
    i+1,                                    // index
    ..., html.EscapeString(digest),         // digest (now column 2)
    ..., humanize.IBytes(uint64(size)),     // size (now column 3)
    ...)                                    // download
```

---

## Issue 2: Remove Docker Pull Box

### Files to Modify

**File: [`internal/explore/templates.go`](internal/explore/templates.go:258)**

Remove lines 258 and 261 which render the docker pull box:
```go
{{if .DockerPull}}<p><span style="background:#2a2a3e;padding:4px 8px;border-radius:4px;border:1px solid #444;">docker pull {{.DockerPull}} <button onclick="navigator.clipboard.writeText('docker pull {{.DockerPull}}')" style="background:#444;color:inherit;border:1px solid #666;padding:2px 6px;border-radius:4px;cursor:pointer;margin-left:4px;font:inherit;">Copy</button></span></p>{{end}}
```

**File: [`internal/explore/templates.go`](internal/explore/templates.go:310)**

Remove the `DockerPull` field from the `HeaderData` struct:
```go
DockerPull           string  // Remove this line
```

**File: [`internal/explore/explore.go`](internal/explore/explore.go:1560)**

Remove the line that sets the DockerPull field:
```go
header.DockerPull = ref.String()  // Remove this line
```

---

## Implementation Checklist

- [ ] Update header row in `renderManifestTables` to reorder column labels
- [ ] Update data row Printf in `renderManifestTables` to swap digest and size columns
- [ ] Remove DockerPull template rendering from `bodyTemplate` (2 occurrences)
- [ ] Remove `DockerPull` field from `HeaderData` struct
- [ ] Remove `header.DockerPull = ref.String()` assignment in `manifestHeader`
