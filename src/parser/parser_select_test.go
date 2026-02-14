package parser

import (
	"learning/db/tokenizer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSelect_doesNotStartWithSelectToken_error(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.INSERT, "insert")
	t2 := newToken(tokenizer.IDENT, "id")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, result)
	assert.EqualError(t, err, "This is not a SELECT command")
}

func TestParseSelect_missingSymbolOrIdentifierAfterSelectToken_error(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.FROM, "from")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, result)
	assert.EqualError(t, err, "SELECT must contain at least one INDETIFIER or SYMBOL")
}

func TestParseSelect_missingFromTokenAfterSelectFields_error(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.IDENT, "id")
	t3 := newToken(tokenizer.SELECT, "select")
	t4 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, result)
	assert.EqualError(t, err, "SELECT command requires FROM clause")
}

func TestParseSelect_missingTableNameAfterFromToken_createsValidStatement(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.IDENT, "id")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.FROM, "from")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, result)
	assert.EqualError(t, err, "FROM requires a valid TableName")
}

func TestParseSelect_selectCommandWithMultipleIdentifiers_createsValidStatement(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.IDENT, "id")
	t3 := newToken(tokenizer.IDENT, "name")
	t4 := newToken(tokenizer.FROM, "from")
	t5 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4, t5}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, err)
	assert.Equal(t, result.Type, SELECT)
	assert.Equal(t, result.Fields, []string{"id", "name"})
	assert.Equal(t, result.TableName, "users")
	assert.Empty(t, result.Values)
}

func TestParse_selectCommandWithSymbol_createsValidStatement(t *testing.T) {
    	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.IDENT, "*")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
    assert.Nil(t, err)
	assert.Equal(t, result.Type, SELECT)
	assert.Equal(t, result.Fields, []string{"*"})
	assert.Equal(t, result.TableName, "users")
	assert.Empty(t, result.Values)
}

func newToken(tokenType tokenizer.TokenType, value string) tokenizer.Token {
    return tokenizer.Token{Type: tokenType, Value: value}
}