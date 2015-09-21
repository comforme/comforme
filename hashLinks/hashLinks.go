package hashLinks

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var messageTemplate *template.Template
var registerTemplate *template.Template
var resetTemplate *template.Template

func init() {
	// Message page template
	messageTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(messageTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(messageTemplate.New("wizardContent").Parse(""))
	template.Must(messageTemplate.New("content").Parse(templates.HashLink))

	// Register page template
	registerTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(registerTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(registerTemplate.New("wizardContent").Parse(registerTemplateText))
	template.Must(registerTemplate.New("content").Parse(templates.HashLink))

	// Reset page template
	resetTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(resetTemplate.New("nav").Parse(templates.NavlessBar))
	template.Must(resetTemplate.New("wizardContent").Parse(resetTemplateText))
	template.Must(resetTemplate.New("content").Parse(templates.HashLink))
}

func ResetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("cache-control", "private, max-age=0, no-cache")

	data := map[string]interface{}{}
	data["pageTitle"] = "Password Reset"

	if !common.CheckParam(req.URL.Query(), "email") ||
		!common.CheckParam(req.URL.Query(), "date") ||
		!common.CheckParam(req.URL.Query(), "code") {
		data["errorMsg"] = common.InvalidLink.Error()
	} else {
		email := req.URL.Query()["email"][0]
		data["email"] = email
		code := req.URL.Query()["code"][0]
		date := req.URL.Query()["date"][0]
		data["formAction"] = fmt.Sprintf("%s?email=%s&date=%s&code=%s", req.URL.Path, email, date, code)

		if !databaseActions.CheckResetLink(
			code,
			email,
			date,
		) {
			data["errorMsg"] = common.InvalidLink.Error()
		} else {
			// Reset user's password
			if req.Method == "POST" {
				newPassword := req.PostFormValue("newPassword")
				newPasswordAgain := req.PostFormValue("newPasswordAgain")

				if len(newPassword) == 0 {
					data["errorMsg"] = "Required field left blank."
				} else if newPassword != newPasswordAgain {
					data["errorMsg"] = "Passwords do not match."
				} else {
					err := databaseActions.SetPassword(email, newPassword)
					if err != nil {
						data["errorMsg"] = err.Error()
					} else { // No error
						sessionid, err := databaseActions.Login(email, newPassword)
						if err != nil {
							data["errorMsg"] = err.Error()
						} else { // No error
							common.SetSessionCookie(res, sessionid)

							// Redirect to home page
							http.Redirect(res, req, "/", http.StatusFound)
							return
						}
					}
				}
			}
			common.ExecTemplate(resetTemplate, res, data)
			return
		}
	}

	common.ExecTemplate(messageTemplate, res, data)
	return
}

func RegisterHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("cache-control", "private, max-age=0, no-cache")

	data := map[string]interface{}{}
	data["pageTitle"] = "Registration"

	if !common.CheckParam(req.URL.Query(), "email") ||
		!common.CheckParam(req.URL.Query(), "date") ||
		!common.CheckParam(req.URL.Query(), "code") {
		data["errorMsg"] = common.InvalidLink.Error()
	} else {
		email := req.URL.Query()["email"][0]
		data["email"] = email
		code := req.URL.Query()["code"][0]
		date := req.URL.Query()["date"][0]
		data["formAction"] = fmt.Sprintf("%s?email=%s&date=%s&code=%s", req.URL.Path, email, date, code)

		if !databaseActions.CheckRegisterLink(
			code,
			email,
			date,
		) {
			data["errorMsg"] = common.InvalidLink.Error()
		} else {
			// Register user (for real this time)

			if req.Method == "POST" {
				username := req.PostFormValue("username")
				data["username"] = username
				newPassword := req.PostFormValue("newPassword")
				newPasswordAgain := req.PostFormValue("newPasswordAgain")
				if len(username) == 0 || len(newPassword) == 0 {
					data["errorMsg"] = "Required field(s) left blank."
				} else if newPassword != newPasswordAgain {
					data["errorMsg"] = "Passwords do not match."
				} else {
					sessionid, err := databaseActions.Register2(username, email, newPassword)
					if err != nil {
						data["errorMsg"] = err.Error()
					} else { // No error
						common.SetSessionCookie(res, sessionid)

						// Redirect to tour
						http.Redirect(res, req, "/tour", http.StatusFound)
						return
					}
				}
			}

			common.ExecTemplate(registerTemplate, res, data)
			return
		}
	}

	common.ExecTemplate(messageTemplate, res, data)
	return
}

const registerTemplateText = `					<form method="post" action="{{.formAction}}">
						<h2>Finish Registering</h2>
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
									Username (between 3 and 20 characters)
									<input type="text" name="username" {{if .username}} value="{{.username}}"{{end}}>
								</label>
							</div>
						</div>
						<div class="row">
							<div class="large-4 medium-6 columns left">
								<label>
									New password (must be at least 6 characters)
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

const resetTemplateText = `					<form method="post" action="{{.formAction}}">
						<h2>Password Reset</h2>
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
						<button type="submit" name="resetPassword" value="true">Update Password</button>
					</form>
				`
