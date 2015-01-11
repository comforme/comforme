package profile

import (
	"html/template"
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var profileTemplate *template.Template

func init() {
	profileTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(profileTemplate.New("content").Parse(profileTemplateText))
}

func ProfileHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		// TODO uncomment when put to use
		//username := req.PostFormValue("username")
		//password := req.PostFormValue("password")
		//newPassword := req.PostFormValue("newPassword")
		//newPasswordConfirmation := req.PostFormValue("newPasswordConfirmation")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(profileTemplate, res, data)
}

const profileTemplateText = `<p>Profile Page</p>`
