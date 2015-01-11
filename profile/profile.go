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
	template.Must(profileTemplate.New("nav").Parse(templates.NavBar))
	template.Must(profileTemplate.New("searchBar").Parse(templates.SearchBar))
	template.Must(profileTemplate.New("communitySearch").Parse(templates.CommunitySearch))
	//template.Must(profileTemplate.New("content").Parse(profileTemplateText))
	template.Must(profileTemplate.New("content").Parse(settingsTemplateText))
}

func ProfileHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path

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

const profileTemplateText = `
<div class="content">
	<p>Profile Page</p>
</div>
`

// TODO replace uppercase placeholders below
const settingsTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns communities-settings">
				<h1><i class="fi-widget"></i> Settings</h1>
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
