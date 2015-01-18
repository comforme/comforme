package static

import (
	"fmt"
	"net/http"
)

func Style(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/css; charset=utf-8")

	fmt.Fprintln(res, styleCSS)
}

func SettingsJS(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/javascript; charset=utf-8")

	fmt.Fprintln(res, settingsJS)
}
