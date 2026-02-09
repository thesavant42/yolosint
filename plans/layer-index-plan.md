# Layer Index Display Plan

## Requirements (from goal.md)

1. Display the current layer number above the path when viewing a single layer
2. Include a download link for that layer using `/gis--layer-download.png` icon (same link rendered in /image? root view)
3. For merged view: show label indicating "Merged View" with `--` instead of download link

## Implementation Details

### 1. templates.go - Add fields to HeaderData (line 293)

**WHY:** The `HeaderData` struct is the data passed to the HTML template. Currently it has no way to know which layer is being viewed or whether it's a merged view. We need to add fields so the template can conditionally render the layer information.

**WHAT:** Add three new fields to the `HeaderData` struct:

```go
type HeaderData struct {
    // ... existing fields ...
    Path                 string
    LayerIndex           int    // NEW: 1-based layer index (0 means not set)
    LayerDownloadURL     string // NEW: download URL for this layer
    IsMergedView         bool   // NEW: true when viewing merged layers
}
```

### 2. templates.go - Update bodyTemplate (line 262)

**WHY:** The `bodyTemplate` renders the page header. Currently it shows `path: /some/path` but has no layer information. Per the mockup in goal.md, the layer info should appear ABOVE the path line.

**WHAT:** Currently shows:
```
{{if .Path}}<p>path: {{.Path}}</p>{{end}}
```

Change to display layer info ABOVE the path:
```html
{{if .IsMergedView}}<p>Merged View --</p>{{else if .LayerIndex}}<p>Layer {{.LayerIndex}} <a href="{{.LayerDownloadURL}}"><img src="/gis--layer-download.png" alt="Download" style="height:16px;vertical-align:middle"/></a> <a href="{{.LayerDownloadURL}}">Download Layer</a></p>{{end}}
{{if .Path}}<p>path: {{.Path}}</p>{{end}}
```

### 3. render.go - Add layer param to links (line 443)

**WHY:** When the user is on the `/image?` view, they see a list of layers with their index numbers (1, 2, 3...) and download links. When they click a layer digest to browse its filesystem, that layer index information is currently LOST because the URL `/fs/repo@sha256:abc` has no layer context. We need to pass the layer index and download URL as query parameters so they survive the navigation.

**WHAT:** In `renderManifestTables`, the layer links are generated at line 443. The `downloadURL` variable already exists at line 441 - this is the same download link currently shown next to each layer in the manifest view. We will pass it forward.

Current code at line 443:
```go
w.Printf(`<tr><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d">%s</a></td>...`,
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest),
    ...
```

Change to include `layer` and `dlurl` query params:
```go
w.Printf(`<tr><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d&layer=%d&dlurl=%s">%s</a></td>...`,
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, i+1, url.QueryEscape(downloadURL), html.EscapeString(digest),
    ...
```

### 4. explore.go - Read layer params in renderHeader (line 1679)

**WHY:** The `renderHeader` function is called when viewing a single layer's filesystem. It creates the `HeaderData` that gets passed to the template. We need to read the `layer` and `dlurl` query parameters that were added in step 3 and populate the new `HeaderData` fields so the template can render them.

**WHAT:** After line 1734 where `header.Path` is set:

```go
header.Path = currentPath

// NEW: Read layer index and download URL from query params
if layerStr := r.URL.Query().Get("layer"); layerStr != "" {
    if layerIdx, err := strconv.Atoi(layerStr); err == nil {
        header.LayerIndex = layerIdx
    }
}
if dlurl := r.URL.Query().Get("dlurl"); dlurl != "" {
    header.LayerDownloadURL = dlurl
}
```

### 5. explore.go - Set IsMergedView in renderDir (line 1773)

**WHY:** The `renderDir` function is called when viewing the merged/combined layers view. Per the requirements, the merged view should show "Merged View --" instead of a layer number and download link. We need to set a flag so the template knows to render the merged view text.

**WHAT:** After line 1834 where `header.Path` is set:

```go
header.Path = currentPath

// NEW: Mark this as merged view
header.IsMergedView = true
```

## Data Flow Example

### Single Layer View

**Starting point:** User is viewing `/?image=cgr.dev/chainguard/static:latest`

The manifest view shows a table of layers:

| # | Digest | Size | Download |
|---|--------|------|----------|
| 1 | sha256:abc123... | 2.1 MiB | [icon] |
| 2 | sha256:def456... | 1.5 MiB | [icon] |
| 3 | sha256:789xyz... | 500 KiB | [icon] |

**Current behavior:** Clicking layer 2's digest links to:
```
/fs/cgr.dev/chainguard/static@sha256:def456...?mt=...&size=...
```
The layer number (2) is lost.

**New behavior:** Clicking layer 2's digest links to:
```
/fs/cgr.dev/chainguard/static@sha256:def456...?mt=...&size=...&layer=2&dlurl=%2Fdownload%2F...
```
The layer number (2) and the download URL are preserved in the query string.

When the filesystem view renders, it reads `layer=2` and displays:
```
Layer 2 [download-icon] Download Layer
path: /
```

### Merged View

**Starting point:** User clicks "combined layers view" from manifest page

User navigates to `/layers/cgr.dev/chainguard/static:latest/`

The `renderDir` function handles this view and sets `IsMergedView=true`.

The template displays:
```
Merged View --
path: /
```

No download link is shown because merged view represents all layers combined.
