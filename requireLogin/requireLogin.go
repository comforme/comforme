package requireLogin

import (
	"log"
	"net/http"

	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
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
				handler(res, req)
				return
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
