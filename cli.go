package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"
	"io"
	"strings"
	"time"
)

const (
	exitCodeSuccess = 0
	exitCodeError   = 1
)

type Cli struct {
	projectID string
	in io.ReadCloser
	out io.Writer
}

func NewCli(projectId string, in io.ReadCloser, out io.Writer) (*Cli, error) {
	return &Cli{
		projectID: projectId,
		in: in,
		out: out,
	}, nil
}

func (c *Cli) RunInteractive() int {
	rl, err := readline.NewEx(&readline.Config{
		Stdin:       c.in,
		HistoryFile: "/tmp/mqli_history",
	})
	if err != nil {
		return c.ExitOnError(err)
	}

	rl.SetPrompt("mql> ")
	for {
		line, err := rl.Readline()
		if err != nil {
			return c.ExitOnError(err)
		}
		if err == io.EOF {
			return c.Exit()
		}
		if strings.ToLower(line) == "exit" || strings.ToLower(line) == "quit" {
			return c.Exit()
		}

		stop := c.PrintProgressingMark()
		client := Client{c.projectID}
		result, err := client.Query(line)
		stop()
		if err != nil {
			c.PrintInteractiveError(err)
			continue
		}

		if len(result.Rows) > 0 {
			table := tablewriter.NewWriter(c.out)
			table.SetAutoFormatHeaders(false)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetAutoWrapText(false)
			for _, row := range result.Rows {
				table.Append(row.Columns)
			}
			table.SetHeader(result.Header)
			table.Render()
			fmt.Fprintf(c.out, "%d points in result\n\n", len(result.Rows))
		} else {
			fmt.Fprintf(c.out, "Empty result\n\n")
		}
	}
}

func (c *Cli) Exit() int {
	fmt.Fprintln(c.out, "Bye")
	return exitCodeSuccess
}

func (c *Cli) ExitOnError(err error) int {
	fmt.Fprintf(c.out, "ERROR: %s\n", err)
	return exitCodeError
}

func (c *Cli) PrintInteractiveError(err error) {
	fmt.Fprintf(c.out, "ERROR: %s\n", err)
}

func (c *Cli) PrintProgressingMark() func() {
	progressMarks := []string{`-`, `\`, `|`, `/`}
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		i := 0
		for {
			<-ticker.C
			mark := progressMarks[i%len(progressMarks)]
			fmt.Fprintf(c.out, "\r%s", mark)
			i++
		}
	}()

	stop := func() {
		ticker.Stop()
		fmt.Fprintf(c.out, "\r") // clear progressing mark
	}
	return stop
}
