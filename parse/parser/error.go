package parser

import (
	"fmt"
	"strings"
)

func (e ParserError) Error() string {
	return e.Message
}

func (p *Parser) ErrorOut(msg string, line int) ParserError {
	// 1) If line <= 0, treat as line = 1
	if line <= 0 {
		line = 1
	}

	// 2) If the source text is empty, just return the basic message.
	if len(p.srcText) == 0 {
		return ParserError{
			Message:  fmt.Sprintf("%s\nError at line 0: (No text in source)", p.filePath),
			FilePath: p.filePath,
			SrcText:  p.srcText,
			Line:     line,
		}
	}

	// 3) Split source text into lines
	lines := strings.Split(p.srcText, "\n")

	// If there's only one empty line, handle separately.
	if len(lines) == 1 && lines[0] == "" {
		return ParserError{
			Message:  fmt.Sprintf("%s\nError at line %d: %s\n\n(No text in source)", p.filePath, line, msg),
			FilePath: p.filePath,
			SrcText:  p.srcText,
			Line:     line,
		}
	}

	// Convert from 1-based line to 0-based index
	lineIndex := line - 1
	if lineIndex >= len(lines) {
		// If line is out of bounds, clamp it to the last line.
		lineIndex = len(lines) - 1
	}

	// 4) Calculate snippet range: 2 lines above, 2 lines below
	start := lineIndex - 2
	if start < 0 {
		start = 0
	}
	end := lineIndex + 3

	if end >= len(lines) {
		end = len(lines) - 1
	}

	// 5) Build the snippet
	var snippetLines []string
	for i := start; i < end; i++ {
		// Each line is labeled with "lineNum | "
		lineNum := i + 1 // convert back to 1-based for display
		snippetLines = append(snippetLines,
			fmt.Sprintf("%d | %s", lineNum, lines[i]),
		)
	}

	// 6) Combine everything into an error message
	errorMsg := fmt.Sprintf("%s\nError at line %d: %s\n\n", p.filePath, line, msg) +
		strings.Join(snippetLines, "\n")

	return ParserError{
		Message:  errorMsg,
		FilePath: p.filePath,
		SrcText:  p.srcText,
		Line:     line,
	}
}
