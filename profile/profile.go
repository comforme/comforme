package profile

import (
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func ProfileHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
        // TODO uncomment when put to use
		//username := req.PostFormValue("username")
		//password := req.PostFormValue("password")
		//newPassword := req.PostFormValue("newPassword")
		//newPasswordConfirmation := req.PostFormValue("newPasswordConfirmation")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(nil, res, data)
}
