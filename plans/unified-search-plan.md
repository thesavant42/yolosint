# Plan: Unified Search Box Implementation

## Summary

Replace 3 separate search forms in [`internal/explore/templates.go`](../internal/explore/templates.go:82) with a single form containing a dropdown selector.

## Current Code (lines 82-104)

```html
<p>
Enter a <strong>public</strong> image, e.g. <tt>"ubuntu:latest"</tt>:
</p>
<form action="/" method="GET" autocomplete="off" spellcheck="false">
<input size="100" type="text" name="image" value="ubuntu:latest"/>
<input type="submit" />
</form>
<p>
<p>
Enter a <strong>public</strong> repository, e.g. <tt>"ubuntu"</tt>:
</p>
<form action="/" method="GET" autocomplete="off" spellcheck="false">
<input size="100" type="text" name="repo" value="ubuntu"/>
<input type="submit" />
</form>
<p>
Search <a href="https://hub.docker.com">Docker Hub</a>
<p>	
</p>
<form action="https://hub.docker.com/search" method="GET" autocomplete="off" spellcheck="false">
<input size="100" type="text" name="q" placeholder="Search Docker Hub..."/>
<input type="submit" value="Search"/>
</form>
```

## Replacement Code

Layout: `[Text Input] [Dropdown] [Submit]` - dropdown next to button

```html
<p>
Search for a <strong>public</strong> container:
</p>
<form id="searchForm" action="/" method="GET" autocomplete="off" spellcheck="false">
<input id="searchInput" size="40" type="text" name="image" value="ubuntu:latest"/> <select id="searchType" onchange="updateSearch()">
  <option value="image" selected>Image</option>
  <option value="repo">Repository</option>
  <option value="dockerhub">Docker Hub</option>
</select> <input type="submit" />
</form>

<script>
function updateSearch() {
  var sel = document.getElementById('searchType');
  var inp = document.getElementById('searchInput');
  var frm = document.getElementById('searchForm');
  var v = sel.value;
  
  if (v === 'image') {
    inp.name = 'image';
    inp.placeholder = 'ubuntu:latest';
    inp.value = 'ubuntu:latest';
    frm.action = '/';
  } else if (v === 'repo') {
    inp.name = 'repo';
    inp.placeholder = 'ubuntu';
    inp.value = 'ubuntu';
    frm.action = '/';
  } else {
    inp.name = 'q';
    inp.placeholder = 'Search Docker Hub...';
    inp.value = '';
    frm.action = 'https://hub.docker.com/search';
  }
}
</script>
```

## Behavior Matrix

| Selection | Input Name | Form Action | Default Value |
|-----------|-----------|-------------|---------------|
| Image | `image` | `/` | `ubuntu:latest` |
| Repository | `repo` | `/` | `ubuntu` |
| Docker Hub | `q` | `https://hub.docker.com/search` | (empty) |

## Implementation Steps

1. Open [`internal/explore/templates.go`](../internal/explore/templates.go)
2. Locate the `landingPage` constant (line 22)
3. Replace lines 82-104 with the new unified form code above
4. Test all three search modes work correctly

## Backward Compatibility

- Direct GET requests to `/?image=...` continue to work (server unchanged)
- Direct GET requests to `/?repo=...` continue to work (server unchanged)
- No server-side changes required
