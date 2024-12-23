package handlers

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title   string
		Content string
	}{
		Title:   "Undoc",
		Content: "Welcome to our website!",
	}

	renderTemplate(w, data, "index.html")
}
