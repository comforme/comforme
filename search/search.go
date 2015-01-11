package search

import (
	"html/template"
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var searchTemplate *template.Template

func init() {
	searchTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(searchTemplate.New("nav").Parse(templates.NavBar))
	template.Must(searchTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(searchTemplate.New("content").Parse(searchTemplateText))
}

func SearchHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		// TODO uncomment when put to use
		//query := req.PostFormValue("query")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(searchTemplate, res, data)
}

const searchTemplateText = `
<div class="content">
    {{template "searchBar" .}}
</div>`
