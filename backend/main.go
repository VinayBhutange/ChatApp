package main

import (
	"log"
	"net/http"
)

func main() {
	hub := newHub()
	go hub.run()

	// Serve static files from the "../frontend" directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Configure the WebSocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Start the server
	log.Println("HTTP server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
