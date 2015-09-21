package login

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/recaptcha"
	"github.com/comforme/comforme/templates"
)

var loginTemplate *template.Template
var recaptchaPublicKey string

const registrationSuccess = "Registration successful. Please check your email."
const resetSuccess = "Password reset successful. Please check your email."

func init() {
	loginTemplate = template.Must(template.New("siteLayout").Parse(templates.SiteLayout))
	template.Must(loginTemplate.New("nav").Parse(""))
	template.Must(loginTemplate.New("content").Parse(loginTemplateText))
	recaptchaPublicKey = os.Getenv("RECAPTCHA_PUBLIC_KEY")
	recaptcha.Init(os.Getenv("RECAPTCHA_PRIVATE_KEY"))
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	data["formAction"] = req.URL.Path
	data["pageTitle"] = "Login"
	data["recaptchaPublicKey"] = recaptchaPublicKey
	data["siteName"] = common.SiteName

	if req.Method == "POST" {
		isSignup := req.PostFormValue("sign-up") == "true"
		isLogin := req.PostFormValue("log-in") == "true"
		isReset := req.PostFormValue("reset-password") == "true"

		email := req.PostFormValue("email")
		data["email"] = email

		if isSignup {
			// Check ReCaptcha
			ipAddress := common.GetIpAddress(req)
			recaptchaResponse := req.PostFormValue("g-recaptcha-response")
			log.Println("recaptchaResponse", recaptchaResponse)
			log.Println("ipAddress", ipAddress)
			err := recaptcha.Check(
				recaptchaResponse,
				ipAddress,
			)
			if err != nil {
				log.Println("reCaptcha failed:", err)
				data["formError"] = err.Error()
			} else {
				log.Println("reCaptcha success:", err)
				err = databaseActions.Register1(email, common.GetBaseURL(req))
				if err != nil {
					data["formError"] = err.Error()
				} else { // No error
					data["successMsg"] = registrationSuccess
				}
			}
		} else if isLogin {
			data["loginSelected"] = "true"

			password := req.PostFormValue("password")

			sessionid, err := databaseActions.Login(email, password)
			if err != nil {
				data["formError"] = err.Error()
			} else { // No error
				common.SetSessionCookie(res, sessionid)

				// Redirect to intended page
				http.Redirect(res, req, req.URL.Path, http.StatusFound)
				return // Not needed, may reduce load on server
			}
		} else if isReset {
			// Check ReCaptcha
			ipAddress := common.GetIpAddress(req)
			recaptchaResponse := req.PostFormValue("g-recaptcha-response")
			log.Println("recaptchaResponse", recaptchaResponse)
			log.Println("ipAddress", ipAddress)
			err := recaptcha.Check(
				recaptchaResponse,
				ipAddress,
			)
			if err != nil {
				log.Println("reCaptcha failed:", err)
				data["formError"] = err.Error()
			} else {
				log.Println("reCaptcha success:", err)
				err := databaseActions.ResetPassword(email, common.GetBaseURL(req))
				if err != nil {
					data["formError"] = err.Error()
				} else {
					data["successMsg"] = resetSuccess
				}
			}
		}
	}

	common.ExecTemplate(loginTemplate, res, data)
}

const loginTemplateText = `
    <div class="content sign-up-and-log-in">
		<h1 class="text-center">Welcome to Community for Me!</h1>
		<div class="row">
			<div class="large-3 medium-3 show-for-medium-up columns">&nbsp;</div>
			<div class="large-6 medium-6 columns" style="min-width: 320px;">{{if .formError}}
				<div class="alert-box alert">
					{{.formError}}
				</div>{{end}}{{if .successMsg}}
				<div class="alert-box success">
					{{.successMsg}}
				</div>{{end}}
				<section class="login-tabs sign-up-and-log-in">
					<dl class="tabs" data-tab>
						<dd{{if not .loginSelected}} class="active"{{end}}><a href="#sign-up-form">Sign Up</a></dd>
						<dd{{if .loginSelected}} class="active"{{end}}><a href="#log-in-form">Log In</a></dd>
					</dl>
					<div class="tabs-content">
						<div class="content{{if not .loginSelected}} active{{end}}" id="sign-up-form">
							<form method="post" action="{{.formAction}}">
								<noscript>
									<small class="error">This site requires JavaScript to function!</small>
								</noscript>
								<div{{if .registerEmailError}} class="error"{{end}}>
									<input type="email" name="email" placeholder="Email"{{if .email}} value="{{.email}}"{{end}}>{{if .registerEmailError}}
									<small class="error">{{.registerEmailError}}</small>{{end}}
								</div>
								<div class="g-recaptcha" data-sitekey="{{.recaptchaPublicKey}}"></div>
								<div>
									<button type="submit" class="button expand" name="sign-up" value="true">Sign Up</button>
									<button type="submit" class="button tiny expand" name="reset-password" value="true">Reset Password</button>
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
									<button type="submit" class="button expand" name="log-in" value="true">Log In</button>
								</div>
							</form>
						</div>
					</div>
				</section>
			</div>
			<div class="large-2 medium-2 show-for-medium-up columns">&nbsp;</div>
			<div class="large-12 columns">
				<h2>What is {{.siteName}}?</h2>
					<div>
						<p>{{.siteName}} (Community for Me) is a community-rated and identity-oriented social network/service listing. Users can find accepting communities and services based on a wide array of keywords. Users can also start their own communities categorized by aforementioned keywords. {{.siteName}} makes it easier for an individual to find communities and services which accept them for who they are.</a></p>
					</div>
				</div>
			</div>
		</div>
	</div>
	<a href="https://github.com/comforme/comforme"><img style="position: absolute; top: 0; right: 0; border: 0;" src="https://camo.githubusercontent.com/38ef81f8aca64bb9a64448d0d70f1308ef5341ab/68747470733a2f2f73332e616d617a6f6e6177732e636f6d2f6769746875622f726962626f6e732f666f726b6d655f72696768745f6461726b626c75655f3132313632312e706e67" alt="Fork me on GitHub" data-canonical-src="https://s3.amazonaws.com/github/ribbons/forkme_right_darkblue_121621.png"></a>
`
