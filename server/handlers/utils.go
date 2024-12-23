package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// renderTemplate takes any number of template filenames (relative to server/pages).
// It loads them all and executes the "full_layout" template.
func renderTemplate(w http.ResponseWriter, data interface{}, tmpl ...string) {
	// If no templates were provided, return an error
	if len(tmpl) == 0 {
		http.Error(w, "No templates specified", http.StatusInternalServerError)
		return
	}

	// Build a list of full file paths
	var paths []string
	for _, file := range tmpl {
		paths = append(paths, filepath.Join("server/pages", file))
	}

	// Parse all the template files
	templates, err := template.ParseFiles(paths...)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the named template "full_layout"
	err = templates.ExecuteTemplate(w, "full_layout", data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
