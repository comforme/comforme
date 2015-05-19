package pages

import (
	"fmt"
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

func PageHandler(res http.ResponseWriter, req *http.Request, sessionid, email, username string, userID int) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path

	category := bone.GetValue(req, "category")
	slug := bone.GetValue(req, "slug")

	log.Printf("Looking up page with category (%s) and slug (%s)...\n", category, slug)
	page, err := databaseActions.GetPage(category, slug)
	if err != nil {
		http.NotFound(res, req)
		log.Printf("Error looking up page (%s): %s\n", req.URL.Path, err.Error())
		return
	}

	data["page"] = page

	if req.Method == "POST" {
		thoughts := req.PostFormValue("post-your-thoughts")
		if len(thoughts) < common.MinDescriptionLength {
			data["errorMsg"] = fmt.Sprintf("Post must be at least %d characters long.", common.MinDescriptionLength)
			data["thoughts"] = thoughts
			goto renderPosts
		}

		err = databaseActions.CreatePost(sessionid, thoughts, page)

		if err == nil {
			data["successMsg"] = "Post successfully added."
		} else {
			data["errorMsg"] = err.Error()
			data["thoughts"] = thoughts
		}
	}

renderPosts:
	log.Printf("Looking up posts for page id (%d)...\n", page.Id)
	posts, err := databaseActions.GetPosts(sessionid, page)
	if err != nil {
		http.NotFound(res, req)
		log.Printf("Error looking up posts for page (%d): %s\n", page.Id, err.Error())
		return
	}

	data["posts"] = posts

	common.ExecTemplate(pageTemplate, res, data)
}

const pageTemplateText = `
	<div class="content">
		<div class="row">
			<div class="columns">
				<h1><a href="{{.formAction}}">{{.page.Title}}</a></h1>
				<p>
					{{.page.Description}}
				</p>
			</div>
		</div>{{if .page.Address}}
		<div class="row">
			<div class="columns">
				<p>
					<strong>Address:</strong> <span>{{.page.Address}}</span>
				</p>
			</div>
		</div>{{end}}{{if .page.Website}}
		<div class="row">
			<div class="columns">
				<p>
					<strong>Website:</strong> <span>{{.page.Website}}</span>
				</p>
			</div>
		</div>{{end}}
		<div class="row">
			<div class="columns">{{if .successMsg}}
				<div class="alert-box success">{{.successMsg}}</div>{{end}}{{if .errorMsg}}
				<div class="alert-box alert">{{.errorMsg}}</div>{{end}}
				<form method="post" action="{{.action}}">
					<fieldset>
						<legend>
							Post Your Thoughts
						</legend>
						<div class="row">
							<div class="columns">
								<label for="post-your-thoughts">Comment:</label>
								<textarea name="post-your-thoughts" id="post-your-thoughts">{{if .thoughts}}{{.thoughts}}{{end}}</textarea>
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
		<div class="row">{{range $post_number, $post := $.posts}}
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
