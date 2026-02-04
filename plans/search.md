# Search Task 1 Log to sqlite

I want to add search to the application, so that I can search across the history for `EXAMPLE.md` and compare all instances of it, or any artbitrary search term.

## Current Architecture:

This is how the cache metadata is generated and stored:

- [`internal/soci/indexer.go`](internal/soci/indexer.go:22) - The `Indexer` struct parses tar files and creates a [`TOC`](internal/soci/toc.go:11) (Table of Contents) containing file metadata
- [`internal/explore/cache.go`](internal/explore/cache.go:23) - Stores the TOC as `toc.json.gz` and the indexed archive as `.tar.gz`

- The TOC already contains all the file metadata I want to search:
    -  `Name`,
    -  `Size`,
    -  `Offset`,
    -  `Linkname`,
    -  `ModTime`, etc.
    - **Index by layer digest** (tar.gz in cache has layer digest as name)

## Problem Statement

- I already have JSON data with all the file information, *but it's locked up in gzipped tar files*.
- I want to write this data to SQLite so I can actually search my history.

## Suggessted Solution

- Add SQLite storage that captures the TOC data during indexing, before/while writing the tar files.
- sqlite logging is **mandatory, not conditional. It should never be tied to the cache dir.**
- Do not utilize Environment Variables!

### Requirements

- Use [example toc.json](/cache/sha256-2b97a650489f286782ed9553c3c7e16669fe46b3aaacc3cd4944bfde910b002a.0/toc.json) for sample data to build schema
- Must work with Docker Compose: [docker-compose.yml](/docker-compose.yml) AND [docker-compose-synology](/docker-compose.synology.yml)
- Your **plan must be in markdown** as [/plans/plan.md](/plans/plan.md)


--- 

**Summary:**

The TOC data you need is already in [`soci.TOC`](internal/soci/toc.go:11) - filenames, sizes, timestamps, everything. The [`dirCache.Put()`](internal/explore/cache.go:178) method receives this data before writing to tar.

**The fix:** Add to [`internal/explore/cache.go`](internal/explore/cache.go:150)

---

## Sample Data toc.json

This is the beginning part of `toc.json` for `sha256-2b97a650489f286782ed9553c3c7e16669fe46b3aaacc3cd4944bfde910b002a.0.tar.gz` lines 1-34:

```json
{
    "csize": 36021143,
    "usize": 92976640,
    "ssize": 4194304,
    "type": "tar+gzip",
    "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
    "checkpoints": [
        {
            "in": 10,
            "empty": true
        },
        {
            "in": 1058203,
            "out": 4259840,
            "b": 4,
            "nb": 3,
            "wrpos": 374,
            "full": true
        },
        {
            "in": 3868005,
            "out": 8486912,
            "nb": 1,
            "wrpos": 31494,
            "full": true
        },
        {
            "in": 5388478,
            "out": 12746752,
            "b": 4,
            "nb": 3,
            "wrpos": 12364,
            "full": true
        },
...continued
```
...and the portion with the files. Line 174-205 (end):

```json
    ],
    "files": [
        {
            "typeflag": 53,
            "name": "usr/",
            "mode": 493,
            "mod": "2023-12-18T00:00:00Z",
            "offset": 512
        },
        {
            "typeflag": 53,
            "name": "usr/local/",
            "mode": 493,
            "mod": "2023-12-18T00:00:00Z",
            "offset": 1024
        },
        {
            "typeflag": 53,
            "name": "usr/local/bin/",
            "mode": 493,
            "mod": "2024-01-02T07:20:32Z",
            "offset": 1536
        },
        {
            "typeflag": 48,
            "name": "usr/local/bin/bun",
            "size": 92973248,
            "mode": 493,
            "mod": "2024-01-02T02:29:54Z",
            "offset": 2048
        }
    ]
}
```
