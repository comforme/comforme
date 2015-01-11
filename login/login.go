package login

import (
	"html/template"
	"log"
	"net/http"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/database"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/templates"
)

var loginTemplate *template.Template

func init() {
	loginTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(loginTemplate.New("nav").Parse(""))
	template.Must(loginTemplate.New("content").Parse(loginTemplateText))
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	var err error

	data["formAction"] = req.URL.Path
	data["pageTitle"] = "login"

	if req.Method == "POST" {
		email := req.PostFormValue("email")
		username := req.PostFormValue("username")
		password := req.PostFormValue("password")
		isSignup := req.PostFormValue("sign-up") == "true"
		isLogin := req.PostFormValue("log-in") == "true"

		data["username"] = username
		data["email"] = email

		var sessionid string

		if isSignup {
			sessionid, err = databaseActions.Register(username, email)
			if err == databaseActions.UsernameTooShort {
				data["registerUsernameError"] = err.Error()
			} else if err == databaseActions.InvalidEmail {
				data["registerEmailError"] = err.Error()
			} else if err != nil {
				log.Println("Unknown signup error:", err)
				data["formError"] = "Unknown signup error. Check error log."
			}
		} else if isLogin {
			data["loginSelected"] = "true"
			sessionid, err = databaseActions.Login(email, password)
			if err == database.InvalidUsernameOrPassword {
				data["loginError"] = err.Error()
			} else if err != nil {
				log.Println("Unknown signup error:", err)
				data["formError"] = "Unknown signup error. Check error log."
			}
		}

		if err == nil {
			common.SetSessionCookie(res, sessionid)
			
			// Redirect to home page
			http.Redirect(res, req, "/", http.StatusFound)
		}
	}

	common.ExecTemplate(loginTemplate, res, data)
}

const loginTemplateText = `
    <div class="content sign-up-and-log-in">
		<h1 class="text-center">Welcome to Community for Me!<h1>
		<div class="row">
			<div class="large-4 medium-3 small-1 columns">&nbsp;</div>
			<div class="large-4 medium-6 small-10 columns">{{if .formError}}
				<div class="alert-box alert">
					{{.formError}}
				</div>{{end}}
				<section class="login-tabs sign-up-and-log-in">
					<dl class="tabs" data-tab>
						<dd{{if not .loginSelected}} class="active"{{end}}><a href="#sign-up-form">Sign Up</a></dd>
						<dd{{if .loginSelected}} class="active"{{end}}><a href="#log-in-form">Log In</a></dd>
					</dl>
					<div class="tabs-content">
						<div class="content{{if not .loginSelected}} active{{end}}" id="sign-up-form">
							<form method="post" action="{{.formAction}}">
								<div{{if .registerUsernameError}} class="error"{{end}}>
									<input type="text" name="username" placeholder="User Name"{{if .username}} value="{{.username}}"{{end}}>{{if .registerUsernameError}}
									<small class="error">{{.registerUsernameError}}</small>{{end}}
								</div>
								<div{{if .registerEmailError}} class="error"{{end}}>
									<input type="email" name="email" placeholder="Email"{{if .email}} value="{{.email}}"{{end}}>{{if .registerEmailError}}
									<small class="error">{{.registerEmailError}}</small>{{end}}
								</div>
								<div>
									<button type="submit" class="button" name="sign-up" value="true">Submit</button>
								</div>
							</form>
						</div>
						<div class="content{{if .loginSelected}} active{{end}}" id="log-in-form">
							<form method="post" action="{{.formAction}}">
								<div{{if .loginError}} class="error"{{end}}>
									<input type="email" name="email" placeholder="Email"{{if .email}} value="{{.email}}"{{end}}>{{if .loginError}}
									<small class="error">{{.loginError}}</small>{{end}}
								</div>
								<div{{if .loginError}} class="error"{{end}}>
									<input type="password" name="password" placeholder="Password">{{if .loginError}}
									<small class="error">{{.loginError}}</small>{{end}}
								</div>
								<div>
									<button type="submit" class="button" name="log-in" value="true">Submit</button>
								</div>
							</form>
						</div>
					</div>
				</section>
			</div>
			<div class="large-4 medium-3 small-1 columns">&nbsp;</div>
			<div class="large-12 columns">
				<h2>What is Comfor.me?</h2>
					<div>
						<p>Comfor.me (Community for Me) is a community-rated and identity-oriented social network/service listing. Users can find accepting communities and services based on a wide array of keywords. Users can also start their own communities categorized by aforementioned keywords. Comfor.me makes it easier for an individual to find communities and services which accept them for who they are.</a></p>
					</div>
				</div>
			</div>
		</div>
	</div>
`
