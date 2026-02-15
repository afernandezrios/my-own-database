package tokenizer

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTokenize_containsSelectStatement_tokensCreated(t *testing.T) {
	// Arrange
	statement:= "SELECT id FROM users"

	// Act
	result:= Tokenize(statement)

	// Assert
	assert.Len(t, result, 5)
	assertToken(t, result[0], SELECT, "select")
	assertToken(t, result[1], IDENT, "id")
	assertToken(t, result[2], FROM, "from")
	assertToken(t, result[3], IDENT, "users")
	assertToken(t, result[4], EOF, "")
}

func TestTokenize_containsInsertTokenWithNumbersAndStrings_tokensCreated(t *testing.T) {
	// Arrange
	statement:= "INSERT INTO users VALUES (1, \"TIM\")"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[0], INSERT, "insert")
	assertToken(t, result[0], INTO, "into")
	assertToken(t, result[0], IDENT, "users")
	assertToken(t, result[0], VALUES, "values")
	assertToken(t, result[0], LPAREN, "1")
	assertToken(t, result[0], NUMBER, "1")
	assertToken(t, result[0], STRING, "TIM")
	assertToken(t, result[0], RPAREN, "TIM")
}

func TestTokenize_containsSymbol_symbolTokenCreated(t *testing.T) {
	// Arrange
	statement:= "SELECT * FROM users"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[1], SYMBOL, "*")
}

func TestTokenize_differentCasing_tokensCreated(t *testing.T) {
	// Arrange
	statement:= "SeleCT Id FroM USerS"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[0], SELECT, "select")
	assertToken(t, result[1], IDENT, "id")
	assertToken(t, result[2], FROM, "from")
	assertToken(t, result[3], IDENT, "users")
}

func assertToken(t assert.TestingT, token Token, expectedType TokenType, expectedValue string ) {
	assert.Equal(t, expectedType, token.Type)
	assert.Equal(t, expectedValue, token.Value)
}