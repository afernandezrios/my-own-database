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

	statement := &Statement{
		Type:      SELECT,
		Fields:    selectTarget,
		TableName: tableToken.Value,
	}

	if p.peek().Type == tokenizer.WHERE {
		condition, err := p.createWhereCondition()
		if err != nil {
			return nil, err
		}
		statement.Where = condition
	}

	if p.peek().Type != tokenizer.EOF {
		return nil, fmt.Errorf("unexpected token %q after SELECT", p.peek().Value)
	}

	return statement, nil
}

func (p *ParserTree) createWhereCondition() (*Condition, error) {
	p.consume() // WHERE

	column := p.consume()
	if column.Type != tokenizer.IDENT {
		return nil, fmt.Errorf("WHERE requires a column name")
	}

	if p.consume().Type != tokenizer.COMPARATOR {
		return nil, fmt.Errorf("WHERE requires a comparison operator")
	}

	value := p.consume()
	if value.Type != tokenizer.STRING && value.Type != tokenizer.NUMBER {
		return nil, fmt.Errorf("WHERE requires a value")
	}

	return &Condition{Column: column.Value, Value: value.Value}, nil
}

func isValidSelectTarget(token tokenizer.Token) bool {
	return token.Type == tokenizer.IDENT || token.Type == tokenizer.SYMBOL
}