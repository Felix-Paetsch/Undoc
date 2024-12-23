package tokenizer

type TokenType string

const (
	TokenEOF           TokenType = "EOF"            // End of File
	TokenNL            TokenType = "NL"             // New Line
	TokenError         TokenType = "ERROR"          // Error
	TokenHashtag       TokenType = "HASHTAG"        // Hashtag at beginn of line for heading
	TokenSingleHashtag TokenType = "SINGLE_HASHTAG" // Single Hashtag, indication big heading
	TokenString        TokenType = "STRING"         // Generic String
	TokenWhitespace    TokenType = "WHITESPACE"     // Whitespace outside of text
	TokenTagStart      TokenType = "TAG_START"      // { for a tag start
	TokenTagEnd        TokenType = "TAG_END"        // } for a tag end
	TokenTagSeperator  TokenType = "TAG_SEPERATOR"  // , in tag thing
	TokenTag           TokenType = "TAG"
	TokenEmpty         TokenType = "EMPTY"
)

type Token struct {
	Type  TokenType
	Line  int
	Value string
	Src   string
}

type Tokenizer struct {
	Tokens        []Token
	CurrentLine   int
	RemainingText string
	StateManager  *TokenizerStateManager
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		Tokens:        []Token{},
		CurrentLine:   0,
		RemainingText: input,
		StateManager:  NewTokenizerStateManager(),
	}
}

func (t *Tokenizer) Tokenize() {
	if len(t.RemainingText) > 0 {
		t.CurrentLine = 1
	}

	for len(t.RemainingText) > 0 && !t.StateManager.isError() {
		current := t.current() // Declare before switch
		if current == '\n' {
			t.tokenizeNewLine()
			continue
		}

		if isWhitespace(current) {
			t.tokenizeWhitespace()
			continue
		}

		if t.StateManager.afterTagEnd() {
			t.tokenizeError("Expected whitespace or new line.")
			continue
		}

		if current == '#' && (t.StateManager.atLineStart()) {
			t.tokenizeHashtags()
			continue
		}

		if t.StateManager.afterSingleHashtag() && current == '{' {
			t.tokenizeTagStart()
			continue
		}

		if t.StateManager.parsingTags() {
			if current == ',' {
				t.tokenizeTagSeperator()
				continue
			}
			if current == '}' {
				t.tokenizeTagEnd()
				continue
			}
			if current == '{' {
				t.tokenizeError("Unexpected '{', expected '}'")
				continue
			}

			if current == 0 {
				t.tokenizeError("Unexpected EOF, expected '}'")
				continue
			}

			t.tokenizeTag()
			continue
		}

		t.tokenizeText()
	}

	t.addToken(TokenEOF, "", "")
}

func (t *Tokenizer) tokenizeWhitespace() {
	var value = ""
	for len(t.RemainingText) > 0 && isWhitespace(t.current()) && t.current() != '\n' {
		value += string(t.consume())
	}

	t.addToken(TokenWhitespace, "", value)
}

func (t *Tokenizer) tokenizeText() {
	var value string
	for len(t.RemainingText) > 0 && t.current() != '\n' {
		value += string(t.consume())
	}
	t.addToken(TokenString, value, value)
}

func (t *Tokenizer) tokenizeHashtags() {
	var value string
	for len(t.RemainingText) > 0 && t.current() == '#' {
		value += string(t.consume())
	}

	if (t.current() != 0 && !isWhitespace(t.current())) || len(value) > 6 {
		t.addToken(TokenString, value, value)
		return
	}

	if len(value) == 1 {
		t.addToken(TokenSingleHashtag, "", "#")
	} else {
		t.addToken(TokenHashtag, "", value)
	}
}

func (t *Tokenizer) tokenizeNewLine() {
	if t.consume() == '\n' {
		t.addToken(TokenNL, "", "\n")
		t.CurrentLine++
	} else {
		panic("Unreachable state: Expected '\n'")
	}
}

func (t *Tokenizer) tokenizeTagStart() {
	if t.consume() == '{' {
		t.addToken(TokenTagStart, "", "{")
	} else {
		panic("Unreachable state: Expected '{'")
	}
}

func (t *Tokenizer) tokenizeTagEnd() {
	if t.consume() == '}' {
		t.addToken(TokenTagEnd, "", "}")
	} else {
		panic("Unreachable state: Expected '}'")
	}
}

func (t *Tokenizer) tokenizeTagSeperator() {
	if t.consume() == ',' {
		t.addToken(TokenTagSeperator, "", ",")
	} else {
		panic("Unreachable state: Expected ','")
	}
}

func (t *Tokenizer) tokenizeError(val string) {
	t.addToken(TokenError, val, "")
}

func (t *Tokenizer) tokenizeTag() {
	// Find where the tag would end
	endIndex := findTagEndIndex(t.RemainingText)
	if endIndex == 0 {
		panic("Unreachable state: Empty tag value")
	}

	value := t.consumeN(endIndex)
	t.addToken(TokenTag, value, value)
}

// findTagEndIndex returns the position of the first delimiter
// (',' , '{' , '}' , '\n', or end-of-input) in the string.
// If no delimiter is found, it returns len(s).
func findTagEndIndex(s string) int {
	lastNonWhitespace := -1

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 0, ',', '{', '}', '\n':
			return lastNonWhitespace + 1
		default:
			if !isWhitespace(s[i]) {
				lastNonWhitespace = i
			}
		}
	}

	return lastNonWhitespace + 1
}
