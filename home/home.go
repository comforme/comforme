package home

import (
	"github.com/comforme/comforme/common"
	"html/template"
	"net/http"
	// "github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(homeTemplate.New("nav").Parse(templates.NavBar))
	template.Must(homeTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(homeTemplate.New("content").Parse(homeTemplateText))
	template.Must(homeTemplate.New("summary").Parse(summTemplateText))
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	common.ExecTemplate(homeTemplate, res, data)
}

const homeTemplateText = `
<div class="content">
    {{template "searchBar" .}}
    {{template "summary"   .}}
</div>`


const summTemplateText = `
<div class="content">
	<div class="row">
		<div class="large-centered columns">
                	<h2>What is Comfor.me?</h2>
                	<div>
                    	<p>Comfor.me (Community for Me) is a community-rated and identity-oriented social network/service listing. Users can find accepting communities and services based on a wide array of keywords. Users can also start their own communities categorized by aforementioned keywords. Comfor.me makes it easier for an individual to find communities and services which accept them for who they are.</a></p>
                	</div>
            	</div>
	</div>
</div>`
