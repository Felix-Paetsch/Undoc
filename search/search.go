package search

import (
	"strings"
	"sync"
	"undoc/parse/parser"
)

type SearchableDoc struct {
	Docs []parser.DocFile
	mu   sync.RWMutex
}

func NewSearchableDoc() *SearchableDoc {
	return &SearchableDoc{
		Docs: []parser.DocFile{},
	}
}

func (s *SearchableDoc) AddDoc(doc parser.DocFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Docs = append(s.Docs, doc)
}

func (s *SearchableDoc) Search(query string, tags []string) ([]parser.DocFile, []parser.DocFile) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.TrimSpace(strings.ToLower(query))
	for i := range tags {
		tags[i] = strings.TrimSpace(strings.ToLower(tags[i]))
	}

	var titleMatches []parser.DocFile
	var contentMatches []parser.DocFile

	for _, doc := range s.Docs {
		if matchesTags(doc.Tags, tags) {
			if strings.Contains(strings.ToLower(doc.Title), query) {
				titleMatches = append(titleMatches, doc)
			} else if strings.Contains(strings.ToLower(doc.Content), query) {
				contentMatches = append(contentMatches, doc)
			}
		}
	}

	return titleMatches, contentMatches
}

func matchesTags(docTags, searchTags []string) bool {
	if len(searchTags) == 0 {
		return true
	}
	docTagSet := make(map[string]struct{})
	for _, tag := range docTags {
		docTagSet[strings.TrimSpace(strings.ToLower(tag))] = struct{}{}
	}

	for _, tag := range searchTags {
		if _, exists := docTagSet[tag]; !exists {
			return false
		}
	}

	return true
}
