package explore

import (
	"strings"
	"text/template"

	v1 "github.com/thesavant42/yolosint/pkg/forks/github.com/google/go-containerregistry/pkg/v1"
)

var (
	headerTmpl *template.Template
	bodyTmpl   *template.Template
	oauthTmpl  *template.Template
)

func init() {
	headerTmpl = template.Must(template.New("headerTemplate").Parse(headerTemplate))
	bodyTmpl = template.Must(template.New("bodyTemplate").Parse(bodyTemplate))
	oauthTmpl = template.Must(template.New("oauthTemplate").Parse(oauthTemplate))
}

const (
	landingPage = `
<html>
<body>
<head>
<title></title>
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

details > summary {
	list-style: none;
	cursor: pointer;
}
details > summary::before {
	content: "> ";
}
details[open] > summary::before {
	content: "v ";
}
details > summary::-webkit-details-marker {
	display: none;
}

</style>
</head>
<h1><a class="top" href="/"><img class="crane" src="/docdork-32.png"/> <span class="link"></span></a></h1>
<p>
<a href="/yolosint.user.js">yolosint</a> - by <a href="https://github.com/thesavant42/yolosint">@thesavant42</a>
</p>
<p>
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
<p>
<details>
<summary>Interesting examples</summary>
<ul>
  <li><a href="/?image=gcr.io/distroless/static">gcr.io/distroless/static:latest</a></li>
  <li><a href="/?repo=ghcr.io/homebrew/core/crane">ghcr.io/homebrew/core/crane</a></li>
  <li><a href="/?repo=registry.k8s.io">registry.k8s.io</a></li>
  <li><a href="/?image=registry.k8s.io/bom/bom:sha256-499bdf4cc0498bbfb2395f8bbaf3b7e9e407cca605aecc46b2ef1b390a0bc4c4.sig">registry.k8s.io/bom/bom:sha256-499bdf4cc0498bbfb2395f8bbaf3b7e9e407cca605aecc46b2ef1b390a0bc4c4.sig</a></li>
  <li><a href="/?image=docker/dockerfile:1.5.1">docker/dockerfile:1.5.1</a></li>
  <li><a href="/?image=tianon/true:oci">tianon/true:oci</a></li>
  <li><a href="/?image=kentico/ems:web-12.0.29">?image=kentico/ems:web-12.0.29</a></li>
</ul>
</details>
<p>
</body>
</html>
`

	oauthTemplate = `
<html>
<body>
<head>
<title></title>
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
</style>
</head>
<h1><a class="top" href="/"><img class="crane" src="/docdork-32.png"/> <span class="link"></span></a></h1>
<p>
It looks like we encountered an auth error:
</p>
<code>
{{.Error}}
</code>
<p>
If you trust <a class="mt" href="https://github.com/jonjohnsonjr">me</a>, click <a href="{{.Redirect}}">here</a> for oauth to use your own credentials.
</p>
</body>
</html>
`

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
	overflow-wrap: break-word;
	padding: 12px;
}

pre {
	white-space: pre-wrap;
}

.indent {
	margin-left: 2em;
}

td {
	vertical-align: top;
}

td:first-child {
	white-space: nowrap;
	width: 1px;
	overflow: visible;
}

</style>
</head>
`

	bodyTemplate = `
<body>
<div>
<h1><a class="top" href="/"><img class="crane" src="/docdork-32.png"/> <span class="link"></span></a>{{if .SaveURL}}<a href="{{.SaveURL}}"><img class="crane" src="/save-32.jpg" alt="save" title="Download file"/></a>{{end}}</h1>
{{ if .Up }}
<p><strong>{{ if and (ne .Up.Parent "docker.io") (ne .Up.Parent "index.docker.io") }}<a class="mt" href="/?repo={{.Up.Parent}}">{{.Up.Parent}}</a>{{else}}{{.Up.Parent}}{{end}}{{.Up.Separator}}{{if .RefHandler }}<a class="mt" href="/{{.RefHandler}}{{.Reference}}{{if .EscapedMediaType}}{{.QuerySep}}mt={{.EscapedMediaType}}{{end}}">{{.Up.Child}}</a>{{else}}{{.Up.Child}}{{end}}</strong>{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</p>
{{if .DockerPull}}<p><span style="background:#2a2a3e;padding:4px 8px;border-radius:4px;border:1px solid #444;">docker pull {{.DockerPull}} <button onclick="navigator.clipboard.writeText('docker pull {{.DockerPull}}')" style="background:#444;color:inherit;border:1px solid #666;padding:2px 6px;border-radius:4px;cursor:pointer;margin-left:4px;font:inherit;">Copy</button></span></p>{{end}}
{{ else }}
	<p><strong>{{.Reference}}</strong>{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</p>
{{if .DockerPull}}<p><span style="background:#2a2a3e;padding:4px 8px;border-radius:4px;border:1px solid #444;">docker pull {{.DockerPull}} <button onclick="navigator.clipboard.writeText('docker pull {{.DockerPull}}')" style="background:#444;color:inherit;border:1px solid #666;padding:2px 6px;border-radius:4px;cursor:pointer;margin-left:4px;font:inherit;">Copy</button></span></p>{{end}}
{{ end }}
{{ if .Descriptor }}
<table>
<tr><td><img src="/ant-design--container-outlined.png" alt="manifest" style="height:16px;vertical-align:middle"/> {{.AbbreviatedMediaType}}</td></tr>
{{if $.Subject}}<tr><td>OCI-Subject</td><td></td><td><a class="mt" href="/?image={{$.Repo}}@{{.Subject}}">{{.Subject}}</a></td></tr>{{end}}
</table>
{{end}}
{{if .Filename}}<h3>{{.Filename}}</h3>{{end}}
</div>
`

	footer = `
</body>
</html>
`
)

type RepoParent struct {
	Parent    string
	Child     string
	Separator string
}

type OauthData struct {
	Error    string
	Redirect string
}

type TitleData struct {
	Title string
}
type CosignTag struct {
	Tag   string
	Short string
}

type HeaderData struct {
	Repo                 string
	CosignTags           []CosignTag
	Reference            string
	Up                   *RepoParent
	Descriptor           *v1.Descriptor
	RefHandler           string
	Handler              string
	EscapedMediaType     string
	QuerySep             string
	MediaTypeLink        string
	SizeLink             string
	HumanSize            string
	Referrers            bool
	Subject              string
	SaveURL              string
	Filename             string
	DockerPull           string
	AbbreviatedMediaType string
}

// AbbreviateMediaType converts a full media type string to a short label
func AbbreviateMediaType(mt string) string {
	switch {
	case strings.Contains(mt, "manifest.list"):
		return "list v2"
	case strings.Contains(mt, "image.index"):
		return "index"
	case strings.Contains(mt, "docker.distribution.manifest.v2"):
		return "v2"
	case strings.Contains(mt, "docker.distribution.manifest.v1"):
		return "v1"
	case strings.Contains(mt, "oci.image.manifest"):
		return "OCI"
	default:
		return mt
	}
}
