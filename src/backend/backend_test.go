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
