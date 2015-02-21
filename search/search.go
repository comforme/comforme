package search

import (
	"html/template"
	"log"
	"net/http"

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

func SearchHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	if req.Method == "POST" {
		cookie, err := req.Cookie("sessionid")
		if err != nil {
			log.Println("Failed to retrieve sessionid in SearchHandler:", err)
			common.Logout(res, req)
			return
		}
		sessionid := cookie.Value
		query := req.PostFormValue("page-search")
		data["query"] = query
		data["pageTitle"] = req.PostFormValue("page-search")
		data["results"], err = databaseActions.SearchPages(sessionid, query)
		if err != nil {
			log.Println("Failed to retrieve search results for "+
				query, err)
		}
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
                {{if .query}}
                    <div class="alert-box secondary">Results for <span style="color:red">{{.query}}</span></div>
                {{end}}
            </div>
        </div>

        <div class="row">
            {{ range .results}}
            <div class="columns">
                <h2><a href="/page/{{.Category}}/{{.Slug}}">{{.Title}}</a></h2>
                <div>
                    <p>{{.Description}}</p>
                </div>
            </div>
            {{ end }}
        </div>
    </div>
</div>
`
