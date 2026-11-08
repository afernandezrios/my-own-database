package parser

import (
	"fmt"
	"learning/db/tokenizer"
)

func (p *ParserTree) createInsertStatement() (*Statement, error) {

	if p.consume().Type != tokenizer.INSERT {
		return nil, fmt.Errorf("expected INSERT")
	}

	if p.consume().Type != tokenizer.INTO {
		return nil, fmt.Errorf("expected INTO after INSERT")
	}

	table := p.consume()
	if table.Type != tokenizer.IDENT {
		return nil, fmt.Errorf("INSERT requires a table name")
	}

	if p.consume().Type != tokenizer.VALUES {
		return nil, fmt.Errorf("expected VALUES after table name")
	}

	if p.consume().Type != tokenizer.LPAREN {
		return nil, fmt.Errorf("expected '(' after VALUES")
	}

	var values []string
	for {
		token := p.consume()
		if token.Type != tokenizer.STRING && token.Type != tokenizer.NUMBER {
			return nil, fmt.Errorf("expected a value, got %q", token.Value)
		}
		values = append(values, token.Value)

		switch p.consume().Type {
		case tokenizer.COMMA:
			continue
		case tokenizer.RPAREN:
		default:
			return nil, fmt.Errorf("expected ',' or ')' in VALUES list")
		}
		break
	}

	if p.peek().Type != tokenizer.EOF {
		return nil, fmt.Errorf("unexpected token %q after INSERT", p.peek().Value)
	}

	return &Statement{
		Type:      INSERT,
		TableName: table.Value,
		Values:    values,
	}, nil

}
