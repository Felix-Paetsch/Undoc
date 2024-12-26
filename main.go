package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"undoc/parse"
	"undoc/search"
	"undoc/server"
)

func main() {
	docsFolder := "Docs" // Default folder
	if len(os.Args) > 1 {
		docsFolder = os.Args[1] // Use provided folder if available
	}

	info, err := os.Stat(docsFolder)
	if os.IsNotExist(err) || !info.IsDir() {
		fmt.Printf("Folder '%s' not found or is not a directory.\n", docsFolder)
		os.Exit(1)
	}

	// Initialize search storage
	docStore := search.NewSearchableDoc()

	// Walk through the given folder
	err = filepath.Walk(docsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() &&
			strings.HasSuffix(path, ".md") &&
			!strings.HasSuffix(path, ".undoc.md") {

			content, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Printf("Failed to read file '%s'\n", path)
				os.Exit(1)
			}

			doc, parseErr := parse.ParseDocFile(path, string(content))
			if parseErr != nil {
				log.Fatalf("\n\n%s\n", parseErr.Error())
			}

			// Add document to searchable store
			docStore.AddDoc(doc)
			fmt.Printf("Parsed doc for '%s': %+v\n", path, doc.Title)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the folder: %v\n", err)
	}

	// Start server with the docStore
	srv := &server.Server{DocStore: docStore}
	srv.Start()
}
