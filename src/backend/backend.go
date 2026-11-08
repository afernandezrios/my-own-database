package backend

import (
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

var catalog = &Catalog{Tables: make(map[string]*Table)}

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

func (t *Table) InsertRow(values []string) error {
	if len(t.Columns) != len(values) {
		return fmt.Errorf("Values doesn't match Row columns")
	}

	t.Rows = append(t.Rows, Row{values: values})
	return nil
}
