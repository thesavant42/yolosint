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

**WHAT:** In `renderManifestTables`, the layer links are generated at line 443:

```go
w.Printf(`<tr><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d">%s</a></td>...`,
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest),
    ...
```

Change the `/fs/` link to include `layer` and `dlurl` query params:

```go
// Build download URL for this layer (already exists at line 441)
downloadURL := fmt.Sprintf("/download/%s@%s?filename=%s", w.repo, digest, url.QueryEscape(downloadFilename))

// Add layer index and download URL to the /fs/ link
w.Printf(`<tr><td>%d</td><td><a href="/%s%s@%s%smt=%s&size=%d&layer=%d&dlurl=%s">%s</a></td>...`,
    i+1,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, i+1, url.QueryEscape(downloadURL), html.EscapeString(digest),
    ...
```

### 4. explore.go - Read layer params in renderHeader (line 1679)

**WHY:** The `renderHeader` function is called when viewing a single layer's filesystem (e.g., `/fs/repo@sha256:abc/some/path`). It creates the `HeaderData` that gets passed to the template. We need to read the `layer` and `dlurl` query parameters that were added in step 3 and populate the new `HeaderData` fields so the template can render them.

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

**WHY:** The `renderDir` function is called when viewing the merged/combined layers view (e.g., `/layers/repo:tag/`). This is different from viewing a single layer. Per the requirements, the merged view should show "Merged View --" instead of a layer number and download link. We need to set a flag so the template knows to render the merged view text.

**WHAT:** After line 1834 where `header.Path` is set:

```go
header.Path = currentPath

// NEW: Mark this as merged view
header.IsMergedView = true
```

## Data Flow

### Single Layer View:
1. User views manifest at `/?image=ubuntu:latest`
2. `renderManifestTables` generates layer list with links like `/fs/repo@sha256:abc?layer=3&dlurl=/download/repo@sha256:abc?filename=...`
3. User clicks layer 3
4. `renderHeader` reads `layer=3` and `dlurl=...` from query params, sets `header.LayerIndex=3` and `header.LayerDownloadURL=...`
5. `bodyTemplate` renders "Layer 3 [icon] Download Layer" above the path

### Merged View:
1. User views manifest at `/?image=ubuntu:latest`
2. User clicks "combined layers view" link to `/layers/ubuntu:latest/`
3. User navigates merged filesystem
4. `renderDir` sets `header.IsMergedView=true`
5. `bodyTemplate` renders "Merged View --" (no download link since this is merged)
