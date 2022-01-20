package main

import (
	"os"

	"github.com/rqtx/pddoc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
