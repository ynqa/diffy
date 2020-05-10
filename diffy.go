package main

import (
	"os"

	"github.com/ynqa/diffy/cmd"
)

func main() {
	cmd := cmd.New()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
