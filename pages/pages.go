package pages

import (
	"html/template"
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var pagesTemplate *template.Template

func init() {
	pagesTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(pagesTemplate.New("content").Parse(pagesTemplateText))
}

func PagesHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		// TODO uncomment when put to use
		//title := req.PostFormValue("title")
		//description := req.PostFormValue("description")
		//address := req.PostFormValue("address")
		//categories := req.PostFormValue("categories")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(pagesTemplate, res, data)
}

const pagesTemplateText = `<p>Pages!</p>`
