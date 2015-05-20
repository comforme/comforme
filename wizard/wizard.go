package wizard

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/requireLogin"
	"github.com/comforme/comforme/templates"
)

var communitiesTemplate *template.Template

func init() {
	// Community selection page template
	communitiesTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(communitiesTemplate.New("nav").Parse(templates.NavBar))
	template.Must(communitiesTemplate.New("content").Parse(templates.Wizard))
	template.Must(communitiesTemplate.New("wizardContent").Parse(communitiesTemplateText))
	template.Must(communitiesTemplate.New("communitiesContent").Parse(templates.Communities))
}

func WizardHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params, userInfo common.UserInfo) {

	data := map[string]interface{}{}

	var err error
	data["communitiesCols"], err = databaseActions.GetCommunityColumns(userInfo.UserID)
	if err != nil {
		log.Println("Error listing communities:", err)
		common.Logout(res, req)
		return
	}

	common.ExecTemplate(communitiesTemplate, res, data)
}

const communitiesTemplateText = `					<section class="wizard-tabs">
						<ul class="tabs small-block-grid-2 medium-block-grid-3 xlarge-block-grid-6" data-tab role="tablist">
							<li class="tab-title active" role="presentational"><a href="#tab-communities" role="tab" tabindex="0" aria-selected="false" controls="tab-communities" id="tab-communities-navigation">Your Communities</a></li>
							<li class="tab-title" role="presentational"><a href="#tab-nav-bar" role="tab" tabindex="0" aria-selected="true" controls="tab-nav-bar" id="tab-nav-bar-navigation">The Nav Bar</a></li>
							<li class="tab-title" role="presentational"><a href="#tab-finding-resources" role="tab" tabindex="0" aria-selected="false" controls="tab-finding-resources" id="tab-finding-resources-navigation">Finding Resources</a></li>
							<li class="tab-title" role="presentational"><a href="#tab-reviews" role="tab" tabindex="0" aria-selected="false" controls="tab-reviews" id="tab-reviews-navigation">Reviews</a></li>
							<li class="tab-title" role="presentational"><a href="#tab-adding-resources" role="tab" tabindex="0" aria-selected="false" controls="tab-adding-resources" id="tab-adding-resources-navigation">Adding Resources</a></li>
							<li class="tab-title" role="presentational"><a href="#tab-settings" role="tab" tabindex="0" aria-selected="false" controls="tab-settings" id="tab-settings-navigation">Changing Settings</a></li>
						</ul>
					</section>
					<div class="tabs-content">
						<section role="tabtab" aria-hidden="false" class="content active" id="tab-communities">
							{{ template "communitiesContent" .}}
							<button onClick="location.href='#tab-nav-bar'">Next</button>
						</section>
						<section role="tabtab" aria-hidden="true" class="content" id="tab-nav-bar">
							<h2>The Nav Bar</h2>
							<p>The nav bar is at the top of every page when you are logged in.</p>
							<button onClick="location.href='#tab-finding-resources'">Next</button>
						</section>
						<section role="tabtab" aria-hidden="true" class="content" id="tab-finding-resources">
							<h2>Finding Resources</h2>
							<p>Search for resources on the home page.</p>
							<button onClick="location.href='#tab-reviews'">Next</button>
						</section>
						<section role="tabtab" aria-hidden="true" class="content" id="tab-reviews">
							<h2>Reviews</h2>
							<p>View a resource to read its reviews or post your own. Reviews are sorted by how many communities you have in common with the reviewer.</p>
							<button onClick="location.href='#tab-adding-resources'">Next</button>
						</section>
						<section role="tabtab" aria-hidden="true" class="content" id="tab-adding-resources">
							<h2>Adding Resources</h2>
							<p>Add new resources on the add new resource page.</p>
							<button onClick="location.href='#tab-settings'">Next</button>
						</section>
						<section role="tabtab" aria-hidden="true" class="content" id="tab-settings">
							<h2>Changing Settings</h2>
							<p>Change your username or password, or update your communities on the settings page.</p>
							<button onclick="location.href='/'">Finish</button>
						</section>
					</div>
				`
