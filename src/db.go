package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"learning/db/parser"
	"learning/db/tokenizer"
	"os"
	"strings"
)

func StartREPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, "db > ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		cleanedLine := strings.TrimSpace(line)

		if cleanedLine == "exit" {
			fmt.Fprintln(out, "Bye!")
			return
		}

		tokens := tokenizer.Tokenize(cleanedLine)
		statement, err := parser.Parse(tokens)

		if err != nil {
			if errors.Is(err, parser.ErrUnsupportedStatement) {
				fmt.Fprintf(out, "Unrecognized command: %s\n", cleanedLine)
				continue
			}
			fmt.Fprintln(out, err.Error())
			return
		}

		fmt.Fprintln(out, "STATEMENT TYPE: ", statement.Type)
		fmt.Fprintln(out, "STATEMENT FIELDS: ", statement.Fields)
		fmt.Fprintln(out, "STATEMENT TABLE: ", statement.TableName)
		fmt.Fprintln(out, "STATEMENT VALUES: ", statement.Values)
	}
}

func main() {
	fmt.Println("CLI started. Type 'exit' to close the application")
	StartREPL(os.Stdin, os.Stdout)
}
