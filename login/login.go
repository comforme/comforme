package login

import (
	"net/http"

	"github.com/comforme/comforme/common"
	"github.com/comforme/comforme/databaseActions"
)

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		email := req.PostFormValue("email")
		password := req.PostFormValue("password")
		sessionid, err := databaseActions.Login(email, password)
		if err != nil {

		}

		common.SetSessionCookie(req, sessionid)
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(nil, res, data)
}
