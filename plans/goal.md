# Save Single Layer from Layer view


## Current Status

![current view](/plans/current.png)

### User Journey
- When I visit a `?image=namespace/repo:tag` root link, such as http://192.168.1.37:8042/?image=smoshysmosh%2Fcomposer%3Alatest :
    - I am able to view a [list of the fies](http://192.168.1.37:8042/size/smoshysmosh/composer@sha256:0b342bf7e74f5fbef621a6d413a2a56088d405057e0e165e383904fe37fe28a2?mt=application%2Fvnd.docker.image.rootfs.diff.tar.gzip&size=26470178),
    - I am able to browse the [filesystem by clicking this link](http://192.168.1.37:8042/fs/smoshysmosh/composer@sha256:0b342bf7e74f5fbef621a6d413a2a56088d405057e0e165e383904fe37fe28a2/?mt=application%2Fvnd.docker.image.rootfs.diff.tar.gzip&size=26470178)

### Problem Statement 

- There's no mechanism to download the layer `.tar.gzip` file as a whole layer

## Proposed Solution

Add a text link to download the layer from the registry, as indicated in this screenshot mockup: The text label `[x]` is after each layer's `sha256:digest`. 

![mockup](/plans/mockup.png)

Instead of the default sha256 layer digest as the file name I'd like to save it as:

`namespace-repo-tag-idx.tar.gzip`. where `idx` is the layer number being saved.

---

You hit Docker Hub’s registry API directly: first GET a Bearer token for `repository:<namespace>/<repo>:pull`, then use that token to fetch the manifest, then download the specific layer blob by its digest. The blob URL is always `https://registry-1.docker.io/v2/<namespace>/<repo>/blobs/<digest>` and returns the raw tar.gz filesystem layer. 
---




Based on my analysis of the codebase, I now understand the requirements and the current architecture. Let me create a plan for implementing the layer download feature.

## Summary of Analysis

**Current State:**
- The layer list is rendered in [`renderManifestTables()`](internal/explore/render.go:362) in `render.go`
- Each layer row shows: index, size link, index again, and the digest link
- There's an existing `SaveURL` pattern used for file downloads in [`HeaderData`](internal/explore/templates.go:313)

**Goal:**
- Add `[x]` download links after each layer's digest
- Download filename format: `namespace-repo-tag-idx.tar.gzip`