package tokenizer

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	SELECT TokenType = iota
	INSERT
	INTO
	VALUES
	CREATE
	TABLE
	LPAREN
	RPAREN
	COMMA
	STRING
	NUMBER
	IDENT
	FROM
	SYMBOL
	WHERE
	COMPARATOR
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

var keywordsRegistry = map[string]Token{
	"select": {Type: SELECT, Value: "select"},
	"insert": {Type: INSERT, Value: "insert"},
	"into":   {Type: INTO, Value: "into"},
	"values": {Type: VALUES, Value: "values"},
	"create": {Type: CREATE, Value: "create"},
	"table":  {Type: TABLE, Value: "table"},
	"from":   {Type: FROM, Value: "from"},
	"where":  {Type: WHERE, Value: "where"},
	"*":      {Type: SYMBOL, Value: "*"},
	"=":      {Type: COMPARATOR, Value: "="},
}

func Tokenize(statement string) []Token {
	var tokens []Token
	runes := []rune(statement)
	i := 0

	for i < len(runes) {
		ch := runes[i]

		switch {
		case unicode.IsSpace(ch):
			i++

		case ch == '(':
			tokens = append(tokens, Token{Type: LPAREN, Value: "("})
			i++

		case ch == ')':
			tokens = append(tokens, Token{Type: RPAREN, Value: ")"})
			i++

		case ch == ',':
			tokens = append(tokens, Token{Type: COMMA, Value: ","})
			i++

		case ch == '"' || ch == '\'':
			quote := ch
			i++
			start := i
			for i < len(runes) && runes[i] != quote {
				i++
			}
			tokens = append(tokens, Token{Type: STRING, Value: string(runes[start:i])})
			if i < len(runes) {
				i++ // skip the closing quote
			}

		case unicode.IsNumber(ch):
			start := i
			for i < len(runes) && unicode.IsNumber(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{Type: NUMBER, Value: string(runes[start:i])})

		default:
			start := i
			for i < len(runes) && isIdentifierRune(runes[i]) {
				i++
			}
			if i == start {
				i++ // unknown punctuation still advances the cursor
			}
			value := strings.ToLower(string(runes[start:i]))
			if token, exists := keywordsRegistry[value]; exists {
				tokens = append(tokens, token)
			} else {
				tokens = append(tokens, Token{Type: IDENT, Value: value})
			}
		}
	}

	tokens = append(tokens, Token{Type: EOF})
	return tokens
}

func isIdentifierRune(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}
