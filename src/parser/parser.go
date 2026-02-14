package parser

import (
	"fmt"
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

func Parse(tokens []tokenizer.Token) (*Statement, error) {
	parserTree := ParserTree{tokens: tokens, position: 0}
	switch parserTree.peek().Type {
	case tokenizer.SELECT:
		return parserTree.createSelectStatement()
	case tokenizer.INSERT:
		return parserTree.createInsertStatement()
	default:
		return &Statement{}, fmt.Errorf("Syntax not valid. Not supported operation!")
	}
}