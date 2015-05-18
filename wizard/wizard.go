package wizard

import (
	"html/template"
	"log"
	"net/http"
	//"os"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var communitiesTemplate *template.Template

func init() {
	communitiesTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(communitiesTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(communitiesTemplate.New("communities").Parse(templates.Communities))
	template.Must(communitiesTemplate.New("content").Parse(communitiesTemplateText))
}

func WizardHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	cookie, err := req.Cookie("sessionid")
	if err != nil {
		log.Println("Failed to retrieve sessionid:", err)
		common.Logout(res, req)
		return
	}
	sessionid := cookie.Value

	data["communitiesCols"], err = databaseActions.GetCommunityColumns(sessionid)
	if err != nil {
		log.Println("Error listing communities:", err)
		common.Logout(res, req)
		return
	}

	common.ExecTemplate(communitiesTemplate, res, data)
}

const communitiesTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings Wizard</h1>
                {{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
                {{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>
{{ template "communities" .}}
					<form action="/">
						<button type="submit">Finish</button>
					</form>
				</section>
			</div>
		</div>
	</div>
	<script src="/static/js/settings_js"></script>
`
