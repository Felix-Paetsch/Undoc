package handlers

import (
	"net/http"
	"path/filepath"
	"undoc/parse/parser"
	"undoc/search"
)

// DocHandler serves the "doc.html" template with the requested DocFile.
func DocHandler(w http.ResponseWriter, r *http.Request, docStore *search.SearchableDoc) {
	// Get file path from query parameters
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Find the requested document
	var doc *parser.DocFile
	for _, d := range docStore.Docs {
		if filepath.Clean(d.FilePath) == filepath.Clean(filePath) {
			doc = &d
			break
		}
	}

	// Handle document not found
	if doc == nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Render template using the utility function
	renderTemplate(w, doc, "doc.html")
}
