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

	communities, err := databaseActions.ListCommunities(sessionid)
	if err != nil {
		log.Println("Error listing communities:", err)
		common.Logout(res, req)
		return
	}
	
	perCol := len(communities) / 4
	extra := len(communities) % 4
	cut1 := perCol
	if extra >= 1 {
		cut1++
	}
	cut2 := cut1 + perCol
	if extra >= 2 {
		cut2++
	}
	cut3 := cut2 + perCol
	if extra >= 3 {
		cut3++
	}
	data["communitiesCols"] = [][]common.Community{
		communities[0:cut1],
		communities[cut1:cut2],
		communities[cut2:cut3],
		communities[cut3:],
	}

	openSessions, err := databaseActions.OtherSessions(sessionid)
	if err != nil {
		log.Println("Error listing communities:", err)
		common.Logout(res, req)
		return
	} else {
		data["openSessions"] = openSessions
	}

	if req.Method == "POST" {
		//username := req.PostFormValue("username")
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
				data["errorMsg"] = "Failed to validate password."
			}
		} else {
			data["errorMsg"] = "Passwords do not match"
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
				<h1><i class="fi-widget"></i> Settings</h1>
                {{if .successMsg}}<div class="alert-box success">{{.successMsg}}</div>{{end}}
                {{if .errorMsg}}<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<section>
					<h2>User Information</h2>
					<div class="row">
						<div class="large-4 columns left">
							<h6>Email:</h6> {{.email}}
						</div>
						<div class="large-4 columns left">
							<h6>Username:</h6> {{.username}}
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
					<h2>Your Communities</h2>
					<h6>Check all that apply.</h6>
					<div class="row">{{range $col_number, $communitiesCol := $.communitiesCols}}
						<div class="large-3 medium-6 small-12 columns left">{{range $line_number, $community := $communitiesCol}}
							<div>
								<label>
									<input class="communityCheckbox" type="checkbox" name="{{$community.Id}}"{{if eq $community.IsMember true}} checked="checked"{{end}} value="{{$community.Name}}">
									{{$community.Name}}
								</label>
							</div>{{end}}
						</div>{{end}}
					</div>
				</section>
				<section>
					<h2>Sessions</h2>
					<h6>You currently have <span id="numOpenSessions">{{.openSessions}}</span> sessions open in addition to this one.</h6>
					<button onclick="logoutOtherSessions(this)" name="logout-sessions">Logout Other Sessions</button>
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
