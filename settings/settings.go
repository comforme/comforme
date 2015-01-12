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

	cookie, err := req.Cookie("sessionid")
	if err != nil {
		log.Println("Failed to retrieve sessionid:", err)
		common.Logout(res, req)
		return
	}
	sessionid := cookie.Value

	communities, err := databaseActions.ListCommunities(sessionid)
	if err != nil {
		log.Println("Error listing communities:", err)
	} else {
		log.Printf("communities: %+v\n", communities)
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
		data["communitiesCol1"] = communities[0:cut1]
		data["communitiesCol2"] = communities[cut1:cut2]
		data["communitiesCol3"] = communities[cut2:cut3]
		data["communitiesCol4"] = communities[cut3:]
		
		log.Printf("communitiesCol1: %+v\n", data["communitiesCol1"])
		log.Printf("communitiesCol2: %+v\n", data["communitiesCol2"])
		log.Printf("communitiesCol3: %+v\n", data["communitiesCol3"])
		log.Printf("communitiesCol4: %+v\n", data["communitiesCol4"])
	}

	if req.Method == "POST" {
		//username := req.PostFormValue("username")
		oldPassword := req.PostFormValue("oldPassword")
		newPassword := req.PostFormValue("newPassword")
		newPasswordAgain := req.PostFormValue("newPasswordAgain")
		if len(oldPassword) != 0 || len(newPassword) != 0 {
			if newPassword == newPasswordAgain {
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
				<form action="{{.formAction}}" method="post">
					<section>
						<h2>Password Change</h2>
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
