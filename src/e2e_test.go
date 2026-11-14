package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var binaryPath string

// TestMain builds the real application binary once so the e2e tests can run it
// as a black box: stdin in, stdout out.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "db-e2e-*")
	if err != nil {
		panic("failed to create temp dir for e2e binary: " + err.Error())
	}

	binaryPath = filepath.Join(dir, "db")

	var buildOutput bytes.Buffer
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Stdout = &buildOutput
	build.Stderr = &buildOutput
	if err := build.Run(); err != nil {
		os.RemoveAll(dir)
		panic("failed to build e2e binary: " + err.Error() + "\n" + buildOutput.String())
	}

	exitCode := m.Run()
	os.RemoveAll(dir)
	os.Exit(exitCode)
}

// runApp executes the application with the given stdin and returns its stdout.
func runApp(t *testing.T, input string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Stdin = strings.NewReader(input)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.NoError(t, err, "stderr: %s", stderr.String())

	return stdout.String()
}

func TestE2E_createInsertSelect_printsConfirmationAndTable(t *testing.T) {
	// Arrange
	input := `CREATE TABLE users (id, name)
INSERT INTO users VALUES (1, 'Tim')
INSERT INTO users VALUES (2, 'Alex')
SELECT * FROM users
exit
`
	expected := `CLI started. Type 'exit' to close the application
db > Table 'users' created.
db > 1 row inserted.
db > 1 row inserted.
db > id   | name
-----+------
1    | Tim
2    | Alex
db > Bye!
`

	// Act
	output := runApp(t, input)

	// Assert
	assert.Equal(t, expected, output)
}

func TestE2E_insertIntoMissingTable_printsError(t *testing.T) {
	// Arrange
	input := `INSERT INTO non_existent_table VALUES (1)
`
	expected := `CLI started. Type 'exit' to close the application
db > Error: Table 'non_existent_table' does not exist.
`

	// Act
	output := runApp(t, input)

	// Assert
	assert.Equal(t, expected, output)
}
