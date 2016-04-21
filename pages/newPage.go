package pages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"


	"github.com/comforme/comforme/algoliaUtil"
	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var newPageTemplate *template.Template

func init() {
	newPageTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(newPageTemplate.New("nav").Parse(templates.NavBar))
	template.Must(newPageTemplate.New("content").Parse(newPageTemplateText))
	template.Must(newPageTemplate.New("dropdown").Parse(templates.Dropdown))
}

func NewPageHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params, userInfo common.UserInfo) {
	data := map[string]interface{}{}
	data["siteName"] = common.SiteName
	data["formAction"] = req.URL.Path
	title := req.PostFormValue("title")
	description := req.PostFormValue("description")
	address := req.PostFormValue("address")
	website := req.PostFormValue("website")

	data["title"] = title
	data["description"] = description
	data["address"] = address
	data["website"] = website

	data["categoryDropdown"] = map[string]interface{}{}
	data["categoryDropdown"].(map[string]interface{})["name"] = "category"
	options, err := databaseActions.ListCategories()
	if err != nil {
		data["errorMsg"] = err.Error()
		goto render
	}
	data["categoryDropdown"].(map[string]interface{})["options"] = options
	data["categoryDropdown"].(map[string]interface{})["selected"] = req.PostFormValue("category")

	if req.Method == "POST" {

		if len(title) <= 1 {
			data["errorMsg"] = "Title must be more than 1 character long."
			goto render
		}
		if len(description) < common.MinDescriptionLength {
			data["errorMsg"] = fmt.Sprintf("Description must be at least %d characters long.", common.MinDescriptionLength)
			goto render
		}

		category, err := strconv.ParseInt(req.PostFormValue("category"), 0, 0)
		if err != nil || category < 0 {
			log.Println("Invalid category:", req.PostFormValue("category"))
			data["errorMsg"] = "Invalid category."
			goto render
		}

		categorySlug, pageSlug, err := databaseActions.CreatePage(userInfo.UserID, title, description, address, website, int(category))
		if err == nil {
			log.Printf("Created %s!\n", title)

			// Update algolia index
			page, err := databaseActions.GetPage(categorySlug, pageSlug)
			err = algoliaUtil.ExportPageRecord(page)
		  if err != nil {
				return
			}

			http.Redirect(res, req, "/page/"+categorySlug+"/"+pageSlug, http.StatusFound)
			return
		} else {
			data["errorMsg"] = err.Error()
		}
	}

render:
	common.ExecTemplate(newPageTemplate, res, data)
}

const newPageTemplateText = `
<div class="row">
	<div class="large-centered medium-centered large-8 medium-8 columns">
		<div class="content" id="add-page-form">{{if .successMsg}}
			<div class="alert-box success">{{.successMsg}}</div>{{end}}{{if .errorMsg}}
			<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
			<form method="POST" action="{{.formAction}}" align="center">
				<fieldset>
					<legend>Create a Resource New Page</legend>
					<div>
						<input type="text" name="title" placeholder="Resource page title"{{if .title}} value="{{ .title }}"{{end}} align="center" />
					</div>
					<div>
						<textarea name="description" placeholder="Unbiased description of resource" rows="15">{{if .description}}{{ .description }}{{end}}</textarea>
					</div>
					<div>
						<input type="text" name="address" placeholder="Physical address of resource (if applicable)"{{if .address}} value="{{ .address }}"{{end}} />
					</div>
					<div>
						<input type="text" name="website" placeholder="Resource's website (if applicable)"{{if .website}} value="{{ .website }}"{{end}} />
					</div>
					<div>
						{{template "dropdown" .categoryDropdown}}
					</div>
					<div style="text-align:center">
						<button type="submit" class="button" name="sign-up" value="true">Submit</button>
					</div>
				</fieldset>
			</form>
		</div>
	</div>
</div>
`
