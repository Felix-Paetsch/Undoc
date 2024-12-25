package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings" // Import strings package
)

func renderTemplate(w http.ResponseWriter, data interface{}, tmpl ...string) {
	if len(tmpl) == 0 {
		http.Error(w, "No templates specified", http.StatusInternalServerError)
		return
	}

	// Build full file paths
	var paths []string
	for _, file := range tmpl {
		paths = append(paths, filepath.Join("server/pages", file))
	}

	// Helper functions
	toJSON := func(data interface{}) string {
		bytes, err := json.Marshal(data)
		if err != nil {
			return "[]"
		}
		return string(bytes)
	}

	replace := func(input, old, new string) string {
		return strings.ReplaceAll(input, old, new)
	}

	trimBackticks := func(input string) string {
		// Trim only if the whole title is not wrapped in backticks
		if strings.HasPrefix(input, "`") && strings.HasSuffix(input, "`") && len(strings.Trim(input, "`")) > 0 {
			return strings.Trim(input, "`") // Trim surrounding backticks
		}
		return input // Return unchanged if fully wrapped in backticks
	}

	// Register template functions
	tmplFunc := template.New("").Funcs(template.FuncMap{
		"json":          toJSON,
		"replace":       replace,
		"trimBackticks": trimBackticks, // Register the function
	})

	// Parse and render templates
	templates, err := tmplFunc.ParseFiles(paths...)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "full_layout", data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
