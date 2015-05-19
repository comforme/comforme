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
	template.Must(settingsTemplate.New("communities").Parse(templates.Communities))
	//template.Must(settingsTemplate.New("content").Parse(settingsTemplateText))
	template.Must(settingsTemplate.New("content").Parse(settingsTemplateText))
}

func SettingsHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path
	data["pageTitle"] = "Settings"

	cookie, err := req.Cookie("sessionid")
	if err != nil {
		log.Println("Failed to retrieve sessionid:", err)
		common.Logout(res, req)
		return
	}
	sessionid := cookie.Value

	email, err := databaseActions.GetEmail(sessionid)
	if err != nil {
		log.Println("Error getting email:", err)
		common.Logout(res, req)
		return
	}
	data["email"] = email

	username, err := databaseActions.GetUsername(sessionid)
	if err != nil {
		log.Println("Error getting username:", err)
		common.Logout(res, req)
		return
	}
	data["username"] = username

	data["communitiesCols"], err = databaseActions.GetCommunityColumns(sessionid)
	if err != nil {
		log.Println("Error listing communities:", err)
		common.Logout(res, req)
		return
	}

	openSessions, err := databaseActions.OtherSessions(sessionid)
	if err != nil {
		log.Println("Error getting open sessions:", err)
		common.Logout(res, req)
		return
	} else {
		data["openSessions"] = openSessions
	}

	if req.Method == "POST" {
		if req.PostFormValue("password-update") == "true" {
			oldPassword := req.PostFormValue("oldPassword")
			newPassword := req.PostFormValue("newPassword")
			newPasswordAgain := req.PostFormValue("newPasswordAgain")
			if len(oldPassword) == 0 || len(newPassword) == 0 {
				data["errorMsg"] = "Both old and new password required to change password"
			} else if newPassword == newPasswordAgain {
				err := databaseActions.ChangePassword(sessionid, oldPassword, newPassword)
				if err == nil {
					data["successMsg"] = "Password changed."
					if req.URL.Path != "/settings" {
						http.Redirect(res, req, req.URL.Path, http.StatusFound)
						return
					}
				} else {
					data["errorMsg"] = err.Error()
				}
			} else {
				data["errorMsg"] = "Passwords do not match."
			}
		} else if req.PostFormValue("username-update") == "true" {
			usernameChangePassword := req.PostFormValue("usernameChangePassword")
			newUsername := req.PostFormValue("newUsername")

			err := databaseActions.ChangeUsername(sessionid, newUsername, usernameChangePassword)
			if err != nil {
				data["newUsername"] = newUsername
				data["errorMsg"] = err.Error()
			} else {
				data["successMsg"] = "Username changed."
				data["username"] = newUsername
			}
		}
	}

	if data["errorMsg"] == nil {
		cookie, err := req.Cookie("sessionid")
		if err != nil {
			log.Println("Error reading cookie:", err)
			common.Logout(res, req)
		}

		isRequired, err := databaseActions.PasswordChangeRequired(cookie.Value)
		if err != nil {
			log.Println("Error checking if password reset is required:", err)
			common.Logout(res, req)
		}

		if isRequired {
			data["errorMsg"] = "Password change required."
		}
	}

	common.ExecTemplate(settingsTemplate, res, data)
}

// TODO replace uppercase placeholders below
const settingsTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings</h1>{{if .successMsg}}
				<div class="alert-box success">{{.successMsg}}</div>{{end}}{{if .errorMsg}}
				<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>
					<h2>User Information</h2>
					<div class="row">
						<div class="large-4 columns left">
							<h5>Email:</h5> {{.email}}
						</div>
						<div class="large-4 columns left">
							<h5>Username:</h5> {{.username}}
						</div>
					</div>
				</section>
				<section>
					<h2>Password Change</h2>
					<form action="{{.formAction}}" method="post">
						<div class="row">
							<div class="large-4 columns left">
								<label>
									Old password (Initial password sent via email)
									<input type="password" name="oldPassword">
								</label>
							</div>
							<div class="large-4 columns left">
								<label>
									New password
									<input type="password" name="newPassword">
								</label>
							</div>
							<div class="large-4 columns left">
								<label>
									New password (again)
									<input type="password" name="newPasswordAgain">
								</label>
							</div>
						</div>
						<button type="submit" name="password-update" value="true">Update Password</button>
					</form>
				</section>
				<section>
{{ template "communities" .}}
				</section>
				<section>
					<h2>Sessions</h2>
					<h6>You currently have <span id="numOpenSessions">{{.openSessions}}</span> sessions open in addition to this one.</h6>
					<button onclick="logoutOtherSessions(this)" name="logout-sessions">Logout Other Sessions</button>
				</section>
				<section>
					<h2>Username Change</h2>
					<form action="{{.formAction}}" method="post">
						<div class="row">
							<div class="large-4 columns left">
								<label>
									Password
									<input type="password" name="usernameChangePassword">
								</label>
							</div>
							<div class="large-4 columns left">
								<label>
									New username
									<input type="text" name="newUsername"{{if .newUsername}} value="{{.newUsername}}"{{end}}>
								</label>
							</div>
						</div>
						<button type="submit" name="username-update" value="true">Update Username</button>
					</form>
				</section>{{/*
				<section>
					<div class="row">
						<div class="columns">
							<h2>Find Communities</h2>
							{{template "communitySearch" . }}
						</div>
					</div>
				</section>*/}}
			</div>
		</div>
	</div>
	<script src="/static/js/settings_js"></script>
`
