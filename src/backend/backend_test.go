package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCatalog() *Catalog {
	return &Catalog{Tables: make(map[string]*Table)}
}

func TestCreateTable_newTable_addsTableToCatalog(t *testing.T) {
	// Arrange
	catalog := newCatalog()

	// Act
	table, err := catalog.CreateTable("users", []string{"id", "name"})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, table)
	assert.Contains(t, catalog.Tables, "users")
	assert.Equal(t, "users", catalog.Tables["users"].Name)
	assert.Equal(t, []string{"id", "name"}, catalog.Tables["users"].Columns)
}

func TestCreateTable_existingTable_returnsErrorAndKeepsOriginal(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	_, err := catalog.CreateTable("users", []string{"id"})
	require.NoError(t, err)

	// Act
	_, err = catalog.CreateTable("users", []string{"id", "name"})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, []string{"id"}, catalog.Tables["users"].Columns)
}

func TestCreateTable_zeroValueCatalog_addsTable(t *testing.T) {
	// Arrange
	var catalog Catalog

	// Act
	table, err := catalog.CreateTable("users", []string{"id"})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, table)
	assert.Contains(t, catalog.Tables, "users")
}

func TestCreateTable_returnedTable_isTheOneStoredInCatalog(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)

	// Act
	err = table.InsertRow([]string{"1", "Tim"})

	// Assert
	require.NoError(t, err)
	assert.Len(t, catalog.Tables["users"].Rows, 1)
}

func TestInsertRow_matchingColumnCount_appendsRow(t *testing.T) {
	// Arrange
	table := &Table{Name: "users", Columns: []string{"id", "name"}}

	// Act
	err := table.InsertRow([]string{"1", "Tim"})

	// Assert
	require.NoError(t, err)
	require.Len(t, table.Rows, 1)
	assert.Equal(t, []string{"1", "Tim"}, table.Rows[0].values)
}

func TestInsertRow_multipleRows_appendsAllRows(t *testing.T) {
	// Arrange
	table := &Table{Name: "users", Columns: []string{"id", "name"}}

	// Act
	err := table.InsertRow([]string{"1", "Tim"})
	require.NoError(t, err)
	err = table.InsertRow([]string{"2", "Sarah"})

	// Assert
	require.NoError(t, err)
	require.Len(t, table.Rows, 2)
	assert.Equal(t, []string{"1", "Tim"}, table.Rows[0].values)
	assert.Equal(t, []string{"2", "Sarah"}, table.Rows[1].values)
}

func TestInsertRow_valueCountMismatch_returnsErrorAndDoesNotAppend(t *testing.T) {
	testCases := []struct {
		name   string
		values []string
	}{
		{name: "too few values", values: []string{"1"}},
		{name: "too many values", values: []string{"1", "Tim", "extra"}},
		{name: "no values", values: []string{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			table := &Table{Name: "users", Columns: []string{"id", "name"}}

			// Act
			err := table.InsertRow(tc.values)

			// Assert
			assert.Error(t, err)
			assert.Empty(t, table.Rows)
		})
	}
}

func newUsersTable(t *testing.T) *Table {
	t.Helper()

	table := &Table{Name: "users", Columns: []string{"id", "name"}}
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))
	require.NoError(t, table.InsertRow([]string{"2", "Sarah"}))

	return table
}

func TestSelectRows_subsetOfFields_returnsProjectedRows(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{"name"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Tim"}, {"Sarah"}}, rows)
}

func TestSelectRows_fieldsInRequestedOrder_returnsValuesInThatOrder(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{"name", "id"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Tim", "1"}, {"Sarah", "2"}}, rows)
}

func TestSelectRows_star_expandsAllColumns(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{"*"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "Tim"}, {"2", "Sarah"}}, rows)
}

func TestSelectRows_duplicateFields_returnsDuplicatedValues(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{"id", "id"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "1"}, {"2", "2"}}, rows)
}

func TestSelectRows_unknownField_returnsError(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{"age"})

	// Assert
	assert.Nil(t, rows)
	assert.ErrorContains(t, err, "age")
}

func TestSelectRows_emptyFields_returnsError(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	rows, err := table.SelectRows([]string{})

	// Assert
	assert.Nil(t, rows)
	assert.Error(t, err)
}

func TestSelectRows_emptyTable_returnsNoRows(t *testing.T) {
	// Arrange
	table := &Table{Name: "users", Columns: []string{"id", "name"}}

	// Act
	rows, err := table.SelectRows([]string{"id"})

	// Assert
	require.NoError(t, err)
	assert.Empty(t, rows)
}

func TestSelectRows_returnedRowsAreCopies_doesNotMutateTable(t *testing.T) {
	// Arrange
	table := newUsersTable(t)
	rows, err := table.SelectRows([]string{"id", "name"})
	require.NoError(t, err)

	// Act
	rows[0][0] = "changed"

	// Assert
	assert.Equal(t, "1", table.Rows[0].values[0])
}

func TestCatalogSelectRows_existingTable_returnsRows(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))

	// Act
	rows, err := catalog.SelectRows("users", []string{"id"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1"}}, rows)
}

func TestCatalogSelectRows_missingTable_returnsError(t *testing.T) {
	// Arrange
	catalog := newCatalog()

	// Act
	rows, err := catalog.SelectRows("missing", []string{"id"})

	// Assert
	assert.Nil(t, rows)
	assert.Error(t, err)
}
