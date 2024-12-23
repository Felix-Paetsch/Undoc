package server

import (
	"fmt"
	"net/http"

	"undoc/server/handlers"
)

// Start initializes and runs the HTTP server
func Start() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	// Routes
	http.HandleFunc("/", handlers.IndexHandler)

	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
