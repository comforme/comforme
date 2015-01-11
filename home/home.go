package home

import (
	"net/http"
        "html/template"
	"github.com/comforme/comforme/common"
	// "github.com/comforme/comforme/databaseActions"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	var data map[string]interface{}

	// TODO: Add template and compile it.
	tmpl, _ := template.New("test").ParseFiles("/templates/templates.go")
	common.ExecTemplate(tmpl, res, data)
}
