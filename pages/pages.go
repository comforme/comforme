package pages

import (
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func PagesHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
        // TODO uncomment when put to use
		//title := req.PostFormValue("title")
		//description := req.PostFormValue("description")
		//address := req.PostFormValue("address")
        //categories := req.PostFormValue("categories")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(nil, res, data)
}
