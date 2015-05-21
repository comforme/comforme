package home

import (
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
	topPages, err := databaseActions.GetTopPages()
	if err != nil {
		log.Println("Failed to retrieve top results:", err)
	} else {
		data["topPages"] = topPages
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
	<div class="row">
		<div class="columns left">
			Lost? <a href="/wizard">Take the tour</a> again.
		</div>
	</div>
	<div class="row">
		<div class="columns left">
			<h2>Top Resources:</h2>
		</div>{{range .topPages}}
		<div class="columns left large-3 medium-4 small-6 xsmall-12">
			<h3><a href="/page/{{.CategorySlug}}/{{.PageSlug}}">{{.Title}}</a></h3>
		</div>{{ end }}
	</div>
</div>`
