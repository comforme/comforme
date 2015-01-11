package login

import (
	"fmt"
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
	template.Must(loginTemplate.New("content").Parse(loginTemplateText))
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	var err error

	if req.Method == "POST" {
		email := req.PostFormValue("email")
		username := req.PostFormValue("username")
		password := req.PostFormValue("password")
		isSignup := req.PostFormValue("sign-up") == "true"
		isLogin := req.PostFormValue("log-in") == "true"

		data["username"] = username
		data["email"] = email
		data["pageTitle"] = "login"

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
			sessionid, err = databaseActions.Login(email, password)
			if err == database.InvalidUsernameOrPassword {
				data["loginUsernameError"] = err.Error()
			} else if err != nil {
				log.Println("Unknown signup error:", err)
				data["formError"] = "Unknown signup error. Check error log."
			}
		}

		if err == nil {
			common.SetSessionCookie(req, sessionid)
			fmt.Fprintln(res, "Success!")
			return
		}
	}

	common.ExecTemplate(loginTemplate, res, data)
}

const loginTemplateText = `
    <div class="content">
		<div class="row">
			<div class="large-4 medium-3 small-1 columns">&nbsp;</div>
			<div class="large-4 medium-6 small-10 columns">{{if .formError}}
				<div class="alert-box alert">
					{{.formError}}
				</div>{{end}}
				<section class="login-tabs">
					<dl class="tabs" data-tab>
						<dd class="active"><a href="#sign-up-form">Sign Up</a></dd>
						<dd><a href="#log-in-form">Log In</a></dd>
					</dl>
					<div class="tabs-content">
						<div class="content active" id="sign-up-form">
							<form method="post" action="/">
								<div>
									<input type="text" name="username" placeholder="User Name"{{if .username}} value="{{.username}}"{{end}}>{{if .registerUsernameError}}
									<small class="error">{{.registerUsernameError}}</small>{{end}}
								</div>
								<div>
									<input type="email" name="email" placeholder="Email"{{if .email}} value="{{.email}}"{{end}}>{{if .registerEmailError}}
									<small class="error">{{.registerEmailError}}</small>{{end}}
								</div>
								<div>
									<button type="submit" class="button" name="sign-up" value="true">Submit</button>
								</div>
							</form>
						</div>
						<div class="content" id="log-in-form">
							<form method="post" action="/">
								<div>
									<input type="email" name="email" placeholder="Email"{{if .email}} value="{{.email}}"{{end}}>{{if .loginEmailError}}
									<small class="error">{{.loginEmailError}}</small>{{end}}
								</div>
								<div>
									<input type="password" name="password" placeholder="Password">
								</div>
								<div>
									<button type="submit" class="button" name="log-in" value="true">
								</div>
							</form>
						</div>
					</div>
				</section>
			</div>
			<div class="large-4 medium-3 small-1 columns">&nbsp;</div>
		</div>
	</div>
`
