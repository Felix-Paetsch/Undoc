package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"undoc/parse"
	"undoc/server"
)

func main() {
	docsFolder := "./Docs"

	info, err := os.Stat(docsFolder)
	if os.IsNotExist(err) || !info.IsDir() {
		fmt.Println("Folder './Docs' not found or is not a directory.")
		os.Exit(1)
	}

	filepath.Walk(docsFolder, func(path string, info os.FileInfo, err error) error {
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

			fmt.Printf("Parsed doc for '%s': %+v\n", path, doc.Title)
		}
		return nil
	})

	server.Start()
}
