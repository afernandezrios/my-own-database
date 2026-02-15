package parser

import (
	"fmt"
	"learning/db/tokenizer"
)

func (p *ParserTree) createInsertStatement() (*Statement, error) {
	
	
	if p.consume().Type != tokenizer.INSERT {
		return nil, fmt.Errorf("This is not an INSERT command")
	}

	if p.consume().Type != tokenizer.INTO {
		return nil, fmt.Errorf("This is not an INSERT command")
	}

	tableNameToken := p.consume()
	if tableNameToken.Type != tokenizer.IDENT {
		return nil, fmt.Errorf("Insert command must contain a table")
	}

	if p.consume().Type != tokenizer.VALUES {
		return nil, fmt.Errorf("Missing VALUE in INSERT command")
	}

	values:= p.consume().Value

	return &Statement{
		Type: INSERT,
		TableName: tableNameToken.Value,
		Values: values,
	}, nil

}
