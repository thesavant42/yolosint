# Image View UI Fix

## Problem
Icons on separate lines due to table wrappers.

## Solution
Remove the table wrappers. Output icons directly.

---

## Changes

### templates.go

**Current:**
```html
<table>
<tr><td><img .../> {{.AbbreviatedMediaType}}</td></tr>
</table>
```

**New:**
```html
<img src="/ant-design--container-outlined.png" alt="manifest" style="height:16px;vertical-align:middle"/> {{.AbbreviatedMediaType}}
```

### render.go  

**Current:**
```go
w.Print(`<table><tr><td>`)
w.Printf(`<a href="..."><img .../></a>`)
w.Print(`</td></tr></table>`)
```

**New:**
```go
w.Printf(`<a href="..." title="Config: %s"><img src="/eos-icons--init-container-outlined.png" alt="config" style="height:16px;vertical-align:middle"/></a>`, ...)
```

---

## Result
```
[manifest-icon] OCI [config-icon]
```
