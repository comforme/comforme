package home

import (
	"fmt"
	"github.com/comforme/comforme/common"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(homeTemplate.New("nav").Parse(templates.NavBar))
	template.Must(homeTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(homeTemplate.New("content").Parse(homeTemplateText))
}

func HomeHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params, userInfo common.UserInfo) {
	data := map[string]interface{}{}
	pages, err := databaseActions.GetTopPages()
	if err != nil {
		log.Println("Failed to retrieve top results:", err)
	} else {
		data["Top"] = fmt.Sprintf("%+v", pages)
	}

	common.ExecTemplate(homeTemplate, res, data)
}

const homeTemplateText = `
<div class="content">
	<div class="row">
		<div class="columns">
			<h1>Search</h1>
		</div>
	</div>
	<div class="row">
		<div class="columns">
			{{template "searchBar" .}}
			<p>{{.Top}}</p>
		</div>
	</div>
</div>`
