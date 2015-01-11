package settings

import (
	"html/template"
	"log"
	"net/http"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var settingsTemplate *template.Template

func init() {
	settingsTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(settingsTemplate.New("nav").Parse(templates.NavBar))
	template.Must(settingsTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(settingsTemplate.New("communitySearch").Parse(templates.CommunitySearch))
	//template.Must(settingsTemplate.New("content").Parse(settingsTemplateText))
	template.Must(settingsTemplate.New("content").Parse(settingsTemplateText))
}

func SettingsHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path

	if req.Method == "POST" {
		//username := req.PostFormValue("username")
		oldPassword := req.PostFormValue("oldPassword")
		newPassword := req.PostFormValue("newPassword")
		newPasswordAgain := req.PostFormValue("newPasswordAgain")
		if len(oldPassword) != 0 || len(newPassword) != 0 {
			if newPassword == newPasswordAgain {
				cookie, err := req.Cookie("sessionid")
				if err == nil {
					sessionId := cookie.Value
					err := databaseActions.ChangePassword(sessionId, oldPassword, newPassword)
					if err == nil {
						data["successMsg"] = "Password changed."
					} else {
						data["errorMsg"] = "Failed to validate password."
					}
				} else {
					log.Println("Failed to retrieve sessionid:", err)
					data["errorMsg"] = "Failed to update password"
				}
			} else {
				data["errorMsg"] = "Passwords do not match"
			}
		}
	}

	common.ExecTemplate(settingsTemplate, res, data)
}

// TODO replace uppercase placeholders below
const settingsTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings</h1>
                {{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
                {{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<form action="{{.formAction}}" method="post">
					<section>
						<h2>Password Change</h2>
						<div class="row">
							<div class="large-4 columns left">
								<label>
									<input type="password" name="oldPassword">
									Old password (Initial password sent via email)
								</label>
							</div>
							<div class="large-4 columns left">
								<label>
									<input type="password" name="newPassword">
									New password
								</label>
							</div>
							<div class="large-4 columns left">
								<label>
									<input type="password" name="newPasswordAgain">
									New password (again)
								</label>
							</div>
						</div>
					</section>
					<section>
						<h2>Your Communities</h2>
						<div class="row">
							<div class="large-3 medium-6 small-12 columns left">
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
							</div>
							<div class="large-3 medium-6 small-12 columns left">
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
							</div>
							<div class="large-3 medium-6 small-12 columns left">
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
							</div>
							<div class="large-3 medium-6 small-12 columns left">
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
								<div>
									<label>
										<input type="checkbox" name="WHAT-TO-NAME-THIS" checked="checked" value="NAME-OR-ID-OF-CATEGORY-REPLACE-ME">
										NAME-OF-CATEGORY
									</label>
								</div>
							</div>
						</div>
						<button type="submit" name="user-communites-update" value="true">Update</button>
					</section>
					<section>
						<div class="row">
							<div class="columns">
								<h2>Find Communities</h2>
								{{template "communitySearch" . }}
							</div>
						</div>
					</section>
				</form>
			</div>
		</div>
	</div>
`
