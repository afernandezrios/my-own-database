package backend

import (
	"strconv"
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
	assert.EqualError(t, err, `Table 'users' already exists.`)
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

func TestCreateTable_idColumn_createsIndex(t *testing.T) {
	// Arrange
	catalog := newCatalog()

	// Act
	withID, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	withoutID, err := catalog.CreateTable("logs", []string{"message"})
	require.NoError(t, err)

	// Assert
	assert.NotNil(t, withID.index)
	assert.Nil(t, withoutID.index)
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

func TestInsertRow_indexedTable_maintainsIndex(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)

	// Act
	require.NoError(t, table.InsertRow([]string{"7", "Tim"}))
	require.NoError(t, table.InsertRow([]string{"2", "Sarah"}))

	// Assert
	offset, found := table.index.Search(7)
	require.True(t, found)
	assert.Equal(t, 0, offset)
	offset, found = table.index.Search(2)
	require.True(t, found)
	assert.Equal(t, 1, offset)
	_, found = table.index.Search(99)
	assert.False(t, found)
}

func TestInsertRow_nonNumericID_returnsErrorAndDoesNotAppend(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)

	// Act
	err = table.InsertRow([]string{"abc", "Tim"})

	// Assert
	assert.EqualError(t, err, `invalid value 'abc' for column 'id' in table 'users': id must be an integer`)
	assert.Empty(t, table.Rows)
	_, found := table.index.Search(0)
	assert.False(t, found)

	// A failed insert must not corrupt the index for later inserts.
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))
	_, rows, selectErr := table.SelectRows([]string{"id"}, &Filter{Column: "id", Value: "1"})
	require.NoError(t, selectErr)
	assert.Equal(t, [][]string{{"1"}}, rows)
}

func TestInsertRow_duplicateID_indexPointsAtNewestRow(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))
	require.NoError(t, table.InsertRow([]string{"1", "Sarah"}))

	// Act
	_, pointRows, err := table.SelectRows([]string{"*"}, &Filter{Column: "id", Value: "1"})
	require.NoError(t, err)
	_, allRows, err := table.SelectRows([]string{"*"}, nil)
	require.NoError(t, err)

	// Assert: the point lookup resolves to the newest row; the scan still sees both.
	assert.Equal(t, [][]string{{"1", "Sarah"}}, pointRows)
	assert.Equal(t, [][]string{{"1", "Tim"}, {"1", "Sarah"}}, allRows)
}

func TestInsertRows_missingTable_returnsError(t *testing.T) {
	// Arrange
	catalog := newCatalog()

	// Act
	err := catalog.InsertRows("missing", []string{"1"})

	// Assert
	assert.EqualError(t, err, `Table 'missing' does not exist.`)
}

func newUsersTable(t *testing.T) *Table {
	t.Helper()

	table, err := newCatalog().CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))
	require.NoError(t, table.InsertRow([]string{"2", "Sarah"}))

	return table
}

func TestSelectRows_subsetOfFields_returnsProjectedRows(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"name"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, columns)
	assert.Equal(t, [][]string{{"Tim"}, {"Sarah"}}, rows)
}

func TestSelectRows_fieldsInRequestedOrder_returnsValuesInThatOrder(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"name", "id"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"name", "id"}, columns)
	assert.Equal(t, [][]string{{"Tim", "1"}, {"Sarah", "2"}}, rows)
}

func TestSelectRows_star_expandsAllColumns(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"*"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id", "name"}, columns)
	assert.Equal(t, [][]string{{"1", "Tim"}, {"2", "Sarah"}}, rows)
}

func TestSelectRows_duplicateFields_returnsDuplicatedValues(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"id", "id"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id", "id"}, columns)
	assert.Equal(t, [][]string{{"1", "1"}, {"2", "2"}}, rows)
}

func TestSelectRows_unknownField_returnsError(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"age"}, nil)

	// Assert
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.EqualError(t, err, `column 'age' does not exist in table 'users'`)
}

func TestSelectRows_emptyFields_returnsError(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{}, nil)

	// Assert
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.Error(t, err)
}

