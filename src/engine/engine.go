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
	case parser.CREATE_TABLE:
		return executeCreateTable(statement)
	default:
		return fmt.Errorf("Unrecognized statement in %q", command)
	}
}

func executeSelect(statement *parser.Statement) error {
	columns, rows, err := catalog.SelectRows(statement.TableName, statement.Fields)
	if err != nil {
		return err
	}

	printTable(columns, rows)
	return nil
}

const minColumnWidth = 4

func printTable(columns []string, rows [][]string) {
	widths := make([]int, len(columns))
	for i, column := range columns {
		widths[i] = max(minColumnWidth, len(column))
	}

	for _, row := range rows {
		for i, value := range row {
			widths[i] = max(widths[i], len(value))
		}
	}

	fmt.Println(formatRow(columns, widths))
	fmt.Println(formatSeparator(widths))

	for _, row := range rows {
		fmt.Println(formatRow(row, widths))
	}
}

func formatRow(values []string, widths []int) string {
	cells := make([]string, len(values))
	for i, value := range values {
		cells[i] = fmt.Sprintf("%-*s", widths[i], value)
	}

	return strings.TrimRight(strings.Join(cells, " | "), " ")
}

func formatSeparator(widths []int) string {
	segments := make([]string, len(widths))
	for i, width := range widths {
		segmentWidth := width + 1
		if i > 0 {
			segmentWidth = width + 2
		}
		segments[i] = strings.Repeat("-", segmentWidth)
	}

	return strings.Join(segments, "+")
}

func executeInsert(statement *parser.Statement) error {
	if err := catalog.InsertRows(statement.TableName, statement.Values); err != nil {
		return err
	}

	fmt.Println("1 row inserted.")
	return nil
}

func executeCreateTable(statement *parser.Statement) error {
	_, err := catalog.CreateTable(statement.TableName, statement.Columns)
	if err != nil {
		return err
	}

	fmt.Printf("Table '%s' created.\n", statement.TableName)
	return nil
}
