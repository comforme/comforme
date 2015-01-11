package home

import (
	"net/http"

	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}

	// TODO: Add template and compile it.
	common.ExecTemplate(nil, res, data)
}
