package requireLogin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/ajax"
	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
	"github.com/comforme/comforme/settings"
)

func RequireLogin(handler func(http.ResponseWriter, *http.Request, httprouter.Params, common.UserInfo)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		cookie, err := req.Cookie("sessionid")
		if err == nil {
			sessionid := cookie.Value
			userInfo, err := databaseActions.GetUserInfo(sessionid)
			if err == nil {
				log.Printf("User with email %s logged in.", userInfo.Email)
				isRequired, err := databaseActions.PasswordChangeRequired(sessionid)
				if err == nil {
					if isRequired {
						settings.SettingsHandler(res, req, ps, userInfo)
						return
					} else {
						handler(res, req, ps, userInfo)
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

func AjaxRequireLogin(handler func(http.ResponseWriter, *http.Request, httprouter.Params, common.UserInfo)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		cookie, err := req.Cookie("sessionid")
		if err == nil {
			sessionid := cookie.Value
			userInfo, err := databaseActions.GetUserInfo(sessionid)
			if err == nil {
				log.Printf("User with email %s logged in.", userInfo.Email)
				handler(res, req, ps, userInfo)
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
