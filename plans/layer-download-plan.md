# Layer Download Feature Implementation Plan

## Overview

Add download links `[x]` in a new column after each layer's sha256 digest in the manifest view (as shown in mockup).

## Mockup Reference

The `[x]` link appears as a 5th column after the digest:

```
| idx | size      | idx | sha256:digest...                                      | [x] |
| 1   | 2.3 MiB   | 1   | sha256:49388a8bc9c86a6f56d228954eede699c64fce6c671... | [x] |
| 2   | 1.2 MiB   | 2   | sha256:8e5f7c337501eeff98e19cc2d10f4c50fd8bfcce7fe... | [x] |
...
```

## Implementation - Single File Change

**File:** [`internal/explore/render.go`](internal/explore/render.go)

**Function:** [`renderManifestTables()`](internal/explore/render.go:362)

### Change 1: Add filename base parsing (after line 363)

```go
image := w.u.Query().Get("image")

// Build filename base from image reference (replace special chars with -)
filenameBase := strings.ReplaceAll(image, "/", "-")
filenameBase = strings.ReplaceAll(filenameBase, ":", "-")
filenameBase = strings.ReplaceAll(filenameBase, "@", "-")
```

### Change 2: Update table header (line 391)

Current:
```go
w.Print(`<table><tr><td colspan="2"><strong>LIST VIEW:</strong></td><td colspan="2"><strong>LAYERS VIEW</strong> [<a href="/layers/` + image + `/">COMBINED LAYERS VIEW</a>]</td></tr>`)
```

Change to (add empty 5th column):
```go
w.Print(`<table><tr><td colspan="2"><strong>LIST VIEW:</strong></td><td colspan="2"><strong>LAYERS VIEW</strong> [<a href="/layers/` + image + `/">COMBINED LAYERS VIEW</a>]</td><td></td></tr>`)
```

### Change 3: Add download link column to each layer row (lines 412-418)

Current Printf (line 413-417):
```go
w.Printf(`<tr><td>%d</td><td><a href="/size/%s@%s?mt=%s&size=%d">%s</a></td><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d">%s</a></td></tr>`,
    i+1,
    w.repo, digest, url.QueryEscape(mt), size, humanize.IBytes(uint64(size)),
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest))
```

Change to (add 5th column with `[x]` download link):
```go
// Construct download filename: namespace-repo-tag-idx.tar.gzip
downloadFilename := fmt.Sprintf("%s-%d.tar.gzip", filenameBase, i+1)

w.Printf(`<tr><td>%d</td><td><a href="/size/%s@%s?mt=%s&size=%d">%s</a></td><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d">%s</a></td><td><a href="/blob/%s@%s?mt=%s&size=%d" download="%s" title="Download %s">[x]</a></td></tr>`,
    i+1,
    w.repo, digest, url.QueryEscape(mt), size, humanize.IBytes(uint64(size)),
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest),
    w.repo, digest, url.QueryEscape(mt), size, html.EscapeString(downloadFilename), html.EscapeString(downloadFilename))
```

## Summary of Changes

| Line | Change |
|------|--------|
| 363-367 | Add `filenameBase` variable construction after `image` |
| 391 | Add empty `<td></td>` to table header for download column |
| 412-418 | Add `downloadFilename` variable and add 5th `<td>` with `[x]` link |

## Key Points

1. **No new routes needed** - Uses existing `/blob/` endpoint with HTML `download` attribute
2. **Single file change** - Only `internal/explore/render.go` needs modification  
3. **Filename format** - `namespace-repo-tag-idx.tar.gzip` (e.g., `smoshysmosh-composer-latest-1.tar.gzip`)
4. **`download` attribute** - Tells browser to download with custom filename instead of navigating
