package main

import (
	"backend/internal"
	"log"
	"net/http"
)

func main() {
	hub := internal.NewHub()
	go hub.Run()

	// Create the router
	router := internal.NewRouter(hub)

	// Start the server
	log.Println("HTTP server starting on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
