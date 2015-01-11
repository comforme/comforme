package requireLogin

import(
	"net/http"
	"github.com/comforme/comforme/databaseActions"
	"github.com/comforme/comforme/login"
)

func RequireLogin(handler func (http.ResponseWriter, *http.Request)) (func (http.ResponseWriter, *http.Request)) {
	return func (res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("sessionid")
		if err != nil {
			_, err := databaseActions.GetEmail(cookie.Value)
			if err != nil {
				handler(res, req)
				return
			}
		}
		login.LoginHandler(res, req)
	}
}
