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
	assert.Equal(t, "Tim\nSarah\n", output)
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
	assert.Equal(t, "1, Tim\n", output)
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
	assert.ErrorContains(t, processErr, `table "users" does not exist`)
	assert.Empty(t, output)
}

func TestProcess_selectUnknownColumn_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("select age from users")

	// Assert
	assert.ErrorContains(t, err, `column "age" does not exist`)
}

func TestProcess_insertExistingTable_addsRow(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("insert into users values ('1', 'Tim')")

	// Assert
	require.NoError(t, err)
	rows, err := catalog.SelectRows("users", []string{"id", "name"})
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"1", "Tim"}}, rows)
}

func TestProcess_insertTwice_appendsBothRows(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	require.NoError(t, Process("insert into users values ('1', 'Tim')"))
	err := Process("insert into users values ('2', 'Sarah')")

	// Assert
	require.NoError(t, err)
	rows, err := catalog.SelectRows("users", []string{"name"})
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Tim"}, {"Sarah"}}, rows)
}

func TestProcess_insertMissingTable_returnsError(t *testing.T) {
	// Arrange
	resetCatalog(t)

	// Act
	err := Process("insert into users values ('1', 'Tim')")

	// Assert
	assert.ErrorContains(t, err, `table "users" does not exist`)
}

func TestProcess_insertValueCountMismatch_returnsErrorAndAddsNoRow(t *testing.T) {
	// Arrange
	resetCatalog(t)
	createUsersTable(t)

	// Act
	err := Process("insert into users values ('1')")

	// Assert
	assert.Error(t, err)
	rows, selectErr := catalog.SelectRows("users", []string{"id", "name"})
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
	assert.Equal(t, "1, Tim\n", output)
}
