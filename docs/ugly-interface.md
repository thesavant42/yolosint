User: I understand the fake json is intentional, but it's ugly as hell

Answer
I get the sentiment—the current HTML-styled JSON view is functional but visually stark. The styling comes from the CSS in `templates.go` and the rendering logic in `render.go`’s `jsonOutputter` [1](#2-0) [2](#2-1) . You could improve aesthetics in a few ways without breaking navigation:

- Enhance the CSS in `headerTemplate`/`bodyTemplate` for better typography, spacing, and syntax highlighting [1](#2-0) .
- Add a new `render=table` mode in `renderContent` that outputs HTML tables for top-level arrays (e.g., `layers`, `annotations`), similar to how `render=history` already works [3](#2-2) [4](#2-3) .
- Extend `renderMap` to detect specific keys and render them as tables while keeping the rest as styled JSON [5](#2-4) .

## Details
- Styling hooks: The page uses minimal CSS in `templates.go` with monospace font and basic link styles [6](#2-5) . Adding more sophisticated syntax highlighting or spacing here would directly affect the “JSON-like” view.
- Rendering pipeline: `jsonOutputter` emits HTML with linkified digests and annotations [7](#2-6) . You could add a parallel outputter for tables without disrupting this.
- Precedent for tables: `render=history` already produces clean HTML tables with digest links and sizes [4](#2-3) . Reusing that pattern for other sections would be straightforward.

## Notes
- Any change should preserve the clickable links (digests → /fs/, annotations → specialized renderers), as they’re core to navigation [8](#2-7) .
- A full visual overhaul would likely require new CSS classes and possibly a new template, but the current structure supports incremental improvements [9](#2-8) .


### Citations

**File:** internal/explore/templates.go (L203-291)
```go
	headerTemplate = `
<html>
<head>
<title>{{.Title}}</title>
<link rel="icon" href="/favicon.svg">
<style>
.mt:hover {
	text-decoration: underline;
}

.mt {
	color: inherit;
	text-decoration: inherit;
}

.link {
	position: relative;
	bottom: .125em;
}

.crane {
	height: 1em;
	width: 1em;
}

.top {
	color: inherit;
	text-decoration: inherit;
}

:root {
  color-scheme: light dark;
}

body {
	font-family: monospace;
	width: fit-content;
	overflow-wrap: anywhere;
	padding: 12px;
}

pre {
	white-space: pre-wrap;
}

.indent {
	margin-left: 2em;
}

.noselect {
	user-select: none;
	-webkit-user-select: none;
	width: fit-content;
	overflow-wrap: none;
	padding-right: 1em;
	text-align: right;
	white-space: nowrap;
}

td {
	vertical-align: top;
}

// Adapted from https://medium.com/allenhwkim/how-to-build-tabs-only-with-css-844718d7de2f
input + label { display: inline-block } /* show labels in line */
input { display: none; }                /* hide radio buttons */
input ~ .tab { display: none }          /* hide contents */

/* show contents only for selected tab */
#tab1:checked ~ .tab.content1,
#tab2:checked ~ .tab.content2 { display: block; }

input + label {             /* box with rounded corner */
	display: inline-block;
  border: 1px solid #999;
  padding: 4px 12px;
  border-radius: 4px 4px 0 0;
  position: relative;
  top: 1px;
  opacity: 50%;
}
input:checked + label {     /* white background for selected tab */
  opacity: 100%;
}
input ~ .tab {          /* grey line between tab and contents */
  border-top: 1px solid #999;
  padding-top: 0.5em;
}
</style>
```

**File:** internal/explore/templates.go (L295-329)
```go
	bodyTemplate = `
<body>
<div>
<h1><a class="top" href="/"><img class="crane" src="/favicon.svg"/> <span class="link">Registry Explorer</span></a></h1>
{{ if .Up }}
<h2>{{ if and (ne .Up.Parent "docker.io") (ne .Up.Parent "index.docker.io") }}<a class="mt" href="/?repo={{.Up.Parent}}">{{.Up.Parent}}</a>{{else}}{{.Up.Parent}}{{end}}{{.Up.Separator}}{{if .RefHandler }}<a class="mt" href="/{{.RefHandler}}{{.Reference}}{{if .EscapedMediaType}}{{.QuerySep}}mt={{.EscapedMediaType}}{{end}}">{{.Up.Child}}</a>{{else}}{{.Up.Child}}{{end}}{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</h2>
{{ else }}
	<h2>{{.Reference}}{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</h2>
{{ end }}
{{ if .Descriptor }}
<input type="radio" name="tabs" id="tab1" checked />
<label for="tab1">HTTP</label>
<input type="radio" name="tabs" id="tab2" />
<label for="tab2">OCI</label>
<div class="tab content1">
Content-Type: {{if .MediaTypeLink}}<a class="mt" href="{{.MediaTypeLink}}">{{.Descriptor.MediaType}}</a>{{else}}{{.Descriptor.MediaType}}{{end}}<br>
Docker-Content-Digest: <a class="mt" href="/{{.Handler}}{{$.Repo}}@{{.Descriptor.Digest}}{{if .EscapedMediaType}}{{.QuerySep}}mt={{.EscapedMediaType}}{{end}}&size={{.Descriptor.Size}}">{{.Descriptor.Digest}}</a><br>
<span{{if .HumanSize}} title="{{.HumanSize}}"{{end}}>Content-Length: {{if .SizeLink}}<a class="mt" href="{{.SizeLink}}">{{.Descriptor.Size}}</a>{{else}}{{.Descriptor.Size}}{{end}}</span><br>
{{if $.Subject}}OCI-Subject: <a class="mt" href="/?image={{$.Repo}}@{{.Subject}}">{{.Subject}}</a><br>{{end}}
</div>
<div class="tab content2">
{<br>
&nbsp;&nbsp;"mediaType": "{{.Descriptor.MediaType}}",<br>
&nbsp;&nbsp;"digest": "{{.Descriptor.Digest}}",<br>
&nbsp;&nbsp;"size": {{.Descriptor.Size}}<br>
}<br>
</div>

{{end}}
</div>
{{ if .JQ }}
<h4><span style="padding:0;" class="noselect">$ </span>{{.JQ}}</h4>

{{ end }}
`
```

**File:** internal/explore/render.go (L30-44)
```go
type jsonOutputter struct {
	w    io.Writer
	u    *url.URL
	name string
	repo string
	mt   string

	fresh []bool
	jq    []string
	key   bool
	root  map[string]interface{}
	isMap bool

	dockerHub bool
}
```

**File:** internal/explore/render.go (L80-103)
```go
func (w *jsonOutputter) Linkify(mt string, h v1.Hash, size int64) {
	w.tabf()
	qs := "/?"
	handler := handlerForMT(mt)
	if strings.Contains(handler, "?") {
		qs = "&"
	}
	if strings.HasSuffix(mt, "+wasm") {
		w.Printf(`"<a href="https://github.com/opencontainers/artifacts/issues/61">%s</a>"`, html.EscapeString(h.String()))
	} else if h.String() == emptyDigest {
		w.Printf(`"<a href="/%s%s@%s%smt=%s" title="this is an empty layer that only modifies metadata, so it has no filesystem content">%s</a>"`, handler, w.repo, h.String(), qs, url.QueryEscape(mt), html.EscapeString(h.String()))
	} else if size != 0 {
		image := w.u.Query().Get("image")
		if w.jth(-1) == ".config" && image != "" {
			w.Printf(`"<a href="/%s%s@%s%smt=%s&size=%d&manifest=%s">%s</a>"`, handler, w.repo, h.String(), qs, url.QueryEscape(mt), size, image, html.EscapeString(h.String()))
		} else {
			w.Printf(`"<a href="/%s%s@%s%smt=%s&size=%d">%s</a>"`, handler, w.repo, h.String(), qs, url.QueryEscape(mt), size, html.EscapeString(h.String()))
		}
	} else {
		w.Printf(`"<a href="/%s%s@%s%smt=%s">%s</a>"`, handler, w.repo, h.String(), qs, url.QueryEscape(mt), html.EscapeString(h.String()))
	}
	w.unfresh()
	w.key = false
}
```

**File:** internal/explore/render.go (L479-483)
```go
		if k == "layers" && len(w.jq) == 0 {
			image := w.u.Query().Get("image")
			w.Layers(image, "layers")
		} else if k == "history" && shouldHistory(w.mt) {
			w.History(k)
```

**File:** internal/explore/render.go (L657-713)
```go
		case "mediaType":
			mt := ""
			if err := json.Unmarshal(v, &mt); err != nil {
				log.Printf("Unmarshal mediaType %q: %v", string(v), err)
			} else {
				w.Doc(getLink(mt), strconv.Quote(mt))

				// Don't fall through to renderRaw.
				continue
			}
		case "urls":
			if digest, ok := rawMap["digest"]; ok {
				h := v1.Hash{}
				if err := json.Unmarshal(digest, &h); err != nil {
					log.Printf("Unmarshal digest %q: %v", string(digest), err)
				} else {
					// We got a digest, so we can link to some blob.
					if urls, ok := o["urls"]; ok {
						if ii, ok := urls.([]interface{}); ok {
							if len(ii) == 0 {
								w.Value([]byte("[]"))
								continue
							}
							w.StartArray()
							for _, iface := range ii {
								if original, ok := iface.(string); ok {
									scheme := "https"
									u := original
									if strings.HasPrefix(original, "https://") {
										u = strings.TrimPrefix(original, "https://")
									} else if strings.HasPrefix(original, "http://") {
										u = strings.TrimPrefix(original, "http://")
										scheme = "http"
									}
									w.URL("/"+scheme+"/"+url.PathEscape(u)+"@"+h.String(), original)
								} else {
									// This wasn't a list of strings, render whatever we found.
									b, err := json.Marshal(iface)
									if err != nil {
										return err
									}
									raw := json.RawMessage(b)
									if err := renderRaw(w, &raw); err != nil {
										return err
									}
								}
							}
							w.EndArray()

							// Don't fall through to renderRaw.
							continue
						}
					}
				}
			}
		case "Docker-reference", "docker-reference":
			if js, ok := o[k]; ok {
```

**File:** internal/explore/explore.go (L961-1004)
```go
func (h *handler) renderContent(w http.ResponseWriter, r *http.Request, ref name.Reference, b []byte, output *jsonOutputter, u url.URL) error {
	switch r.URL.Query().Get("render") {
	case "raw":
		fmt.Fprintf(w, "<pre>")
		if _, err := w.Write(b); err != nil {
			return err
		}
		fmt.Fprintf(w, "</pre>")
	case "x509":
		return renderx509(w, b)
	case "cert":
		return renderCert(w, b, u)
	case "der":
		return renderDer(w, b)
	case "history":
		if types.MediaType(r.URL.Query().Get("mt")).IsSchema1() {
			return renderDockerfileSchema1(w, b, ref.Context())
		} else {
			return h.renderDockerfile(w, r, ref, b)
		}
	case "created_by":
		fmt.Fprintf(w, "<pre>")
		if err := renderCreatedBy(w, b); err != nil {
			return err
		}
		fmt.Fprintf(w, "</pre>")
	case "xxd":
		return renderOctets(w, r, b)
	case "timestamp":
		ts, err := timestamp.Parse(b)
		if err != nil {
			return err
		}
		j, err := json.Marshal(ts)
		if err != nil {
			return err
		}
		return renderJSON(output, j)
	default:
		return renderJSON(output, b)
	}

	return nil
}
```
