package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize_containsSelectStatement_tokensCreated(t *testing.T) {
	// Arrange
	statement := "SELECT id FROM users"

	// Act
	result := Tokenize(statement)

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
	statement := "INSERT INTO users VALUES (1, \"TIM\")"

	// Act
	result := Tokenize(statement)

	// Assert
	assert.Len(t, result, 10)
	assertToken(t, result[0], INSERT, "insert")
	assertToken(t, result[1], INTO, "into")
	assertToken(t, result[2], IDENT, "users")
	assertToken(t, result[3], VALUES, "values")
	assertToken(t, result[4], LPAREN, "(")
	assertToken(t, result[5], NUMBER, "1")
	assertToken(t, result[6], COMMA, ",")
	assertToken(t, result[7], STRING, "TIM")
	assertToken(t, result[8], RPAREN, ")")
	assertToken(t, result[9], EOF, "")
}

func TestTokenize_containsSingleQuotedString_stringTokenCreated(t *testing.T) {
	// Arrange
	statement := "INSERT INTO users VALUES (1, 'Tim')"

	// Act
	result := Tokenize(statement)

	// Assert
	assert.Len(t, result, 10)
	assertToken(t, result[7], STRING, "Tim")
}

func TestTokenize_containsSymbol_symbolTokenCreated(t *testing.T) {
	// Arrange
	statement := "SELECT * FROM users"

	// Act
	result := Tokenize(statement)

	// Assert
	assertToken(t, result[1], SYMBOL, "*")
}

func TestTokenize_differentCasing_tokensCreated(t *testing.T) {
	// Arrange
	statement := "SeleCT Id FroM USerS"

	// Act
	result := Tokenize(statement)

	// Assert
	assertToken(t, result[0], SELECT, "select")
	assertToken(t, result[1], IDENT, "id")
	assertToken(t, result[2], FROM, "from")
	assertToken(t, result[3], IDENT, "users")
}

func assertToken(t assert.TestingT, token Token, expectedType TokenType, expectedValue string) {
	assert.Equal(t, expectedType, token.Type)
	assert.Equal(t, expectedValue, token.Value)
}
