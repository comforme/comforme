package logout

import (
	"log"
	"net/http"

	"github.com/comforme/comforme/databaseActions"
)

func LogoutHandler(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("sessionid")
	if err == nil {
		log.Println("Logging out sessionid:", cookie.Value)
		databaseActions.Logout(cookie.Value)
		if err != nil {
			log.Println("Logout session error:", err)
		}
		
		// Delete cookie
		cookie.MaxAge = -1
		http.SetCookie(res, cookie)
	} else {
		log.Println("Unable to logout, no cookie set:", err)
	}

	// Redirect to home page
	http.Redirect(res, req, "/", http.StatusFound)
}
