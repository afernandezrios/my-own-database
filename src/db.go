package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"learning/db/engine"
	"learning/db/parser"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

const prompt = "db > "

func StartREPL(in io.Reader, out io.Writer) {
	if file, ok := in.(*os.File); ok && readline.IsTerminal(int(file.Fd())) {
		startInteractiveREPL(file, out)
		return
	}

	startStreamREPL(in, out)
}

// startInteractiveREPL runs a readline-based REPL when stdin is a terminal.
// readline switches the terminal to raw mode and interprets escape sequences,
// which is what makes the arrow keys move the cursor and edit the line.
func startInteractiveREPL(in *os.File, out io.Writer) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: prompt,
		Stdin:  in,
		Stdout: out,
		Stderr: out,
	})
	if err != nil {
		startStreamREPL(in, out)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			continue
		}
		if err != nil {
			// io.EOF (Ctrl-D) ends the session.
			return
		}
		if processCommand(line, out) {
			return
		}
	}
}

// startStreamREPL reads lines from a non-terminal reader, such as piped input
// or a test buffer, where no line editing is possible or needed.
func startStreamREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, prompt)

		if !scanner.Scan() {
			return
		}

		if processCommand(scanner.Text(), out) {
			return
		}
	}
}

// processCommand executes a single line of input. It reports whether the REPL
// should exit.
func processCommand(line string, out io.Writer) bool {
	cleanedLine := strings.TrimSpace(line)

	if cleanedLine == "exit" {
		fmt.Fprintln(out, "Bye!")
		return true
	}

	err := engine.Process(cleanedLine)
	if err != nil {
		if errors.Is(err, parser.ErrUnsupportedStatement) {
			fmt.Fprintf(out, "Unrecognized command: %s\n", cleanedLine)
			return false
		}
		fmt.Fprintf(out, "Error: %s\n", err.Error())
	}

	return false
}

func main() {
	fmt.Println("CLI started. Type 'exit' to close the application")
	StartREPL(os.Stdin, os.Stdout)
}
