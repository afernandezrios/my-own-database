package parser

import (
	"errors"
	"learning/db/tokenizer"
)

type StatementType int

const (
	SELECT StatementType = iota
	INSERT
)

type Statement struct {
	Type      StatementType
	Fields    []string
	TableName string
	Values    []string
}

// ErrUnsupportedStatement is returned when the first token is not a statement
// type the parser knows how to handle.
var ErrUnsupportedStatement = errors.New("Syntax not valid. Not supported operation!")

func Parse(tokens []tokenizer.Token) (*Statement, error) {
	parserTree := ParserTree{tokens: tokens, position: 0}
	switch parserTree.peek().Type {
	case tokenizer.SELECT:
		return parserTree.createSelectStatement()
	case tokenizer.INSERT:
		return parserTree.createInsertStatement()
	default:
		return nil, ErrUnsupportedStatement
	}
}
