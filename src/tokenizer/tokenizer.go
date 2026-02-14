package tokenizer

import (
	"strings"
)

type TokenType int

const (
	SELECT TokenType = iota
	INSERT
	IDENT
	FROM
	SYMBOL
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

var tokenRegistry = map[string]Token{
	"select": {Type: SELECT, Value: "select"},
	"insert": {Type: INSERT, Value: "insert"},
	"from":   {Type: FROM, Value: "from"},
	"*":      {Type: SYMBOL, Value: "*"},
}

func Tokenize(statement string) []Token {
	words := strings.Fields(statement)

	tokens := make([]Token, len(words))

	for i, word := range words {
		tokens[i] = toToken(word)
	}

	tokens = append(tokens, Token{Type: EOF})
	return tokens
}

func toToken(word string) Token {
	lowerWord := strings.ToLower(word)

	if token, exists := tokenRegistry[lowerWord]; exists {
		return token
	}

	// If it's not registered, assume it's a generic IDENTIFIER (table name, column, ...)
	return Token{Type: IDENT, Value: lowerWord}
}
