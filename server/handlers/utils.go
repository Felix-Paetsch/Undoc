package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
)

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

	// Define a helper function to convert data to JSON
	toJSON := func(data interface{}) string {
		bytes, err := json.Marshal(data)
		if err != nil {
			return "[]" // Fallback to empty array
		}
		return string(bytes)
	}

	// Create a new template and register the "json" function
	tmplFunc := template.New("").Funcs(template.FuncMap{
		"json": toJSON,
	})

	// Parse all the template files
	templates, err := tmplFunc.ParseFiles(paths...)
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
