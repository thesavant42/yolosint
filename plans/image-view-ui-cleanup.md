# Image View UI Cleanup Plan

## Overview
Simplify the image manifest display by replacing verbose text with icons and removing redundant information.

## Current State (from mockup)
```
mediaType application/vnd.docker.distribution.manifest.v2+json
digest    6104 sha256:cdba5cff38087cbac9d4f4c123b14a4df1a51db514285961c8429499fdc1d401
CONFIG: 5.4 KiB sha256:d7c8d5a9296211e607ca986f18de61993d2254201b2f54c788783b273e24182a
```

## Desired State
- **mediaType row**: Icon + abbreviated label (e.g., "v2" for manifest.v2+json)
- **digest row**: Remove entirely
- **CONFIG row**: Replace with just a clickable icon linking to the config

---

## Changes Required

### 1. Add Asset Routes in explore.go

The icons already exist in `internal/explore/assets/`:
- `ant-design--container-outlined.png` - for manifest type indicator
- `eos-icons--init-container-outlined.png` - for config link

Add HTTP routes similar to existing `/gis--layer-download.png` route.

**File:** `internal/explore/explore.go`

Add handlers for both icons (follow existing pattern for gis--layer-download.png).

### 2. Update templates.go - Manifest Display

**File:** `internal/explore/templates.go`
**Location:** Lines 263-265

**Current:**
```
mediaType application/vnd.docker.distribution.manifest.v2+json
digest    6104 sha256:...
```

**New:**
```
[icon] v2
```

Just a single row with the container icon followed by the abbreviated version.

**Implementation:**
```html
<tr><td><img src="/ant-design--container-outlined.png" alt="manifest" style="height:16px;vertical-align:middle"/> {{abbreviateMediaType .Descriptor.MediaType}}</td></tr>
```

Changes:
- Replace entire "mediaType ..." row with icon + short label
- Remove the entire digest row
- Need a template function `abbreviateMediaType` to extract version:
  - `application/vnd.docker.distribution.manifest.v2+json` -> "v2"
  - `application/vnd.oci.image.manifest.v1+json` -> "OCI"
  - `application/vnd.docker.distribution.manifest.list.v2+json` -> "list v2"
  - `application/vnd.oci.image.index.v1+json` -> "index"

### 3. Update render.go - Config Display

**File:** `internal/explore/render.go`
**Location:** Lines 374-397

**Current:**
```
CONFIG:	5.4 KiB	sha256:d7c8d5a9296211e607ca986f18de61993d2254201b2f54c788783b273e24182a
```

**New:**
```
[config icon](hyperlink to config)
```

Just a clickable icon - no text, no size, no digest visible.

**Implementation:**
```go
w.Print(`<table><tr><td>`)
w.Printf(`<a href="/%s%s@%s%smt=%s&size=%d" title="Config: %s"><img src="/eos-icons--init-container-outlined.png" alt="config" style="height:16px;vertical-align:middle"/></a>`,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest))
w.Print(`</td></tr></table>`)
```

Changes:
- Remove "CONFIG:" text label
- Remove size display (5.4 KiB)
- Remove visible digest hash
- Keep only a clickable icon that links to the config (digest in title for tooltip on hover)

---

## Files to Modify

| File | Change |
|------|--------|
| `internal/explore/explore.go` | Add 2 new asset routes |
| `internal/explore/templates.go` | Replace mediaType with icon + short label, remove digest row |
| `internal/explore/render.go` | Replace CONFIG section with icon only |

## Assets Already Present
- `internal/explore/assets/ant-design--container-outlined.png`
- `internal/explore/assets/eos-icons--init-container-outlined.png`

## MediaType Abbreviation Examples
| Full MediaType | Abbreviated |
|----------------|-------------|
| `application/vnd.docker.distribution.manifest.v2+json` | v2 |
| `application/vnd.docker.distribution.manifest.v1+json` | v1 |
| `application/vnd.oci.image.manifest.v1+json` | OCI |
| `application/vnd.docker.distribution.manifest.list.v2+json` | list v2 |
| `application/vnd.oci.image.index.v1+json` | index |
