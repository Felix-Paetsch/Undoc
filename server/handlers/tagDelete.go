package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func TagDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	// Extract the tag to delete and query
	tagToDelete := strings.TrimSpace(strings.ToLower(r.FormValue("delete_tag")))
	query := r.FormValue("current_query")

	// Filter out the tag to delete
	var filteredTags []string
	for _, tag := range currentTags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != tagToDelete {
			filteredTags = append(filteredTags, tag)
		}
	}

	// Prepare data for rendering
	data := struct {
		Tags  []string
		Query string
	}{
		Tags:  filteredTags,
		Query: query,
	}

	// Render updated tags
	renderTemplate(w, data, "htmx_responses/update_tags.html", "htmx_responses/partials/query_actions.html")
}
