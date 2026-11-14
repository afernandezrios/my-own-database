package parser

import (
	"testing"

	"learning/db/tokenizer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTableTokens builds the token stream for "CREATE TABLE users ..." and
// appends the EOF token the tokenizer always produces.
func createTableTokens(tail ...tokenizer.Token) []tokenizer.Token {
	tokens := []tokenizer.Token{
		newToken(tokenizer.CREATE, "create"),
		newToken(tokenizer.TABLE, "table"),
		newToken(tokenizer.IDENT, "users"),
	}
	tokens = append(tokens, tail...)
	return append(tokens, newToken(tokenizer.EOF, ""))
}

func TestParseCreateTable_validStatement_createsStatement(t *testing.T) {
	testCases := []struct {
		name            string
		tokens          []tokenizer.Token
		expectedTable   string
		expectedColumns []string
	}{
		{
			name: "single column",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:   "users",
			expectedColumns: []string{"id"},
		},
		{
			name: "two columns",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.IDENT, "name"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:   "users",
			expectedColumns: []string{"id", "name"},
		},
		{
			name: "three columns",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.IDENT, "name"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.IDENT, "age"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:   "users",
			expectedColumns: []string{"id", "name", "age"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			parser := ParserTree{tokens: tc.tokens, position: 0}

			// Act
			result, err := parser.createCreateTableStatement()

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, CREATE_TABLE, result.Type)
			assert.Equal(t, tc.expectedTable, result.TableName)
			assert.Equal(t, tc.expectedColumns, result.Columns)
			assert.Empty(t, result.Fields)
			assert.Empty(t, result.Values)
		})
	}
}

func TestParseCreateTable_invalidStatement_returnsError(t *testing.T) {
	testCases := []struct {
		name        string
		tokens      []tokenizer.Token
		expectedErr string
	}{
		{
			name:        "empty token stream",
			tokens:      []tokenizer.Token{newToken(tokenizer.EOF, "")},
			expectedErr: "expected CREATE",
		},
		{
			name: "does not start with CREATE",
			tokens: []tokenizer.Token{
				newToken(tokenizer.SELECT, "select"),
				newToken(tokenizer.TABLE, "table"),
				newToken(tokenizer.IDENT, "users"),
				newToken(tokenizer.EOF, ""),
			},
			expectedErr: "expected CREATE",
		},
		{
			name: "missing TABLE",
			tokens: []tokenizer.Token{
				newToken(tokenizer.CREATE, "create"),
				newToken(tokenizer.IDENT, "users"),
				newToken(tokenizer.EOF, ""),
			},
			expectedErr: "expected TABLE after CREATE",
		},
		{
			name: "missing table name",
			tokens: []tokenizer.Token{
				newToken(tokenizer.CREATE, "create"),
				newToken(tokenizer.TABLE, "table"),
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.EOF, ""),
			},
			expectedErr: "CREATE TABLE requires a table name",
		},
		{
			name: "table name is not an identifier",
			tokens: []tokenizer.Token{
				newToken(tokenizer.CREATE, "create"),
				newToken(tokenizer.TABLE, "table"),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.EOF, ""),
			},
			expectedErr: "CREATE TABLE requires a table name",
		},
		{
			name: "missing opening parenthesis",
			tokens: createTableTokens(
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedErr: "expected '(' after table name",
		},
		{
			name: "empty column list",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedErr: `expected a column name, got ")"`,
		},
		{
			name: "column is not an identifier",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedErr: `expected a column name, got "1"`,
		},
		{
			name: "missing comma between columns",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.IDENT, "name"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedErr: "expected ',' or ')' in column list",
		},
		{
			name: "trailing comma",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedErr: `expected a column name, got ")"`,
		},
		{
			name: "missing closing parenthesis",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
			),
			expectedErr: "expected ',' or ')' in column list",
		},
		{
			name: "tokens after closing parenthesis",
			tokens: createTableTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.RPAREN, ")"),
				newToken(tokenizer.IDENT, "garbage"),
			),
			expectedErr: `unexpected token "garbage" after CREATE TABLE`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			parser := ParserTree{tokens: tc.tokens, position: 0}

			// Act
			result, err := parser.createCreateTableStatement()

			// Assert
			assert.Nil(t, result)
			assert.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestParse_createTable_returnsStatement(t *testing.T) {
	// Arrange
	tokens := []tokenizer.Token{
		newToken(tokenizer.CREATE, "create"),
		newToken(tokenizer.TABLE, "table"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.LPAREN, "("),
		newToken(tokenizer.IDENT, "id"),
		newToken(tokenizer.COMMA, ","),
		newToken(tokenizer.IDENT, "name"),
		newToken(tokenizer.RPAREN, ")"),
		newToken(tokenizer.EOF, ""),
	}

	// Act
	result, err := Parse(tokens)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, CREATE_TABLE, result.Type)
	assert.Equal(t, "users", result.TableName)
	assert.Equal(t, []string{"id", "name"}, result.Columns)
}
