package home

import (
	"net/http"
        "html/template"
	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	// TODO: Add template and compile it.
	tmpl, _ := template.New("test").Parse(rootPageTemplateHtml)
	common.ExecTemplate(tmpl, res, data)
}

//Generic page template

const rootPageTemplateHtml = `
<!DOCTYPE html>
<html>
  <head>
	<title>{{.PageTitle}}</title>
	<style type="text/css" media="all">
		<!--
			.your_entry { background: yellow; }
			.other_entry { background: #00FFFF; }
			.errortext { color: red; font-weight: bold }
			body { font-family: Helvetica, Arial, Sans-Serif; color: red }
		-->
	</style>
  </head>
  <body>
	I really love Wine!
  </body>
</html>
`
