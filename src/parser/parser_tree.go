package parser

import (
	"learning/db/tokenizer"
)

type ParserTree struct {
	tokens   []tokenizer.Token
	position int
}

func (p *ParserTree) peek() tokenizer.Token {
	if p.position >= len(p.tokens) {
		return tokenizer.Token{Type: tokenizer.EOF}
	}
	return p.tokens[p.position]
}

func (p *ParserTree) consume() tokenizer.Token {
	token := p.peek()
	p.position++
	return token
}
