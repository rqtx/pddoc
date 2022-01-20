package main

import (
	"os"

	"github.com/rqtx/pdoc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
