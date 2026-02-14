package main

import (
	"strings"
)

type TokenType int

const (
	SELECT TokenType = iota
	IDENT
	FROM
	SYMBOL
)

type Token struct {
	Type  TokenType
	Value string
}

var tokenRegistry = map[string]Token{
	"select": {Type: SELECT, Value: "select"},
	"from":   {Type: FROM, Value: "from"},
	"*":      {Type: SYMBOL, Value: "*"},
}

func Tokenize(statement string) []Token {
	words := strings.Fields(statement)

	tokens := make([]Token, len(words))

	for i, word := range words {
		tokens[i] = toToken(word)
	}

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
