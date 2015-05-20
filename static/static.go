package static

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func StaticHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("cache-control", "public, no-transform")
	path := "./staticFiles/" + ps.ByName("filepath")
	log.Println("Serving file:", path)
	http.ServeFile(res, req, path)
}
