package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestStartREPL_unrecognizedCommand_printsError(t *testing.T) {
	// Arrange
	input := strings.NewReader("hello\nexit\n")
	var output bytes.Buffer

	// Act
	StartREPL(input, &output)

	// Assert
	got := output.String()
	expectedError := "Unrecognized command: hello"

	if !strings.Contains(got, expectedError) {
		t.Errorf("Expected output to contain %q, but got:\n%q", expectedError, got)
	}
}

func TestStartREPL_unrecognizedCommand_continuesReadingCommands(t *testing.T) {
	// Arrange
	input := strings.NewReader("hello\nexit\n")
	var output bytes.Buffer

	// Act
	StartREPL(input, &output)

	// Assert
	got := output.String()
	if !strings.HasSuffix(got, "Bye!\n") {
		t.Errorf("Expected REPL to keep reading after an unrecognized command and exit on 'exit', but got:\n%q", got)
	}
}

func TestStartREPL_exitCommand_closeProgramWithByeText(t *testing.T) {
	// Arrange
	input := strings.NewReader("exit")
	var output bytes.Buffer

	// Act
	StartREPL(input, &output)

	// Assert
	got := output.String()
	newPrompt := "Bye!\n"

	if !strings.HasSuffix(got, newPrompt) {
		t.Errorf("Expected output to ends with %q, but got:\n%q", newPrompt, got)
	}
}
