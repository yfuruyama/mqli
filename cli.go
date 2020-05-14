package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
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

		stop := c.PrintProgressingMark()
		client := Client{c.projectID}
		resp, err := client.Query(line)
		stop()
		if err != nil {
			c.PrintInteractiveError(err)
			continue
		}

		fmt.Fprintf(c.out, "%#v\n", resp)
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
