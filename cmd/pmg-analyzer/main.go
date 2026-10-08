package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.1.0-dev"

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && args[0] == "version" {
		fmt.Fprintln(stdout, "pmg-analyzer "+version)
		return 0
	}
	fmt.Fprintln(stderr, "usage: pmg-analyzer version")
	return 2
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
