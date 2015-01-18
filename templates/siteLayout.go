package templates

const SiteLayout = `<!DOCTYPE html>
<html>
<head>
	<link href="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/css/normalize.min.css" rel="stylesheet" type="text/css" />
	<link href="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/css/foundation.min.css" rel="stylesheet" type="text/css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/vendor/jquery.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/foundation.min.js"></script>
	<meta charset="utf-8" />
	<title>ComFor.Me{{if .pageTitle}} - {{.pageTitle}}{{end}}</title>
    <link rel="stylesheet" href="/static/style_css" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/foundicons/3.0.0/foundation-icons.css" />
	<script scr="https://cdnjs.cloudflare.com/ajax/libs/foundation/5.5.0/js/login.js"></script>
	<script src='https://www.google.com/recaptcha/api.js'></script>
</head>
<body>
    {{ template "nav" .}}
    {{ template "content" .}}
	<script>$(document).foundation();</script>
</body>
</html>
`
