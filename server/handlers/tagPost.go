package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func filterTags(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range input {
		v = strings.TrimSpace(strings.ToLower(v)) // Fixes trim() and toLowerCase()
		if !seen[v] && v != "" {                  // Checks for duplicates and ignores empty strings
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func TagPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Log all form parameters
	for k, v := range r.Form {
		fmt.Printf("%s: %v\n", k, v)
	}

	// Parse current_tags as JSON
	var currentTags []string
	currentTagsJSON := r.FormValue("current_tags")
	if err := json.Unmarshal([]byte(currentTagsJSON), &currentTags); err != nil {
		http.Error(w, "Invalid tag format", http.StatusBadRequest)
		return
	}

	// Extract new tag and query
	newTag := strings.ToLower(r.FormValue("new_tag"))
	query := r.FormValue("current_query")

	// Append the new tag and ensure uniqueness
	tags := append(currentTags, newTag)
	tags = filterTags(tags)

	// Prepare data for rendering
	data := struct {
		Tags  []string
		Query string
	}{
		Tags:  tags,
		Query: query,
	}

	// Render updated tags
	renderTemplate(w, data, "htmx_responses/update_tags.html", "htmx_responses/partials/query_actions.html")
}
