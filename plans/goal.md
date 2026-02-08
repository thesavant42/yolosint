# Refactor Template

## Summary of Task

Refactor: [rendergo](/internal/explore/render.go)

### Problem Statement

Current view is cluttered and needs alignment of the elements. The follow screencapture demonstrates the current state of the view when browsing `http://192.168.1.37:8042/?image=smoshysmosh/composer@sha256:4de9e962f0780c5bf92917340931f28f20556326ec6b102eaa24984622141858&mt=application%2Fvnd.docker.distribution.manifest.v2%2Bjson&size=3458`

![current](/plans/current.png)


---

## Mockup

This is a mockup of the same screenshot, but with the elements restructured the way that I would like them.

![mockup](/plans/mockup.png)

### List of changes

 1. **Truncated** the `sha256` digest **beneath the logo** to *12* characters
 2. Truncate Manifest Type (there aren't that many types, easy to abbreviate)
 3. **Moved** `combined layers view` and `referrers` to the same line, below the topmost digest.
 4. **Removed** the "`Docker pull`" box from this view (not the tag, this command will fail)
 5. `CONFIG` label is now capitalized and placed in alignment with the "`LIST VIEW:`" column label. 
 6. `CONFIG` size is aligned with file size
 7. `sha256` digest for config is aligned with other `sha256` digests
 8. **Removed** the redundant Layer IDX column. Now there are *4* columns: The `idx`, the `file size`, the `sha256` digest, and the download `buttons`.
 9. **Removed** redundant `digest` text label and sha256:digest, it's the same one we're truncating below the logo, so no need to have it twice.

## Task 

1 . Make a plan with the required lines of code to change.
- Include sufficient detail for handoff to a coding agent
    - line numnbers where applicable

### Validation Criteria
- Must match the mockup screenshsot 1:1
- validate against this doc, [goal.md](/plans/goal.md)

