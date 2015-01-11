package requireLogin

import (
	"log"
	"net/http"

	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
	"github.com/comforme/comforme/settings"
)

const DebugMode = false

func RequireLogin(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if DebugMode {
			log.Printf("Entering debug mode...")
			handler(res, req)
			return
		}
		cookie, err := req.Cookie("sessionid")
		if err == nil {
			email, err := databaseActions.GetEmail(cookie.Value)
			if err == nil {
				log.Printf("User with email %s logged in.", email)
				isRequired, err := databaseActions.PasswordChangeRequired(cookie.Value)
				if err == nil {
					if isRequired {
						settings.SettingsHandler(res, req)
						return
					} else {
						handler(res, req)
						return
					}
				} else {
					log.Println("Error checking if password reset is required:", err)
				}
			} else {
				log.Println("Error checking email:", err)

				// Delete bad cookie
				cookie.MaxAge = -1
				http.SetCookie(res, cookie)
			}
		} else {
			log.Println("Error reading cookie:", err)
		}
		login.LoginHandler(res, req)
	}
}
