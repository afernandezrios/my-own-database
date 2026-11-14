package engine

import (
	"fmt"
	"learning/db/backend"
	"learning/db/parser"
	"learning/db/tokenizer"
	"strings"
)

var catalog = &backend.Catalog{Tables: make(map[string]*backend.Table)}

func Process(command string) error {
	tokens := tokenizer.Tokenize(command)

	statement, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	switch statement.Type {
	case parser.SELECT:
		return executeSelect(statement)
	case parser.INSERT:
		return executeInsert(statement)
	default:
		return fmt.Errorf("Unrecognized statement in %q", command)
	}
}

func executeSelect(statement *parser.Statement) error {
	rows, err := catalog.SelectRows(statement.TableName, statement.Fields)
	if err != nil {
		return err
	}

	for _, row := range rows {
		fmt.Println(strings.Join(row, ", "))
	}

	return nil
}

func executeInsert(statement *parser.Statement) error {
	return catalog.InsertRows(statement.TableName, statement.Values)
}
