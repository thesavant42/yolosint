# Save Single Layer from Layer view


## Current Status

~~![current view](/plans/current.png)~~
**Has links, but they don't use the authentication flow**

1. Authenticate to the registry, 
2. use the JWT from the registry to make an authenticated request to save the layer
    - The registry serves the file, we save it.
    - The same flow exists for viewing layer images, saving layer files.


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
- The layer list is rendered in [`renderManifestTables()`](internal/explore/render.go:362) in `render.go`
- Each layer row shows: index, size link, index again, and the digest link
- There's an existing `SaveURL` pattern used for file downloads in [`HeaderData`](internal/explore/templates.go:313)

**Goal:**
- Add `[x]` download links after each layer's digest
- Download filename format: `namespace-repo-tag-idx.tar.gzip`

### Example code

Here's the sequence using a real image, but done in curl to demonstrate. 

```bash
# Get token
TOKEN=$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:smoshysmosh/composer:pull" | jq -r .token)
# Use token
curl -L -H "Authorization: Bearer $TOKEN" \
  "https://registry-1.docker.io/v2/smoshysmosh/composer/blobs/sha256:0b342bf7e74f5fbef621a6d413a2a56088d405057e0e165e383904fe37fe28a2" \
  -o namespace-repo-tag-layer.tar.gz
```


I see the problem clearly now. The current code at [`renderManifestTables()`](internal/explore/render.go:362) creates a direct `<a href>` link to the registry URL:

```go
downloadURL := fmt.Sprintf("https://%s/v2/%s/blobs/%s", registry, repoPath, digest)
```

This bypasses authentication entirely. The browser can't add Bearer tokens to a simple link click. The solution needs a server-side proxy endpoint that:
1. Receives the download request
2. Authenticates to the registry (gets Bearer token)
3. Fetches the blob
4. Streams it to the user with `Content-Disposition` header

---


## Analysis Summary

**The Problem:** The current code at [`renderManifestTables()`](internal/explore/render.go:362) creates direct `<a href>` links to registry URLs like `https://registry-1.docker.io/v2/<repo>/blobs/<digest>`. This bypasses authentication entirely - browsers cannot inject Bearer tokens into simple link clicks.

**The Solution:** Create a server-side proxy endpoint that:
1. Receives the download request from the user
2. Authenticates to the registry (gets Bearer token)
3. Fetches the blob with proper auth headers
4. Streams it back to the user with `Content-Disposition` header

**Existing Patterns:**
- [`SaveURL`](internal/explore/templates.go:313) pattern already exists for file downloads
- Handler has auth infrastructure: [`tokens map[string]token`](internal/explore/explore.go:66), keychain support
- Routes like [`/fs/`](internal/explore/explore.go:123), [`/size/`](internal/explore/explore.go:124) already fetch authenticated blobs
- 
- 
Now I see how it works. [`transportFromCookie()`](internal/explore/cookies.go:18) already:
- Pings registry (line 37)
- Gets/caches Bearer token via `transport.NewBearer()` (line 55)
- Token is stored in `tok.TokenResponse.Token` (cached at line 78)

The simplification would be to **extract** the token-getting logic into a separate `getToken()` function that `transportFromCookie()` can call, rather than duplicating code.