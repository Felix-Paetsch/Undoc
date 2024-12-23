package parse

import (
	"undoc/parse/parser"
	"undoc/parse/tokenizer"
)

func ParseDocFile(filePath, content string) (parser.DocFile, error) {
	// Tokenize
	t := tokenizer.NewTokenizer(content)
	t.Tokenize()

	// Create parser
	p := parser.NewParser(filePath, content, t.Tokens)
	return p.Parse()
}
