package parser

import "undoc/parse/tokenizer"

func (p *Parser) current() tokenizer.Token {
	if len(p.tokens) == 0 {
		return tokenizer.Token{
			Type: tokenizer.TokenEmpty,
			Line: -1, Value: "", Src: "",
		}
	}
	return p.tokens[0]
}

func (p *Parser) next() tokenizer.Token {
	if len(p.tokens) < 1 {
		return tokenizer.Token{
			Type: tokenizer.TokenEmpty,
			Line: -1, Value: "", Src: "",
		}
	}
	return p.tokens[1]
}

func (p *Parser) nextN(n int) tokenizer.Token {
	if len(p.tokens) < n {
		return tokenizer.Token{
			Type: tokenizer.TokenEmpty,
			Line: -1, Value: "", Src: "",
		}
	}
	return p.tokens[n]
}

func (p *Parser) consume() tokenizer.Token {
	if len(p.tokens) == 0 {
		return tokenizer.Token{
			Type: tokenizer.TokenEmpty,
			Line: -1, Value: "", Src: "",
		}
	}
	ch := p.tokens[0]       // Get the first character
	p.tokens = p.tokens[1:] // Remove the first character
	return ch
}

func (p *Parser) consumeN(n int) []tokenizer.Token {
	if len(p.tokens) < n {
		n = len(p.tokens) // Consume only available characters
	}
	consumed := p.tokens[:n] // Get the first n characters
	p.tokens = p.tokens[n:]  // Remove them from the text
	return consumed
}
