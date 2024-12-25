package handlers

import (
	"encoding/json"
	"net/http"
	"undoc/parse/parser"
	"undoc/search"
)

// SearchHandler processes search requests
func SearchHandler(w http.ResponseWriter, r *http.Request, docStore *search.SearchableDoc) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract tags
	var currentTags []string
	currentTagsJSON := r.FormValue("current_tags")
	if err := json.Unmarshal([]byte(currentTagsJSON), &currentTags); err != nil {
		http.Error(w, "Invalid tag format", http.StatusBadRequest)
		return
	}

	// Extract query
	query := r.FormValue("current_query")

	// Perform search
	titleMatches, contentMatches := docStore.Search(query, currentTags)

	// Prepare data for rendering
	data := struct {
		Tags           []string
		Query          string
		TitleMatches   []parser.DocFile
		ContentMatches []parser.DocFile
	}{
		Tags:           currentTags,
		Query:          query,
		TitleMatches:   titleMatches,
		ContentMatches: contentMatches,
	}

	// Render updated template with search results
	renderTemplate(w, data, "htmx_responses/update_query.html", "htmx_responses/partials/query_actions.html")
}
