package main

import (
	"html/template"
)

const rootTemplateHTML = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Sign In with Apple Test Server</title>
</head>
<body>
<div><a href="{{.}}" target="_blank">Sign In with Apple</a></div>
<div><code>{{.}}</code></div>
<hr>
<div>
    <form method="POST" target="_blank" action="/validate">
        <label for="token">Validate refresh token:</label>
        <input type="text" id="token" name="token">
        <input type="submit" value="Submit">
    </form>
</div>
</body>
</html>
`

const resultTemplateHTML = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Sign In with Apple Test Server</title>
</head>
<body>
<h1>Result</h1>
<pre>{{.}}</pre>
</body>
</html>
`

var (
	rootTemplate   = template.Must(template.New("root").Parse(rootTemplateHTML))
	resultTemplate = template.Must(template.New("result").Parse(resultTemplateHTML))
)
