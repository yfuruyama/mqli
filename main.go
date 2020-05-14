package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var projectID string

	flag.StringVar(&projectID, "p", "", "")
	flag.Parse()

	if projectID == "" {
		flag.Usage()
		os.Exit(1)
	}

	cli, err := NewCli(projectID, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := cli.RunInteractive()
	os.Exit(exitCode)
}
