package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-zoo/bone"

	"github.com/comforme/comforme/home"
	"github.com/comforme/comforme/login"
	"github.com/comforme/comforme/logout"
	"github.com/comforme/comforme/pages"
	"github.com/comforme/comforme/settings"
	"github.com/comforme/comforme/requireLogin"
	"github.com/comforme/comforme/search"
	"github.com/comforme/comforme/static"
)

func main() {
	log.Println("Starting server on port " + os.Getenv("PORT") + "...")
	mux := bone.New()

	mux.Handle(
		"/login",
		http.HandlerFunc(
			login.LoginHandler,
		),
	)

	mux.Handle(
		"/settings",
		http.HandlerFunc(
			requireLogin.RequireLogin(settings.SettingsHandler),
		),
	)

	mux.Handle(
		"/pages",
		http.HandlerFunc(
			requireLogin.RequireLogin(pages.PagesHandler),
		),
	)

	mux.Handle(
		"/search",
		http.HandlerFunc(
			requireLogin.RequireLogin(search.SearchHandler),
		),
	)

	mux.Handle(
		"/style_css",
		http.HandlerFunc(
			static.Style,
		),
	)

	mux.Handle(
		"/logout",
		http.HandlerFunc(
			logout.LogoutHandler,
		),
	)

	mux.Handle(
		"/",
		http.HandlerFunc(
			requireLogin.RequireLogin(home.HomeHandler),
		),
	)

	// Start the server
	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		panic(err)
	}
}
