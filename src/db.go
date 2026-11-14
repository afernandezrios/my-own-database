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

		err := engine.Process(cleanedLine)
		if err != nil {
			if errors.Is(err, parser.ErrUnsupportedStatement) {
				fmt.Fprintf(out, "Unrecognized command: %s\n", cleanedLine)
				continue
			}
			fmt.Fprintln(out, err.Error())
			return
		}
	}
}

func main() {
	fmt.Println("CLI started. Type 'exit' to close the application")
	StartREPL(os.Stdin, os.Stdout)
}
