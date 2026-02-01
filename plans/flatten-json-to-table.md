# Flatten JSON to Table

## render.go changes

| Line | Current | Replace With |
|------|---------|--------------|
| 159 | `"{"` | `"<table>"` |
| 171 | `"}"` | `"</table>"` |
| 179 | `"["` | `"<table>"` |
| 191 | `"]"` | `"</table>"` |
| 146 | `'"%s":'` | `"<tr><td>%s</td><td>"` |
| 207 | `","` | `"</td></tr>"` |
| 211 | `w.div()` | `(remove)` |
| 226 | `<div class="indent">` | `(remove)` |

Remove quotes from string output in link methods (Linkify, BlueDoc, URL, Blob, LinkImage, LinkRepo).

## templates.go changes

1. `headerTemplate` CSS (line 215) already has `td { vertical-align: top; }` - no change needed

2. `bodyTemplate` (lines 233-237) - convert descriptor display from fake JSON to table:

**Current:**
```
{<br>
&nbsp;&nbsp;"mediaType": "{{.Descriptor.MediaType}}",<br>
&nbsp;&nbsp;"digest": "<a ...>{{.Descriptor.Digest}}</a>",<br>
&nbsp;&nbsp;"size": ...
}
```

**Replace with:**
```
<table>
<tr><td>mediaType</td><td>{{.Descriptor.MediaType}}</td></tr>
<tr><td>digest</td><td><a ...>{{.Descriptor.Digest}}</a></td></tr>
<tr><td>size</td><td>...</td></tr>
</table>
```
