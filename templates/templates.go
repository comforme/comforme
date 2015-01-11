package templates

//Generic page template

const rootPageTemplateHtml = `
<!DOCTYPE html>
<html>
  <head>
	<title>{{.PageTitle}}</title>
	<style type="text/css" media="all">
		<!--
			.your_entry { background: white; }
			.other_entry { background: #00FFFF; }
			.errortext { color: red; font-weight: bold }
			body { font-family: Helvetica, Arial, Sans-Serif; }
		-->
	</style>
  </head>
  <body>
	Hello World!
  </body>
</html>
`
