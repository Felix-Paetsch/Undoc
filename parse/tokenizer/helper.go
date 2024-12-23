package tokenizer

import "strings"

func (t *Tokenizer) current() byte {
	if len(t.RemainingText) == 0 {
		return 0
	}
	return t.RemainingText[0]
}

func (t *Tokenizer) next() byte {
	if len(t.RemainingText) < 1 {
		return 0
	}
	return t.RemainingText[1]
}

func (t *Tokenizer) nextN(n int) byte {
	if len(t.RemainingText) < n {
		return 0
	}
	return t.RemainingText[n]
}

func (t *Tokenizer) consume() byte {
	if len(t.RemainingText) == 0 {
		return 0 // Return 0 if no characters are left
	}
	ch := t.RemainingText[0]              // Get the first character
	t.RemainingText = t.RemainingText[1:] // Remove the first character
	return ch
}

func (t *Tokenizer) consumeN(n int) string {
	if len(t.RemainingText) < n {
		n = len(t.RemainingText) // Consume only available characters
	}
	consumed := t.RemainingText[:n]       // Get the first n characters
	t.RemainingText = t.RemainingText[n:] // Remove them from the text
	return consumed
}

func (t *Tokenizer) addToken(tokenType TokenType, value string, src string) {
	token := Token{
		Type:  tokenType,
		Line:  t.CurrentLine,
		Value: value,
		Src:   src,
	}

	t.Tokens = append(t.Tokens, token)
	t.StateManager.processToken(token)
}

const Whitespace = " \t\r\n" // Define valid whitespace characters

func isWhitespace(b byte) bool {
	// Check if the byte exists in the Whitespace string
	return strings.ContainsRune(Whitespace, rune(b))
}
