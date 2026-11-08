package parser

import (
	"testing"

	"learning/db/tokenizer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// insertTokens builds the token stream for "INSERT INTO users VALUES ..." and
// appends the EOF token the tokenizer always produces.
func insertTokens(tail ...tokenizer.Token) []tokenizer.Token {
	tokens := []tokenizer.Token{
		newToken(tokenizer.INSERT, "insert"),
		newToken(tokenizer.INTO, "into"),
		newToken(tokenizer.IDENT, "users"),
		newToken(tokenizer.VALUES, "values"),
	}
	tokens = append(tokens, tail...)
	return append(tokens, newToken(tokenizer.EOF, ""))
}

func TestParseInsert_validStatement_createsStatement(t *testing.T) {
	testCases := []struct {
		name           string
		tokens         []tokenizer.Token
		expectedTable  string
		expectedValues []string
	}{
		{
			name: "single number value",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:  "users",
			expectedValues: []string{"1"},
		},
		{
			name: "number and string values",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.STRING, "Tim"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:  "users",
			expectedValues: []string{"1", "Tim"},
		},
		{
			name: "multiple string values",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.STRING, "Tim"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.STRING, "Sarah"),
				newToken(tokenizer.RPAREN, ")"),
			),
			expectedTable:  "users",
			expectedValues: []string{"Tim", "Sarah"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			parser := ParserTree{tokens: tc.tokens, position: 0}

			// Act
			result, err := parser.createInsertStatement()

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, INSERT, result.Type)
			assert.Equal(t, tc.expectedTable, result.TableName)
			assert.Equal(t, tc.expectedValues, result.Values)
			assert.Empty(t, result.Fields)
		})
	}
}

func TestParseInsert_invalidStatement_returnsError(t *testing.T) {
	testCases := []struct {
		name   string
		tokens []tokenizer.Token
	}{
		{
			name:   "empty token stream",
			tokens: []tokenizer.Token{newToken(tokenizer.EOF, "")},
		},
		{
			name: "does not start with INSERT",
			tokens: []tokenizer.Token{
				newToken(tokenizer.SELECT, "select"),
				newToken(tokenizer.INTO, "into"),
				newToken(tokenizer.IDENT, "users"),
				newToken(tokenizer.VALUES, "values"),
				newToken(tokenizer.EOF, ""),
			},
		},
		{
			name: "missing INTO",
			tokens: []tokenizer.Token{
				newToken(tokenizer.INSERT, "insert"),
				newToken(tokenizer.IDENT, "users"),
				newToken(tokenizer.VALUES, "values"),
				newToken(tokenizer.EOF, ""),
			},
		},
		{
			name: "missing table name",
			tokens: []tokenizer.Token{
				newToken(tokenizer.INSERT, "insert"),
				newToken(tokenizer.INTO, "into"),
				newToken(tokenizer.VALUES, "values"),
				newToken(tokenizer.EOF, ""),
			},
		},
		{
			name: "table name is not an identifier",
			tokens: []tokenizer.Token{
				newToken(tokenizer.INSERT, "insert"),
				newToken(tokenizer.INTO, "into"),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.VALUES, "values"),
				newToken(tokenizer.EOF, ""),
			},
		},
		{
			name: "missing VALUES",
			tokens: []tokenizer.Token{
				newToken(tokenizer.INSERT, "insert"),
				newToken(tokenizer.INTO, "into"),
				newToken(tokenizer.IDENT, "users"),
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.RPAREN, ")"),
				newToken(tokenizer.EOF, ""),
			},
		},
		{
			name: "missing opening parenthesis",
			tokens: insertTokens(
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.RPAREN, ")"),
			),
		},
		{
			name: "empty value list",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.RPAREN, ")"),
			),
		},
		{
			name: "value is not a literal",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.IDENT, "id"),
				newToken(tokenizer.RPAREN, ")"),
			),
		},
		{
			name: "missing comma between values",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.NUMBER, "2"),
				newToken(tokenizer.RPAREN, ")"),
			),
		},
		{
			name: "trailing comma",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.COMMA, ","),
				newToken(tokenizer.RPAREN, ")"),
			),
		},
		{
			name: "missing closing parenthesis",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
			),
		},
		{
			name: "tokens after closing parenthesis",
			tokens: insertTokens(
				newToken(tokenizer.LPAREN, "("),
				newToken(tokenizer.NUMBER, "1"),
				newToken(tokenizer.RPAREN, ")"),
				newToken(tokenizer.IDENT, "garbage"),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			parser := ParserTree{tokens: tc.tokens, position: 0}

			// Act
			result, err := parser.createInsertStatement()

			// Assert
			assert.Nil(t, result)
			assert.Error(t, err)
		})
	}
}
