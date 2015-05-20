package search

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var searchTemplate *template.Template

func init() {
	searchTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(searchTemplate.New("nav").Parse(templates.NavBar))
	template.Must(searchTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(searchTemplate.New("content").Parse(searchTemplateText))
}

func SearchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params, userInfo common.UserInfo) {
	data := map[string]interface{}{}
	if common.CheckParam(req.URL.Query(), "q") {
		query := req.URL.Query()["q"][0]
		log.Println("Performing search for:", query)
		data["query"] = query
		data["pageTitle"] = query
		var err error
		data["results"], err = databaseActions.SearchPages(userInfo.SessionID, query)
		if err != nil {
			log.Println("Failed to retrieve search results for "+
				query, err)
		} else {
			log.Printf("Search results for %s:\n%+v\n", query, data["results"])
		}
	} else {
		data["pageTitle"] = "Search"
	}

	common.ExecTemplate(searchTemplate, res, data)
}

// TODO add description limits and ellipses link to full page
const searchTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns">
				<h1>Search</h1>
				{{template "searchBar" .}}
				<div class="alert-box secondary">{{if .results}}
					Results for <span style="color:red">{{.query}}</span>{{else}}
					<span style="color:red">No matches found for "{{.query}}".</span> Would you like to <a href="/newPage">add a new resource</a>?{{end}}
				</div>
			</div>
		</div>
		<div class="row">{{range .results}}
			<div class="columns">
				<h2><a href="/page/{{.CategorySlug}}/{{.PageSlug}}">{{.Title}}</a></h2>
				<div>
					<p>{{.Description}}</p>
				</div>
			</div>{{ end }}
		</div>
	</div>
</div>
`
