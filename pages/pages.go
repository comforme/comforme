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
	template.Must(pagesTemplate.New("nav").Parse(templates.NavBar))
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

const pagesTemplateText = `
<div class="row">
<div class="content" id="add-page-form">
	<form method="POST" action="/">
		<div>
			<input type="text" name="title" placeholder="page title">
		</div>
		<div>
			<input type="text" name="description" placeholder="description">
		</div>
		<div>
			<input type="text" name="address/location" placeholder="address">
		</div>
		<div>
			<input type="text" name="categories" placeholder="categories">
		</div>
		<div>
			<button type="submit" class="button" name="sign-up" value="true">Submit</button>
		</div>
	</form>
</div>
</div>		
`
