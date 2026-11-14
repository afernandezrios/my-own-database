package parser

import (
	"fmt"
	"learning/db/tokenizer"
)

func (p *ParserTree) createCreateTableStatement() (*Statement, error) {

	if p.consume().Type != tokenizer.CREATE {
		return nil, fmt.Errorf("expected CREATE")
	}

	if p.consume().Type != tokenizer.TABLE {
		return nil, fmt.Errorf("expected TABLE after CREATE")
	}

	table := p.consume()
	if table.Type != tokenizer.IDENT {
		return nil, fmt.Errorf("CREATE TABLE requires a table name")
	}

	if p.consume().Type != tokenizer.LPAREN {
		return nil, fmt.Errorf("expected '(' after table name")
	}

	var columns []string
	for {
		token := p.consume()
		if token.Type != tokenizer.IDENT {
			return nil, fmt.Errorf("expected a column name, got %q", token.Value)
		}
		columns = append(columns, token.Value)

		switch p.consume().Type {
		case tokenizer.COMMA:
			continue
		case tokenizer.RPAREN:
		default:
			return nil, fmt.Errorf("expected ',' or ')' in column list")
		}
		break
	}

	if p.peek().Type != tokenizer.EOF {
		return nil, fmt.Errorf("unexpected token %q after CREATE TABLE", p.peek().Value)
	}

	return &Statement{
		Type:      CREATE_TABLE,
		TableName: table.Value,
		Columns:   columns,
	}, nil

}
