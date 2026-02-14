package parser

import (
	"fmt"
	"learning/db/tokenizer"
)

func (p *ParserTree) createSelectStatement() (*Statement, error) {

	if p.consume().Type != tokenizer.SELECT {
		return nil, fmt.Errorf("This is not a SELECT command")
	}

	var selectTarget []string
	for isValidSelectTarget(p.peek()) {
		selectTarget = append(selectTarget, p.consume().Value)
	}

	if len(selectTarget) == 0 {
		return nil, fmt.Errorf("SELECT must contain at least one INDETIFIER or SYMBOL")
	}

	if p.consume().Type != tokenizer.FROM {
		return nil, fmt.Errorf("SELECT command requires FROM clause")
	}

	tableToken := p.consume()
	if tableToken.Type != tokenizer.IDENT {
		return nil, fmt.Errorf("FROM requires a valid TableName")
	}

	return &Statement{
		Type:      SELECT,
		Fields:    selectTarget,
		TableName: tableToken.Value,
	}, nil
}

func isValidSelectTarget(token tokenizer.Token) bool {
	return token.Type == tokenizer.IDENT || token.Type == tokenizer.SYMBOL
}