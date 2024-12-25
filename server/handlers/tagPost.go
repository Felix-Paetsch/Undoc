package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"undoc/parse/parser"
	"undoc/search"
)

// filterTags ensures uniqueness, trims spaces, and converts tags to lowercase
func filterTags(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range input {
		v = strings.TrimSpace(strings.ToLower(v)) // Normalize tags
		if !seen[v] && v != "" {                  // Avoid duplicates and empty strings
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// TagPostHandler adds a new tag and updates search results
func TagPostHandler(w http.ResponseWriter, r *http.Request, docStore *search.SearchableDoc) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Log form parameters for debugging
	for k, v := range r.Form {
		fmt.Printf("%s: %v\n", k, v)
	}

	// Parse current tags
	var currentTags []string
	currentTagsJSON := r.FormValue("current_tags")
	if err := json.Unmarshal([]byte(currentTagsJSON), &currentTags); err != nil {
		http.Error(w, "Invalid tag format", http.StatusBadRequest)
		return
	}

	// Extract new tag and query
	newTag := strings.TrimSpace(strings.ToLower(r.FormValue("new_tag")))
	query := r.FormValue("current_query")

	// Append the new tag and ensure uniqueness
	tags := append(currentTags, newTag)
	tags = filterTags(tags)

	// Perform search with updated tags
	titleMatches, contentMatches := docStore.Search(query, tags)

	// Prepare data for rendering
	data := struct {
		Tags           []string
		Query          string
		TitleMatches   []parser.DocFile
		ContentMatches []parser.DocFile
	}{
		Tags:           tags,
		Query:          query,
		TitleMatches:   titleMatches,
		ContentMatches: contentMatches,
	}

	// Render updated tags and search results
	renderTemplate(w, data, "htmx_responses/update_tags.html", "htmx_responses/partials/query_actions.html")
}
