package static

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func StaticHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	http.ServeFile(res, req, "../staticFiles/"+ps.ByName("filepath"))
}
