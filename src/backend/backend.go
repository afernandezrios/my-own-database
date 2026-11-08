package backend

import (
	"errors"
	"fmt"
	"slices"
)

type Catalog struct {
	Tables map[string]*Table
}

type Table struct {
	Name    string
	Columns []string
	Rows    []Row
}

type Row struct {
	values []string // only string datatype at this moment
}

func (c *Catalog) CreateTable(name string, columns []string) (*Table, error) {
	if _, exists := c.Tables[name]; exists {
		return nil, fmt.Errorf("Table %q already exists", name)
	}

	if c.Tables == nil {
		c.Tables = make(map[string]*Table)
	}

	var table = &Table{
		Name:    name,
		Columns: slices.Clone(columns),
	}
	c.Tables[name] = table

	return table, nil
}

func (c *Catalog) SelectRows(name string, fields []string) ([][]string, error) {
	table, exists := c.Tables[name]
	if !exists {
		return nil, fmt.Errorf("table %q does not exist", name)
	}

	return table.SelectRows(fields)
}

func (t *Table) SelectRows(fields []string) ([][]string, error) {
	if len(fields) == 0 {
		return nil, errors.New("select requires at least one field")
	}

	if len(fields) == 1 && fields[0] == "*" {
		fields = t.Columns
	}

	indexes := make([]int, len(fields))
	for i, field := range fields {
		index := slices.Index(t.Columns, field)
		if index == -1 {
			return nil, fmt.Errorf("column %q does not exist in table %q", field, t.Name)
		}
		indexes[i] = index
	}

	rows := make([][]string, 0, len(t.Rows))
	for _, row := range t.Rows {
		projected := make([]string, len(indexes))
		for i, index := range indexes {
			projected[i] = row.values[index]
		}
		rows = append(rows, projected)
	}

	return rows, nil
}

func (t *Table) InsertRow(values []string) error {
	if len(t.Columns) != len(values) {
		return fmt.Errorf("Values doesn't match Row columns")
	}

	t.Rows = append(t.Rows, Row{values: values})
	return nil
}
