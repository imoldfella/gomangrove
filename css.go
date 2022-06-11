package main

var pinList = `
<div class="container">
{{range . }}
<a class="card" href='{{.Link}}'>
  <div class="appicon"> <img alt="{{.Title}}" class='cat' src="theme/thumb/{{.Image}}" /></div>
  <div class="caption">{{.Title}}</div>
</a>
{{end}}
</div>

`
var lessonList = `
<div class="container">
{{range . }}
<a class="tile" href="{{.Link}}">
<div class="caption2">{{.Number}}</div>
</a>
{{end}}
</div>
</section>`

var page = `
<!DOCTYPE html>
<html lang="en-US">
<head>
<meta charset="UTF-8" />
<meta name="description" content="Free curriculum">
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Mangrove</title>
<link rel="stylesheet" href="/theme/style.css">
</head>
<body>
<nav class="navbar" role="navigation" aria-label="main navigation">
<div class="navbar-brand">
  <a class="navbar-item" href="/index.html">
	📚 mangrove
  </a>
</div>

</nav>
	{{.Content}}
</body>
</html>
`
