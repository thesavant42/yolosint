# Refactor Image View Plan

## 4 Requirements Mapped to Code Changes


### 1. `mediaType...` shortened to `Manifest v2`
**Target:** Show "manifest.v2" instead of full mediaType string
**File:** [`templates.go:264`](../internal/explore/templates.go:264)
**Change:** Replace `{{.Descriptor.MediaType}}` with shortened text

### 2. `[LIST VIEW]` extends past first column, no wrap
**Target:** Label spans columns and does not wrap
**File:** [`render.go:400`](../internal/explore/render.go:400)
**Change:** Add `style="white-space:nowrap"` and proper colspan

### 3. `CONFIG` capitalized, aligned with row numbers
**Target:** "CONFIG" in uppercase, same alignment as layer indices
**File:** [`render.go:375`](../internal/explore/render.go:375)
**Change:** Change `config` to `CONFIG`

### 4. `[COMBINED LAYERS]` moved to manifest row, before referrers
**Target:** COMBINED LAYERS link appears after manifest.v2, before (referrers)
**File:** [`templates.go:264`](../internal/explore/templates.go:264) - add link
**File:** [`render.go:400`](../internal/explore/render.go:400) - remove from here

---

## Files to Modify

| File | Lines | Changes |
|------|-------|---------|
| `internal/explore/templates.go` | 264 | Shorten mediaType, add COMBINED LAYERS link |
| `internal/explore/render.go` | 375, 400, 429-434 | CONFIG caps, LIST VIEW styling, 4-column layer rows |

---

## Acceptance Criteria
1. Must look like the mockup screenshot
2. All links must continue to work: COMBINED LAYERS, CONFIG, LIST VIEW, LAYERS VIEW, SAVE LAYERS
