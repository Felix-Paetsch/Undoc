package handlers

import (
	"encoding/json"
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse current_tags as JSON
	var currentTags []string
	currentTagsJSON := r.FormValue("current_tags")
	if err := json.Unmarshal([]byte(currentTagsJSON), &currentTags); err != nil {
		http.Error(w, "Invalid tag format", http.StatusBadRequest)
		return
	}
	query := r.FormValue("current_query")

	// Prepare data for rendering
	data := struct {
		Tags  []string
		Query string
	}{
		Tags:  currentTags,
		Query: query,
	}

	// Render updated tags
	renderTemplate(w, data, "htmx_responses/update_query.html", "htmx_responses/partials/query_actions.html")
}
