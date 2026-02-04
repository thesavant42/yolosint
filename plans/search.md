# Task 1 Log to sqlite

**Current Architecture:**


If you ask to read any other files you will be immediately terminated.

- [`internal/soci/indexer.go`](internal/soci/indexer.go:22) - The `Indexer` struct parses tar files and creates a [`TOC`](internal/soci/toc.go:11) (Table of Contents) containing file metadata
- [`internal/explore/cache.go`](internal/explore/cache.go:23) - Stores the TOC as `toc.json.gz` and the indexed archive as `.tar.gz`
- The TOC already contains all the file metadata you want to search: `Name`, `Size`, `Offset`, `Linkname`, `ModTime`, etc.

**The Problem:**
You already have JSON data with all the file information, but it's locked up in gzipped tar files. You want to write this data to SQLite so you can actually search your history.

**The Solution:**
Add SQLite storage that captures the TOC data during indexing, before/while writing the tar files.

**Summary:**

The TOC data you need is already in [`soci.TOC`](internal/soci/toc.go:11) - filenames, sizes, timestamps, everything. The [`dirCache.Put()`](internal/explore/cache.go:178) method receives this data before writing to tar.

**The fix:** Add to [`internal/explore/cache.go`](internal/explore/cache.go:150)

**Index by layer digest**

--- 

## Requirements

Failure to comply with these requiremments will result in termination.

- must work with [example toc.json](/cache/sha256-2b97a650489f286782ed9553c3c7e16669fe46b3aaacc3cd4944bfde910b002a.0/toc.json)
- must work with [docker-compose.yml](/docker-compose.yml) AND [docker-compose-synology](/docker-compose.synology.yml)
- Your plan must be in markdown
- it must include line-number references to the file edits, no guesses allowed.
- it must specifically correlate to a task in this document. 
- It must not be pseudocode. It must contain actual goloang.
- It must not be pseudocode. It must contain actual goloang.