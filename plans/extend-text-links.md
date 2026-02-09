# Plan: Extend Text Links in Image Main Route

## Summary

This plan adds new icon and text link elements to the image main route display, specifically in the [`renderManifestTables`](internal/explore/render.go:362) function.

## Goals from plans/goal.md

1. Add "build history" icon and text link after the config link
2. Make the "Config" icon hyperlink include the text label (increase clickable surface area)  
3. Add `mdi--docker.png` to the embedded binary routes

## Files to Modify

### 1. internal/explore/explore.go

Add route handlers to serve the new icons. These follow the existing pattern around lines 205-232.

**Add handlers for:**
- `/material-symbols--network-intelligence-history.png` - for the history icon
- `/mdi--docker.png` - Docker icon (goal #3)

**Location:** After the existing icon handlers (around line 232, before the robots.txt handler)

**Pattern to follow:**
```go
if r.URL.Path == "/material-symbols--network-intelligence-history.png" {
    w.Header().Set("Cache-Control", "max-age=3600")
    data, _ := Assets.ReadFile("assets/material-symbols--network-intelligence-history.png")
    w.Header().Set("Content-Type", "image/png")
    w.Write(data)
    return
}
if r.URL.Path == "/mdi--docker.png" {
    w.Header().Set("Cache-Control", "max-age=3600")
    data, _ := Assets.ReadFile("assets/mdi--docker.png")
    w.Header().Set("Content-Type", "image/png")
    w.Write(data)
    return
}
```

### 2. internal/explore/render.go

Modify the [`renderManifestTables`](internal/explore/render.go:362) function.

#### Change 1: Extend Config Link (Goal #2)

**Current code (line 403):**
```go
w.Printf(`<a href="/%s%s@%s%smt=%s&size=%d" title="Config: %s"><img src="/eos-icons--init-container-outlined.png" alt="config" style="height:16px;vertical-align:middle"/></a> config`,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest))
```

**New code:**
```go
w.Printf(`<a href="/%s%s@%s%smt=%s&size=%d" title="Config: %s"><img src="/eos-icons--init-container-outlined.png" alt="config" style="height:16px;vertical-align:middle"/> config</a>`,
    handler, w.repo, digest, qs, url.QueryEscape(mt), size, html.EscapeString(digest))
```

The change: Move ` config` inside the `</a>` closing tag.

#### Change 2: Add Build History Link (Goal #1)

**Current code (line 408):**
```go
// Combined layers link with icon (same row as config)
w.Print(` <a href="/layers/` + image + `/"><img src="/f7--layers-alt-fill.png" alt="layers" style="height:16px;vertical-align:middle"/></a><a href="/layers/` + image + `/">combined layers view</a>`)
```

**After this line, add the build history link:**
```go
// Build history link with icon
u := *w.u
qs := u.Query()
qs.Set("render", "history")
qs.Set("mt", mt)
u.RawQuery = qs.Encode()
w.Printf(` <a href="%s"><img src="/material-symbols--network-intelligence-history.png" alt="history" style="height:16px;vertical-align:middle"/> build history</a>`, u.String())
```

This uses the same URL construction pattern as the existing [`History`](internal/explore/render.go:112) method.

## Visual Result

The icon row will display:
```
[manifest-icon] v2+json [config-icon] config [layers-icon] combined layers view [history-icon] build history
```

Where:
- `[config-icon] config` is now a single clickable link (goal #2)
- `[history-icon] build history` is the new addition (goal #1)

## Asset Verification

The required icons already exist in `internal/explore/assets/`:
- `material-symbols--network-intelligence-history.png` - history icon
- `mdi--docker.png` - Docker icon
- `eos-icons--init-container-outlined.png` - config icon (already in use)

## Testing

After implementation, verify by:
1. Loading an image config page (e.g., the URL from goal.md)
2. Confirming the config link now includes the text label
3. Confirming the build history icon and text appear after combined layers view
4. Clicking build history link navigates to the history render view
5. Verifying the new icons load correctly at their URL paths
