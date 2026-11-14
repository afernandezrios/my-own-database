package parser

import (
	"learning/db/tokenizer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestParseSelect_whereClause_createsCondition(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.SYMBOL, "*")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.IDENT, "users")
	t5 := newToken(tokenizer.WHERE, "where")
	t6 := newToken(tokenizer.IDENT, "id")
	t7 := newToken(tokenizer.COMPARATOR, "=")
	t8 := newToken(tokenizer.NUMBER, "1")
	tokens := []tokenizer.Token{t1, t2, t3, t4, t5, t6, t7, t8}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Where)
	assert.Equal(t, "id", result.Where.Column)
	assert.Equal(t, "1", result.Where.Value)
}

func TestParseSelect_noWhereClause_hasNoCondition(t *testing.T) {
	// Arrange
	t1 := newToken(tokenizer.SELECT, "select")
	t2 := newToken(tokenizer.SYMBOL, "*")
	t3 := newToken(tokenizer.FROM, "from")
	t4 := newToken(tokenizer.IDENT, "users")
	tokens := []tokenizer.Token{t1, t2, t3, t4}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Nil(t, result.Where)
}

func TestParseSelect_whereWithoutColumn_returnsError(t *testing.T) {
	// Arrange
	tokens := []tokenizer.Token{
		newToken(tokenizer.SELECT, "select"),
		newToken(tokenizer.SYMBOL, "*"),
		newToken(tokenizer.FROM, "from"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.WHERE, "where"),
		newToken(tokenizer.FROM, "from"),
	}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "WHERE requires a column name")
}

func TestParseSelect_whereWithoutComparator_returnsError(t *testing.T) {
	// Arrange
	tokens := []tokenizer.Token{
		newToken(tokenizer.SELECT, "select"),
		newToken(tokenizer.SYMBOL, "*"),
		newToken(tokenizer.FROM, "from"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.WHERE, "where"),
		newToken(tokenizer.IDENT, "id"),
		newToken(tokenizer.FROM, "from"),
	}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "WHERE requires a comparison operator")
}

func TestParseSelect_whereWithoutValue_returnsError(t *testing.T) {
	// Arrange
	tokens := []tokenizer.Token{
		newToken(tokenizer.SELECT, "select"),
		newToken(tokenizer.SYMBOL, "*"),
		newToken(tokenizer.FROM, "from"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.WHERE, "where"),
		newToken(tokenizer.IDENT, "id"),
		newToken(tokenizer.COMPARATOR, "="),
	}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, "WHERE requires a value")
}

func TestParseSelect_trailingToken_returnsError(t *testing.T) {
	// Arrange
	tokens := []tokenizer.Token{
		newToken(tokenizer.SELECT, "select"),
		newToken(tokenizer.SYMBOL, "*"),
		newToken(tokenizer.FROM, "from"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.WHERE, "where"),
		newToken(tokenizer.IDENT, "id"),
		newToken(tokenizer.COMPARATOR, "="),
		newToken(tokenizer.NUMBER, "1"),
		newToken(tokenizer.IDENT, "extra"),
	}
	parser := ParserTree{tokens: tokens, position: 0}

	// Act
	result, err := parser.createSelectStatement()

	// Assert
	assert.Nil(t, result)
	assert.EqualError(t, err, `unexpected token "extra" after SELECT`)
}

func newToken(tokenType tokenizer.TokenType, value string) tokenizer.Token {
    return tokenizer.Token{Type: tokenType, Value: value}
}