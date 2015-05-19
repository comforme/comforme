package wizard

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"os"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/requireLogin"
	"github.com/comforme/comforme/templates"
)

// Messages
const (
	invalidLink   = "Invalid link."
	accountExists = "You already have an account with this email address."
	linkUsed      = "This link has been previously used and is no longer valid."
)

var communitiesTemplate *template.Template
var messageTemplate *template.Template
var registerTemplate *template.Template

func init() {
	// Message page template
	messageTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(messageTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(messageTemplate.New("wizardContent").Parse(""))
	template.Must(messageTemplate.New("content").Parse(wizardTemplateText))

	// Community selection page template
	communitiesTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(communitiesTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(communitiesTemplate.New("communitiesContent").Parse(templates.Communities))
	template.Must(communitiesTemplate.New("wizardContent").Parse(communitiesTemplateText))
	template.Must(communitiesTemplate.New("content").Parse(wizardTemplateText))

	// Register page template
	registerTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(registerTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(registerTemplate.New("wizardContent").Parse(registerTemplateText))
	template.Must(registerTemplate.New("content").Parse(wizardTemplateText))
}

func WizardHandler(res http.ResponseWriter, req *http.Request) {

	_, err := req.Cookie("sessionid")
	if err == nil { // Found cookie
		requireLogin.RequireLogin(introWizardHandler)(res, req)
		return
	}

	data := map[string]interface{}{}
	data["pageTitle"] = "Wizard"

	// Check for duplicate parameters.
	for _, value := range req.URL.Query() {
		if len(value) != 1 {
			data["errorMsg"] = invalidLink
			common.ExecTemplate(messageTemplate, res, data)
			return
		}
	}

	// Check for action
	if !common.CheckParam(req.URL.Query(), "action") {
		data["errorMsg"] = invalidLink
		common.ExecTemplate(messageTemplate, res, data)
		return
	}
	actionName := req.URL.Query()["action"][0]
	log.Printf("Action: %s\n", actionName)

	if !common.CheckParam(req.URL.Query(), "email") ||
		!common.CheckParam(req.URL.Query(), "date") ||
		!common.CheckParam(req.URL.Query(), "code") {
		data["errorMsg"] = invalidLink
	} else {
		email := req.URL.Query()["email"][0]
		code := req.URL.Query()["code"][0]
		date := req.URL.Query()["date"][0]
		data["formAction"] = fmt.Sprintf("%s?action=%s&email=%s&date=%s&code=%s", req.URL.Path, actionName, email, date, code)

		if actionName == "register" {
			if !common.CheckSecret(
				code,
				email,
				date,
			) {
				data["errorMsg"] = invalidLink
			} else {
				// Register user (for real this time)
				data["email"] = email

				if req.Method == "POST" {
					username := req.PostFormValue("username")
					data["username"] = username
					newPassword := req.PostFormValue("newPassword")
					newPasswordAgain := req.PostFormValue("newPasswordAgain")
					if len(username) == 0 || len(newPassword) == 0 {
						data["errorMsg"] = "Required field left blank."
					} else if newPassword != newPasswordAgain {
						data["errorMsg"] = "Passwords do not match."
					} else {
						sessionid, err := databaseActions.Register2(username, email, newPassword)
						if err != nil {
							data["formError"] = err.Error()
						} else { // No error
							common.SetSessionCookie(res, sessionid)

							// Redirect to logged-in wizard
							http.Redirect(res, req, req.URL.Path, http.StatusFound)
							return
						}
					}
				}

				common.ExecTemplate(registerTemplate, res, data)
				return
			}
		} else if actionName == "reset" {
			if !databaseActions.CheckResetLink(
				code,
				email,
				date,
			) {
				data["errorMsg"] = invalidLink
			} else {
				// Reset user's password
				data["successMsg"] = "Valid link."
			}
		}
	}

	common.ExecTemplate(messageTemplate, res, data)
	return
}

func introWizardHandler(res http.ResponseWriter, req *http.Request) {

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

const wizardTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings Wizard</h1>
                {{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
                {{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>{{ template "wizardContent" .}}</section>
			</div>
		</div>
	</div>
	<script src="/static/js/settings_js"></script>
`

const communitiesTemplateText = `
					{ template "communitiesContent" .}}
					<form action="/">
						<button type="submit">Finish</button>
					</form>
				`

const registerTemplateText = `
					<form method="post" action="{{.formAction}}">
						<h2>Password Change</h2>
						<div class="row">
							<div class="large-4 medium-6 columns left">
								<label>
									Email
									<input type="email" value="{{.email}}" disabled>
								</label>
							</div>
						</div>
						<div class="row">
							<div class="large-4 medium-6 columns left">
								<label>
									Username
									<input type="text" name="username" {{if .username}} value="{{.username}}"{{end}}>
								</label>
							</div>
						</div>
						<div class="row">
							<div class="large-4 medium-6 columns left">
								<label>
									New password
									<input type="password" name="newPassword">
								</label>
							</div>
							<div class="large-4 medium-6 columns left">
								<label>
									New password (again)
									<input type="password" name="newPasswordAgain">
								</label>
							</div>
						</div>
						<button type="submit" name="continue" value="true">Continue</button>
					</form>
				`
