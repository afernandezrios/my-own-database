package backend

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

type Catalog struct {
	Tables map[string]*Table
}

type Table struct {
	Name    string
	Columns []string
	Rows    []Row
	index   *BTreeIndex[int]
}

type Row struct {
	values []string // only string datatype at this moment
}

// Filter selects rows whose Column value equals Value.
type Filter struct {
	Column string
	Value  string
}

// If column "id" present, it creates an index on it
func (c *Catalog) CreateTable(name string, columns []string) (*Table, error) {
	if _, exists := c.Tables[name]; exists {
		return nil, fmt.Errorf("Table '%s' already exists.", name)
	}

	if c.Tables == nil {
		c.Tables = make(map[string]*Table)
	}

	var table = &Table{
		Name:    name,
		Columns: slices.Clone(columns),
	}

	if slices.Contains(columns, "id") {
		table.index = NewBtreeIndex[int]()
	}

	c.Tables[name] = table

	return table, nil
}

func (c *Catalog) InsertRows(name string, values []string) error {
	table, exists := c.Tables[name]
	if !exists {
		return fmt.Errorf("Table '%s' does not exist.", name)
	}

	return table.InsertRow(values)
}

func (c *Catalog) SelectRows(name string, fields []string, filter *Filter) ([]string, [][]string, error) {
	table, exists := c.Tables[name]
	if !exists {
		return nil, nil, fmt.Errorf("Table '%s' does not exist.", name)
	}

	return table.SelectRows(fields, filter)
}

func (t *Table) SelectRows(fields []string, filter *Filter) ([]string, [][]string, error) {
	if len(fields) == 0 {
		return nil, nil, errors.New("select requires at least one field")
	}

	if len(fields) == 1 && fields[0] == "*" {
		fields = t.Columns
	}

	indexes, err := t.resolveColumnIndexes(fields)
	if err != nil {
		return nil, nil, err
	}

	// Validate filter before choosing the search strategy.
	filterIndex := -1
	if filter != nil {
		filterIndex = slices.Index(t.Columns, filter.Column)
		if filterIndex == -1 {
			return nil, nil, fmt.Errorf("column '%s' does not exist in table '%s'", filter.Column, t.Name)
		}
	}

	selectedFields := slices.Clone(fields)

	if t.canUseIndex(filter) {
		rows, err := t.searchUsingIndex(indexes, filter)
		return selectedFields, rows, err
	}

	rows, err := t.sequentialSearch(indexes, filterIndex, filter)
	return selectedFields, rows, err
}

func (t *Table) resolveColumnIndexes(fields []string) ([]int, error) {
	indexes := make([]int, len(fields))

	for i, field := range fields {
		index := slices.Index(t.Columns, field)
		if index == -1 {
			return nil, fmt.Errorf("column '%s' does not exist in table '%s'", field, t.Name)
		}

		indexes[i] = index
	}

	return indexes, nil
}

func (t *Table) canUseIndex(filter *Filter) bool {
	return t.index != nil && filter != nil && filter.Column == "id"
}

func (t *Table) searchUsingIndex(indexes []int, filter *Filter) ([][]string, error) {
	id, err := strconv.Atoi(filter.Value)
	if err != nil {
		return [][]string{}, nil
	}

	offset, found := t.index.Search(id)
	if !found {
		return [][]string{}, nil
	}

	return [][]string{
		projectRow(&t.Rows[offset], indexes),
	}, nil
}

func (t *Table) sequentialSearch(indexes []int, filterIndex int, filter *Filter) ([][]string, error) {
	rows := make([][]string, 0, len(t.Rows))

	for _, row := range t.Rows {
		if filter != nil && row.values[filterIndex] != filter.Value {
			continue
		}

		rows = append(rows, projectRow(&row, indexes))
	}

	return rows, nil
}

func projectRow(row *Row, indexes []int) []string {
	projected := make([]string, len(indexes))

	for i, index := range indexes {
		projected[i] = row.values[index]
	}

	return projected
}

func (t *Table) InsertRow(values []string) error {
	id, err := t.validateRow(values)
	if err != nil {
		return err
	}

	t.Rows = append(t.Rows, Row{values: slices.Clone(values)})

	if t.index != nil {
		t.index.Insert(id, len(t.Rows)-1)
	}

	return nil
}

func (t *Table) validateRow(values []string) (int, error) {
	if len(t.Columns) != len(values) {
		return 0, fmt.Errorf("expected %d values, got %d", len(t.Columns), len(values))
	}

	if t.index == nil {
		return 0, nil
	}

	idColumn := slices.Index(t.Columns, "id")
	if idColumn == -1 {
		return 0, fmt.Errorf("cannot insert row: table %q has an index but no id column", t.Name)
	}

	id, err := strconv.Atoi(values[idColumn])
	if err != nil {
		return 0, fmt.Errorf(
			"invalid value '%s' for column '%s' in table '%s': id must be an integer", values[idColumn], "id", t.Name)
	}

	return id, nil
}
