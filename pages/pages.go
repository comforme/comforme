package pages

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-zoo/bone"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var pageTemplate *template.Template

func init() {
	pageTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(pageTemplate.New("nav").Parse(templates.NavBar))
	template.Must(pageTemplate.New("fullPage").Parse(templates.SearchBar))
	template.Must(pageTemplate.New("content").Parse(pageTemplateText))
}

func PageHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	page := common.Page{}
	//posts := []common.Post{}

	sessionid, err := common.GetSessionId(res, req)
	if err != nil {
		log.Println("Failed to retrieve session id", err)
		return
	}
	data["formAction"] = req.URL.Path
	category := bone.GetValue(req, "category")
	slug := bone.GetValue(req, "pagename")
	page, _, err = databaseActions.GetPage(sessionid, category, slug)
	if err != nil {
		log.Println("Failed to retrieve page", err)
		return
	}
	data["Category"] = page.Category
	data["Slug"] = page.Slug
	data["Title"] = page.Title
	data["Description"] = page.Description

	common.ExecTemplate(pageTemplate, res, data)
}

const pageTemplateText = `
	<div class="content">
		<div class="row">
			<div class="column">
				<h1><a href="/page/{{.Category}}/{{.Slug}}">{{.Title}}</a></h1>
				<p>{{.Description}}</p>
			</div>
		</div>
	</div>
`
