package main

import (
	"bufio"
	"fmt"
	"io"
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

		fmt.Fprintln(out, "Unrecognized command:", cleanedLine)
	}
}

func main() {
	fmt.Println("CLI started. Type 'exit' to close the application")
	StartREPL(os.Stdin, os.Stdout)
}
