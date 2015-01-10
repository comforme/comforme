package search

import (
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func SearchHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}
	if req.Method == "POST" {
		// TODO uncomment when put to use
		//query := req.PostFormValue("query")
	}

	// TODO: Add template and compile it.
	common.ExecTemplate(nil, res, data)
}
