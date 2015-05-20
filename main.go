package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"github.com/comforme/comforme/ajax"
	"github.com/comforme/comforme/home"
	"github.com/comforme/comforme/logout"
	"github.com/comforme/comforme/pages"
	"github.com/comforme/comforme/requireLogin"
	"github.com/comforme/comforme/search"
	"github.com/comforme/comforme/settings"
	"github.com/comforme/comforme/static"
	"github.com/comforme/comforme/wizard"
)

func main() {
	log.Println("Starting server on port " + os.Getenv("PORT") + "...")

	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Current working directory:", dir)

	router := httprouter.New()

	router.GET(
		"/settings",
		requireLogin.RequireLogin(settings.SettingsHandler),
	)
	router.POST(
		"/settings",
		requireLogin.RequireLogin(settings.SettingsHandler),
	)

	router.GET(
		"/wizard",
		wizard.WizardHandler,
	)
	router.POST(
		"/wizard",
		wizard.WizardHandler,
	)

	router.GET(
		"/newPage",
		requireLogin.RequireLogin(pages.NewPageHandler),
	)
	router.POST(
		"/newPage",
		requireLogin.RequireLogin(pages.NewPageHandler),
	)

	router.GET(
		"/page/:category/:slug",
		requireLogin.RequireLogin(pages.PageHandler),
	)
	router.POST(
		"/page/:category/:slug",
		requireLogin.RequireLogin(pages.PageHandler),
	)

	router.POST(
		"/search",
		requireLogin.RequireLogin(search.SearchHandler),
	)

	router.GET(
		"/static/*filepath",
		static.StaticHandler,
	)

	router.POST(
		"/ajax/:action",
		requireLogin.AjaxRequireLogin(ajax.HandleAction),
	)

	router.GET(
		"/logout",
		logout.LogoutHandler,
	)

	router.GET(
		"/",
		requireLogin.RequireLogin(home.HomeHandler),
	)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
