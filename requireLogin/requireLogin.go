package requireLogin

import(
	"net/http"
	"log"
	
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
)

func RequireLogin(handler func (http.ResponseWriter, *http.Request)) (func (http.ResponseWriter, *http.Request)) {
	return func (res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("sessionid")
		if err != nil {
			email, err := databaseActions.GetEmail(cookie.Value)
			if err != nil {
				log.Printf("User with email %s logged in.", email)
				handler(res, req)
				return
			} else {
				log.Println("Error checking email:", err)
			}
		} else {
			log.Println("Error reading cookie:", err)
		}
		login.LoginHandler(res, req)
	}
}