func TestSelectRows_emptyTable_returnsNoRows(t *testing.T) {
	// Arrange
	table := &Table{Name: "users", Columns: []string{"id", "name"}}

	// Act
	columns, rows, err := table.SelectRows([]string{"id"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id"}, columns)
	assert.Empty(t, rows)
}

func TestSelectRows_returnedRowsAreCopies_doesNotMutateTable(t *testing.T) {
	// Arrange
	table := newUsersTable(t)
	_, rows, err := table.SelectRows([]string{"id", "name"}, nil)
	require.NoError(t, err)

	// Act
	rows[0][0] = "changed"

	// Assert
	assert.Equal(t, "1", table.Rows[0].values[0])
}

func TestSelectRows_whereMatch_returnsOnlyMatchingRows(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"*"}, &Filter{Column: "id", Value: "1"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id", "name"}, columns)
	assert.Equal(t, [][]string{{"1", "Tim"}}, rows)
}

func TestSelectRows_whereNoMatch_returnsHeadersAndNoRows(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"*"}, &Filter{Column: "id", Value: "99"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id", "name"}, columns)
	assert.Empty(t, rows)
}

func TestSelectRows_whereID_intSemantics(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))

	// Act: ids are integers, so the filter value "01" addresses the same key as "1".
	_, rows, err := table.SelectRows([]string{"*"}, &Filter{Column: "id", Value: "01"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "Tim"}}, rows)
}

func TestSelectRows_whereUnknownColumn_returnsError(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"*"}, &Filter{Column: "age", Value: "30"})

	// Assert
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.EqualError(t, err, `column 'age' does not exist in table 'users'`)
}

func TestSelectRows_whereIDNonNumericFilter_returnsHeadersAndNoRows(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act
	columns, rows, err := table.SelectRows([]string{"id"}, &Filter{Column: "id", Value: "abc"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id"}, columns)
	assert.Empty(t, rows)
}

func TestSelectRows_whereName_indexPresent_usesScan(t *testing.T) {
	// Arrange
	table := newUsersTable(t)

	// Act: a non-id filter must stay on the scan path even when an index exists.
	columns, rows, err := table.SelectRows([]string{"name"}, &Filter{Column: "name", Value: "Sarah"})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"name"}, columns)
	assert.Equal(t, [][]string{{"Sarah"}}, rows)
}

func TestSelectRows_whereID_afterManyInserts_stillCorrect(t *testing.T) {
	// Arrange: enough rows to force several Rows backing-array reallocations,
	// which would invalidate stored pointers but not offsets.
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	for i := 1; i <= 150; i++ {
		require.NoError(t, table.InsertRow([]string{strconv.Itoa(i), "user-" + strconv.Itoa(i)}))
	}

	// Act & Assert
	for _, id := range []string{"1", "57", "150"} {
		_, rows, selectErr := table.SelectRows([]string{"id", "name"}, &Filter{Column: "id", Value: id})
		require.NoError(t, selectErr)
		assert.Equal(t, [][]string{{id, "user-" + id}}, rows)
	}
}

func TestSelectRows_tableWithoutIDColumn_whereID_returnsError(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("logs", []string{"message"})
	require.NoError(t, err)

	// Act
	_, _, err = table.SelectRows([]string{"message"}, &Filter{Column: "id", Value: "1"})

	// Assert
	assert.EqualError(t, err, `column 'id' does not exist in table 'logs'`)
}

func TestCatalogSelectRows_existingTable_returnsRows(t *testing.T) {
	// Arrange
	catalog := newCatalog()
	table, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
	require.NoError(t, table.InsertRow([]string{"1", "Tim"}))

	// Act
	columns, rows, err := catalog.SelectRows("users", []string{"id"}, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"id"}, columns)
	assert.Equal(t, [][]string{{"1"}}, rows)
}

func TestCatalogSelectRows_missingTable_returnsError(t *testing.T) {
	// Arrange
	catalog := newCatalog()

	// Act
	columns, rows, err := catalog.SelectRows("missing", []string{"id"}, nil)

	// Assert
	assert.Nil(t, columns)
	assert.Nil(t, rows)
	assert.EqualError(t, err, `Table 'missing' does not exist.`)
}
