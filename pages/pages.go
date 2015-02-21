package pages

import (
	"html/template"
	"log"
	"net/http"
	
	"github.com/go-zoo/bone"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var pageTemplate *template.Template

func init() {
	pageTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(pageTemplate.New("nav").Parse(templates.NavBar))
	template.Must(pageTemplate.New("fullPage").Parse(templates.SearchBar))
	template.Must(pageTemplate.New("content").Parse(pageTemplateText))
}

func PageHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path
	
	cookie, err := req.Cookie("sessionid")
	if err != nil {
		log.Println("Failed to retrieve sessionid:", err)
		common.Logout(res, req)
		return
	}
	sessionid := cookie.Value
	
	category := bone.GetValue(req, "category")
	slug := bone.GetValue(req, "slug")
	
	//page, posts, err :=
	_, _, err = databaseActions.GetPage(sessionid, category, slug)
	
	if err != nil {
		http.NotFound(res, req)
		log.Printf("Error looking up page (%s): %s\n", req.URL.Path, err.Error())
		return
	}

	common.ExecTemplate(pageTemplate, res, data)
}

const pageTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns">
				<h1><a href="{{.action}}">{{.title}}</a></h1>
				<p>
					{{.description}}
				</p>
			</div>
		</div>
		<div class="row">
			<div class="columns">
				<p>
					<strong>Address:</strong> <span>{{.address}}</span>
				</p>
			</div>
		</div>
		<div class="row">
			<div class="columns">
				<form method="post" action="{{.action}}">
					<fieldset>
						<legend>
							Post Your Thoughts
						</legend>
						<div class="row">
							<div class="columns">
								<label for="post-your-thoughts">Comment:</label>
								<textarea name="post-your-thoughts" id="post-your-thoughts"></textarea>
							</div>
						</div>
						<div class="row">
							<div class="columns text-right">
								<button type="submit">Comment</button>
							</div>
						</div>
					</fieldset>
				</form>
			</div>
		</div>
		<div class="row">{{range $post_number, $post := $.communitiesCols}}
			<div class="columns">
				<p>
					<strong>
						{{$post.Author}} ({{$post.CommonCategories}})
					</strong>
					<small>
						{{$post.Date}}
					</small>
				</p>
				<p>
					{{$post.Body}}
				</p>
			</div>{{end}}
		</div>
	</div>
`
