package tokenizer

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTokenize_containsSelectToken_selectTokenCreated(t *testing.T) {
	// Arrange
	statement:= "SELECT id FROM users"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[0], SELECT, "select")
}

func TestTokenize_containsInsertToken_insertTokenCreated(t *testing.T) {
	// Arrange
	statement:= "INSERT INTO users VALUES (1, 'TIM')"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[0], INSERT, "insert")
}

func TestTokenize_containsFromToken_fromTokenCreated(t *testing.T) {
	// Arrange
	statement:= "SELECT id FROM users"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[2], FROM, "from")
}

func TestTokenize_containsIdentifiers_identifiersTokenCreated(t *testing.T) {
	// Arrange
	statement:= "SELECT id FROM users"

	// Act
	result:= Tokenize(statement)

	// Assert
	assertToken(t, result[1], IDENT, "id")
	assertToken(t, result[3], IDENT, "users")
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