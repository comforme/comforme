package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting server...")

	// Catchall
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "It works!")
	})

	// Start the server
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
