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
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	common.ExecTemplate(homeTemplate, res, data)
}

const homeTemplateText = `
<div class="content">
    {{template "searchBar" .}}

</div>`

const homeOptTemplateText = 
`<div class="content">
	<div class="row">
		<div class="column">
			<h1>Search for Communities</h1>
			<form method="post" action="/">
				<div class="row collapse">
					<div class="small-10 columns">
						<input type="text" placeholder="enter community name" name="page-search" id="page-search-textbox">
					</div>
				<div class="small-2 columns">
					<button type="submit" class="button postfix">Search</button>
				</div>
				</div>
			</form>
		</div>
	</div>
</div>`

