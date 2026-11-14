package engine

import (
	"io"
	"os"
	"testing"

	"learning/db/backend"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetCatalog(t *testing.T) {
	t.Helper()

	catalog = &backend.Catalog{Tables: make(map[string]*backend.Table)}
	t.Cleanup(func() {
		catalog = &backend.Catalog{Tables: make(map[string]*backend.Table)}
	})
}

func createUsersTable(t *testing.T) {
	t.Helper()

	_, err := catalog.CreateTable("users", []string{"id", "name"})
	require.NoError(t, err)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	originalStdout := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()

	require.NoError(t, writer.Close())

	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())

	return string(output)
}

func TestProcess_selectExistingTable_printsRequestedFields(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))
	require.NoError(t, catalog.InsertRows("users", []string{"2", "Sarah"}))

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("select name from users")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "name\n------\nTim\nSarah\n", output)
}

func TestProcess_selectStar_printsAllColumns(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("select * from users")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "id   | name\n-----+------\n1    | Tim\n", output)
}

func TestProcess_selectMissingTable_returnsErrorAndPrintsNothing(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("select * from users")
	})

	// Assert
	assert.ErrorContains(t, processErr, `Table 'users' does not exist.`)
	assert.Empty(t, output)
}

func TestProcess_selectUnknownColumn_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("select age from users")

	// Assert
	assert.ErrorContains(t, err, `Column 'age' does not exist in table 'users'.`)
}

func TestProcess_selectWhereMatch_printsMatchingRow(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))
	require.NoError(t, catalog.InsertRows("users", []string{"2", "Sarah"}))

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("select * from users where id = 1")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "id   | name\n-----+------\n1    | Tim\n", output)
}

func TestProcess_selectWhereNoMatch_printsHeadersOnly(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))
	require.NoError(t, catalog.InsertRows("users", []string{"2", "Sarah"}))

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("select * from users where id = 99")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "id   | name\n-----+------\n", output)
}

func TestProcess_selectWhereUnknownColumn_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))

	// Act
	err := Process("select * from users where age = 30")

	// Assert
	assert.ErrorContains(t, err, `Column 'age' does not exist in table 'users'.`)
}

func TestProcess_insertExistingTable_addsRow(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("insert into users values ('1', 'Tim')")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "1 row inserted.\n", output)
	_, rows, err := catalog.SelectRows("users", []string{"id", "name"}, nil)
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "Tim"}}, rows)
}

func TestProcess_insertTwice_appendsBothRows(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	var processErr error
	output := captureStdout(t, func() {
		require.NoError(t, Process("insert into users values ('1', 'Tim')"))
		processErr = Process("insert into users values ('2', 'Sarah')")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "1 row inserted.\n1 row inserted.\n", output)
	_, rows, err := catalog.SelectRows("users", []string{"name"}, nil)
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Tim"}, {"Sarah"}}, rows)
}

func TestProcess_insertMissingTable_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	err := Process("insert into users values ('1', 'Tim')")

	// Assert
	assert.ErrorContains(t, err, `Table 'users' does not exist.`)
}

func TestProcess_insertValueCountMismatch_returnsErrorAndAddsNoRow(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("insert into users values ('1')")

	// Assert
	assert.Error(t, err)
	_, rows, selectErr := catalog.SelectRows("users", []string{"id", "name"}, nil)
	require.NoError(t, selectErr)
	assert.Empty(t, rows)
}

func TestProcess_invalidSyntax_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	err := Process("select from users")

	// Assert
	assert.Error(t, err)
}

func TestProcess_selectKeywordsAreCaseInsensitive_returnsRows(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)
	require.NoError(t, catalog.InsertRows("users", []string{"1", "Tim"}))

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("SeLeCt * FrOm UsErS")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "id   | name\n-----+------\n1    | Tim\n", output)
}

func TestProcess_createTable_createsTable(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	var processErr error
	output := captureStdout(t, func() {
		processErr = Process("CREATE TABLE users (id, name)")
	})

	// Assert
	require.NoError(t, processErr)
	assert.Equal(t, "Table 'users' created.\n", output)
	table, exists := catalog.Tables["users"]
	require.True(t, exists)
	assert.Equal(t, []string{"id", "name"}, table.Columns)
	assert.Empty(t, table.Rows)
}

func TestProcess_createTableSingleColumn_createsTable(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	var processErr error
	captureStdout(t, func() {
		processErr = Process("CREATE TABLE users (id)")
	})

	// Assert
	require.NoError(t, processErr)
	table, exists := catalog.Tables["users"]
	require.True(t, exists)
	assert.Equal(t, []string{"id"}, table.Columns)
}

func TestProcess_createTableTwice_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("create table users (id)")

	// Assert
	assert.ErrorContains(t, err, `Table 'users' already exists.`)
}

func TestProcess_createTableThenInsertAndSelect_roundTrip(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	var createErr error
	createOutput := captureStdout(t, func() {
		createErr = Process("CREATE TABLE users (id, name)")
	})

	var insertErr error
	insertOutput := captureStdout(t, func() {
		insertErr = Process("INSERT INTO users VALUES (1, 'Tim')")
	})

	var selectErr error
	selectOutput := captureStdout(t, func() {
		selectErr = Process("SELECT * FROM users")
	})

	// Assert
	require.NoError(t, createErr)
	assert.Equal(t, "Table 'users' created.\n", createOutput)
	require.NoError(t, insertErr)
	assert.Equal(t, "1 row inserted.\n", insertOutput)
	require.NoError(t, selectErr)
	assert.Equal(t, "id   | name\n-----+------\n1    | Tim\n", selectOutput)
}

func TestProcess_createTableInvalidSyntax_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	err := Process("CREATE TABLE users id, name")

	// Assert
	assert.Error(t, err)
	assert.NotContains(t, catalog.Tables, "users")
}
