package requireLogin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/comforme/comforme/ajax"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
	"github.com/comforme/comforme/settings"
)

func RequireLogin(handler func(http.ResponseWriter, *http.Request, string, string, string, int)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("sessionid")
		if err == nil {
			sessionid := cookie.Value
			email, username, userID, err := databaseActions.GetUserInfo(sessionid)
			if err == nil {
				log.Printf("User with email %s logged in.", email)
				isRequired, err := databaseActions.PasswordChangeRequired(sessionid)
				if err == nil {
					if isRequired {
						settings.SettingsHandler(res, req, sessionid, email, username, userID)
						return
					} else {
						handler(res, req, sessionid, email, username, userID)
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

func AjaxRequireLogin(handler func(http.ResponseWriter, *http.Request, string, string, string, int)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("sessionid")
		if err == nil {
			sessionid := cookie.Value
			email, username, userID, err := databaseActions.GetUserInfo(sessionid)
			if err == nil {
				log.Printf("User with email %s logged in.", email)
				handler(res, req, sessionid, email, username, userID)
				return
			}

			log.Println("Error checking email:", err)

			// Delete bad cookie
			cookie.MaxAge = -1
			http.SetCookie(res, cookie)
		} else {
			log.Println("Error reading cookie:", err)
		}

		// Return JSON Error
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintln(res, ajax.JSONLoginError)
	}
}
