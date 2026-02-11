# columns not logging

## Problem statement:
The logging is not working correctly.

I want to be able to search by namespace, by repository, tag, in addition to the currently logged fields.

The issue is about logging to the database, not link generation. 

**You are not supposed to be "carrying forward" `tags` into other routes**
- **you are just supposed to be logging them. **


---

The app begins the user flow as such:

- http://localhost:8042/?image=robdisney/fauxpilotgov:latest


1. exits docker hub via browser plugin and fowwards to api from https://hub.docker.com/api/content/v1/entitlement/robdisney/features
2. Entry point into the applicationhttp://localhost:8042/?image=robdisney%2Ffauxpilotgov%3Alatest
    - Last place the TAG is available programatically via a route. 
    - Forwards to:
3. http://localhost:8042/fs/robdisney/fauxpilotgov@sha256:bb9475e9df95f20a223b74a1bd9401b28ef717d61afa1cb8177d72e60f4e61e8?mt=application%2Fvnd.docker.image.rootfs.diff.tar.gzip&size=19161232
    - The docker image digest, contains the name space and the repository but not the TAG

- this is received from a browser plugin
- This is the namespace, the repository, and the tag
- 
- The code currently looks for this data **incorrectly** in the `sha56:digest`, which will never ever have this mapping.
- The Build configuration typically has this mapping but is not stored in the sqlite database.
- current code only searches for image_ref, does not break down the image_ref into its component parts.



The tag is available at the entry point `/?image=robdisney/fauxpilotgov:latest` but when `extractImageContext()` is called from [`tryNewIndex()`](internal/explore/soci.go:58), it only has the digest reference which cannot recover the tag.



DBMS: modernc sqlite: https://gitlab.com/cznic/sqlite
db: cache/log.db
    - use MCP

- [/internal/explore/sqlitedb.go](/internal/explore/sqlitedb.go)
- [/internal/explore/soci.go](/internal/explore/soci.go)
- [/internal/explore/explore.go](/internal/explore/explore.go)
- [/internal/explore/templates.go](/internal/explore/templates.go)

### Currently logged:
`registry`, `namespace`, `repo`, are all logged correctly.

### Incorrect:
`tag` - column exists empty
`image_ref` - **incorrectly omits the tag**
    `image_ref` = `namespace/repo:tag`