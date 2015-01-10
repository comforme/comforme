package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-zoo/bone"
	
	"github.com/comforme/comforme/login"
	"github.com/comforme/comforme/profile"
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
		"/profile",
		http.HandlerFunc(
			profile.ProfileHandler,
		),
	)

	mux.Handle(
		"/pages",
		http.HandlerFunc(
			profile.ProfileHandler,
		),
	)

	mux.Handle(
		"/",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "It works!")
			},
		),
	)

	// Start the server
	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		panic(err)
	}
}
