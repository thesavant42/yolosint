package explore

import (
	"text/template"

	v1 "github.com/jonjohnsonjr/dagdotdev/pkg/forks/github.com/google/go-containerregistry/pkg/v1"
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
<h1><a class="top" href="/"><img class="crane" src="/favicon.svg"/> <span class="link"></span></a></h1>
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
Search Docker Hub:
<p>	
</p>
<form action="https://hub.docker.com/search" method="GET" autocomplete="off" spellcheck="false">
<input size="100" type="text" name="q" placeholder="Search Docker Hub..."/>
<input type="submit" value="Search"/>
</form>
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
  <li><a href="/?image=ghcr.io/stargz-containers/node:13.13.0-esgz">ghcr.io/stargz-containers/node:13.13.0-esgz</a></li>
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
</style>
</head>
<h1><a class="top" href="/"><img class="crane" src="/favicon.svg"/> <span class="link"></span></a></h1>
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

td {
	vertical-align: top;
}

</style>
</head>
`

	bodyTemplate = `
<body>
<div>
<h1><a class="top" href="/"><img class="crane" src="/favicon.svg"/> <span class="link"></span></a>{{if .SaveURL}}<a href="{{.SaveURL}}"><img class="crane" src="/save-32.jpg" alt="save" title="Download file"/></a>{{end}}</h1>
{{ if .Up }}
<h2>{{ if and (ne .Up.Parent "docker.io") (ne .Up.Parent "index.docker.io") }}<a class="mt" href="/?repo={{.Up.Parent}}">{{.Up.Parent}}</a>{{else}}{{.Up.Parent}}{{end}}{{.Up.Separator}}{{if .RefHandler }}<a class="mt" href="/{{.RefHandler}}{{.Reference}}{{if .EscapedMediaType}}{{.QuerySep}}mt={{.EscapedMediaType}}{{end}}">{{.Up.Child}}</a>{{else}}{{.Up.Child}}{{end}}{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</h2>
{{ else }}
	<h2>{{.Reference}}{{ range .CosignTags }} (<a href="/?image={{$.Repo}}:{{.Tag}}">{{.Short}}</a>){{end}}{{if .Referrers}} <a href="/?referrers={{$.Repo}}@{{$.Descriptor.Digest}}">(referrers)</a>{{end}}</h2>
{{ end }}
{{ if .Descriptor }}
<table>
<tr><td>mediaType</td><td>{{.Descriptor.MediaType}}</td></tr>
<tr><td>digest</td><td><a class="mt" href="/{{.Handler}}{{$.Repo}}@{{.Descriptor.Digest}}{{.QuerySep}}{{if .EscapedMediaType}}mt={{.EscapedMediaType}}&{{end}}size={{.Descriptor.Size}}">{{.Descriptor.Digest}}</a></td></tr>
<tr><td>size</td><td>{{if .SizeLink}}<a class="mt" href="{{.SizeLink}}">{{.Descriptor.Size}}</a>{{else}}{{.Descriptor.Size}}{{end}}</td></tr>
{{if $.Subject}}<tr><td>OCI-Subject</td><td><a class="mt" href="/?image={{$.Repo}}@{{.Subject}}">{{.Subject}}</a></td></tr>{{end}}
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
	Repo             string
	CosignTags       []CosignTag
	Reference        string
	Up               *RepoParent
	Descriptor       *v1.Descriptor
	RefHandler       string
	Handler          string
	EscapedMediaType string
	QuerySep         string
	MediaTypeLink    string
	SizeLink         string
	HumanSize        string
	Referrers        bool
	Subject          string
	SaveURL          string
	Filename         string
}
