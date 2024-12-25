package server

import (
	"fmt"
	"net/http"

	"undoc/search"
	"undoc/server/handlers"
)

type Server struct {
	DocStore *search.SearchableDoc
}

func (s *Server) Start() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/tagPost", func(w http.ResponseWriter, r *http.Request) {
		handlers.TagPostHandler(w, r, s.DocStore)
	})
	http.HandleFunc("/tagDelete", func(w http.ResponseWriter, r *http.Request) {
		handlers.TagDeleteHandler(w, r, s.DocStore)
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.SearchHandler(w, r, s.DocStore)
	})
	http.HandleFunc("/doc", func(w http.ResponseWriter, r *http.Request) {
		handlers.DocHandler(w, r, s.DocStore)
	})

	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
