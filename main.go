package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	
	"github.com/go-zoo/bone"
)

func main() {
	log.Println("Starting server...")
	mux := bone.New()

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
