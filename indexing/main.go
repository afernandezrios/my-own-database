// Command loadtest measures the end-to-end performance of the database: INSERTs
// and point SELECTs through the real SQL engine pipeline. Tables with an "id"
// column carry an in-memory B-tree index, so point lookups avoid full scans.
//
// Run it from the repository root:
//
//	go -C indexing run .
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"learning/db/engine"
)

const (
	smallRows   = 10
	mediumRows  = 10_000
	largeRows   = 100_000
	selectCount = 100_000
)

// result is the outcome of one measured scenario.
type result struct {
	operation string
	rows      int
	count     int
	total     time.Duration
}

func (r result) avg() time.Duration {
	return r.total / time.Duration(r.count)
}

func (r result) opsPerSecond() float64 {
	return float64(r.count) / r.total.Seconds()
}

func main() {
	realStdout := os.Stdout

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		panic("failed to open " + os.DevNull + ": " + err.Error())
	}
	defer devNull.Close()

	// The engine prints confirmations and result tables to stdout. Discard
	// that output while measuring so the terminal isn't flooded; the
	// formatting work is still part of the measured time.
	os.Stdout = devNull

	smallInsert := measureInsert(smallRows)
	mediumInsert := measureInsert(mediumRows)
	largeInsert := measureInsert(largeRows)
	smallSelect := measureSelect(smallRows)
	mediumSelect := measureSelect(mediumRows)
	largeSelect := measureSelect(largeRows)

	os.Stdout = realStdout
	printResults(realStdout, []result{smallInsert, mediumInsert, largeInsert, smallSelect, mediumSelect, largeSelect})
}

// measureInsert times inserting rows one SQL statement at a time. The table is
// created before the clock starts.
func measureInsert(rows int) result {
	table := fmt.Sprintf("users_%d", rows)
	mustRun(fmt.Sprintf("CREATE TABLE %s (id, name)", table))

	start := time.Now()
	for i := 1; i <= rows; i++ {
		mustRun(fmt.Sprintf("INSERT INTO %s VALUES (%d, 'user-%d')", table, i, i))
	}

	return result{operation: "INSERT", rows: rows, count: rows, total: time.Since(start)}
}

// measureSelect times point queries on an already populated table.
func measureSelect(rows int) result {
	table := fmt.Sprintf("users_%d", rows)

	// One untimed query so first-call allocations don't skew the average.
	mustRun(fmt.Sprintf("SELECT * FROM %s WHERE id = 1", table))

	start := time.Now()
	for i := 0; i < selectCount; i++ {
		id := i%rows + 1
		mustRun(fmt.Sprintf("SELECT * FROM %s WHERE id = %d", table, id))
	}

	return result{operation: "SELECT", rows: rows, count: selectCount, total: time.Since(start)}
}

func mustRun(command string) {
	if err := engine.Process(command); err != nil {
		panic(fmt.Sprintf("load test command failed: %q: %v", command, err))
	}
}

func printResults(out io.Writer, results []result) {
	headers := []string{"Operation", "Rows", "Count", "Total", "Avg/op", "Ops/sec"}
	rightAligned := []bool{false, true, true, true, true, true}

	rows := make([][]string, len(results))
	for i, r := range results {
		rows[i] = []string{
			r.operation,
			fmt.Sprintf("%d", r.rows),
			fmt.Sprintf("%d", r.count),
			r.total.String(),
			r.avg().String(),
			fmt.Sprintf("%.0f", r.opsPerSecond()),
		}
	}

	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			widths[i] = max(widths[i], len(cell))
		}
	}

	separator := buildSeparator(widths)

	fmt.Fprintln(out, "Load test results")
	fmt.Fprintln(out)
	fmt.Fprintln(out, separator)
	fmt.Fprintln(out, buildRow(headers, widths, rightAligned))
	fmt.Fprintln(out, separator)
	for _, row := range rows {
		fmt.Fprintln(out, buildRow(row, widths, rightAligned))
	}
	fmt.Fprintln(out, separator)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "SELECT is a point query (WHERE id = N) measured end-to-end through the engine;")
	fmt.Fprintln(out, "Count is the number of operations measured, not the number of result rows.")
}

func buildRow(values []string, widths []int, rightAligned []bool) string {
	cells := make([]string, len(values))
	for i, value := range values {
		if rightAligned[i] {
			cells[i] = fmt.Sprintf("%*s", widths[i], value)
		} else {
			cells[i] = fmt.Sprintf("%-*s", widths[i], value)
		}
	}

	return "| " + strings.Join(cells, " | ") + " |"
}

func buildSeparator(widths []int) string {
	segments := make([]string, len(widths))
	for i, width := range widths {
		segments[i] = strings.Repeat("-", width+2)
	}

	return "+" + strings.Join(segments, "+") + "+"
}
