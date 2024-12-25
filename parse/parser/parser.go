package parser

import (
	"strings"
	"undoc/parse/tokenizer"
)

type DocFile struct {
	FilePath string
	SrcText  string

	Title   string
	Tags    []string
	Content string
}

type ParserError struct {
	Message  string
	FilePath string
	SrcText  string
	Line     int
}

type Parser struct {
	filePath string
	srcText  string
	tokens   []tokenizer.Token
	pos      int // current index into tokens
}

func NewParser(filePath, content string, tokens []tokenizer.Token) *Parser {
	return &Parser{
		filePath: filePath,
		srcText:  content,
		tokens:   tokens,
		pos:      0,
	}
}

func (p *Parser) Parse() (DocFile, error) {
	// 1) Check for token errors
	for _, token := range p.tokens {
		if token.Type == tokenizer.TokenError {
			return DocFile{}, p.ErrorOut(
				token.Value,
				token.Line,
			)
		}
	}

	// Prepare the DocFile we will populate
	doc := DocFile{
		FilePath: p.filePath,
		SrcText:  p.srcText,
	}

	// 2) Expect a single hashtag
	hashtag := p.consume()
	if hashtag.Type != tokenizer.TokenSingleHashtag {
		return DocFile{}, p.ErrorOut("Expected single hashtag (#)", hashtag.Line)
	}

	// 3) Next should be whitespace
	whitespace := p.consume()
	if whitespace.Type != tokenizer.TokenWhitespace {
		return DocFile{}, p.ErrorOut("Expected whitespace after '#'", whitespace.Line)
	}

	// 4) Next token is the Title (must be STRING)
	titleToken := p.consume()
	if titleToken.Type != tokenizer.TokenString {
		return DocFile{}, p.ErrorOut("Expected title (string)", titleToken.Line)
	}
	doc.Title = titleToken.Value

	// 5) Skip any extra whitespace/newlines
	for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
		p.consume()
	}

	// 6) Handle optional tag section
	if p.current().Type == tokenizer.TokenTagStart { // Tag section begins
		p.consume() // Consume '{'

		// Check if immediately closed
		if p.current().Type == tokenizer.TokenTagEnd {
			p.consume()           // Consume '}'
			doc.Tags = []string{} // Empty tags
		} else {
			// Parse tags
			tags, err := p.parseTags()
			if err != nil {
				return DocFile{}, err
			}
			doc.Tags = tags
		}

		// Require newline after tags
		foundNewLine := false
		for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL || p.current().Type == tokenizer.TokenEOF {
			if p.current().Type == tokenizer.TokenNL || p.current().Type == tokenizer.TokenEOF {
				foundNewLine = true
			}
			p.consume()
		}
		if !foundNewLine {
			return DocFile{}, p.ErrorOut("Expected newline after '}'", p.current().Line)
		}
	} else {
		// No tag boundaries; treat as zero tags
		doc.Tags = []string{}
	}

	// 7) Parse remaining content
	var content string
	for p.current().Type != tokenizer.TokenEmpty {
		content += p.consume().Src
	}

	doc.Content = content

	return doc, nil
}

func (p *Parser) parseTags() ([]string, error) {
	var tags []string
	seen := make(map[string]bool) // Track unique tags

	var last_tag_line = p.current().Line
	for {
		// Skip whitespace and newlines
		for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
			p.consume()
		}

		// Check for end of tag list
		if p.current().Type == tokenizer.TokenTagEnd {
			p.consume()
			return tags, nil
		}

		// Validate the tag type
		if p.current().Type != tokenizer.TokenTag {
			return nil, p.ErrorOut("Expected Tag or '}'", p.current().Line)
		}

		// Process the tag
		tag := p.consume()
		normalizedTag := strings.TrimSpace(strings.ToLower(tag.Value)) // Normalize tag

		// Add tag only if it's unique
		if !seen[normalizedTag] {
			tags = append(tags, tag.Value) // Preserve original formatting
			seen[normalizedTag] = true     // Mark as seen
		}
		last_tag_line = tag.Line

		// Skip trailing whitespace/newlines
		for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
			p.consume()
		}

		// Handle tag separators
		if p.current().Type == tokenizer.TokenTagSeperator {
			p.consume()
			continue
		}

		// Handle tag end
		if p.current().Type == tokenizer.TokenTagEnd {
			continue
		}

		// Error if unexpected token appears
		return nil, p.ErrorOut("Expected ',' or '}' after tag", last_tag_line)
	}
}
