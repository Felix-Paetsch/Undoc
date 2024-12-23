package parser

import (
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

	// 6) Expect '{' (TokenTagStart)
	tagOpen := p.consume()
	if tagOpen.Type != tokenizer.TokenTagStart {
		return DocFile{}, p.ErrorOut("Expected '{' to start tag list", tagOpen.Line)
	}

	// 7) Parse at least one tag, separated by commas, until we hit TokenTagEnd
	tags, err := p.parseTags()
	if err != nil {
		return DocFile{}, err
	}
	doc.Tags = tags

	// 8) After '}', skip whitespace/newlines but enforce at least ONE newline
	foundNewLine := false
	for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
		if p.current().Type == tokenizer.TokenNL {
			foundNewLine = true
		}
		p.consume()
	}
	if !foundNewLine {
		return DocFile{}, p.ErrorOut("Expected newline after '}'", p.current().Line)
	}

	var content string
	for p.current().Type != tokenizer.TokenEmpty {
		content += p.consume().Src
	}

	doc.Content = content

	return doc, nil
}

func (p *Parser) parseTags() ([]string, error) {
	var tags []string
	foundAtLeastOne := false

	var last_tag_line = p.current().Line
	for {
		for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
			p.consume()
		}

		if p.current().Type == tokenizer.TokenTagEnd {
			closing := p.consume()
			if !foundAtLeastOne {
				return nil, p.ErrorOut("At least one tag required before '}'", closing.Line)
			}
			return tags, nil
		}

		if p.current().Type != tokenizer.TokenTag {
			return nil, p.ErrorOut("Expected Tag or '}'", p.current().Line)
		}

		tag := p.consume()
		tags = append(tags, tag.Value)
		foundAtLeastOne = true
		last_tag_line = tag.Line

		for p.current().Type == tokenizer.TokenWhitespace || p.current().Type == tokenizer.TokenNL {
			p.consume()
		}

		if p.current().Type == tokenizer.TokenTagSeperator {
			p.consume()
			continue
		}

		if p.current().Type == tokenizer.TokenTagEnd {
			continue
		}

		return nil, p.ErrorOut("Expected ',' or '}' after tag", last_tag_line)
	}
}
