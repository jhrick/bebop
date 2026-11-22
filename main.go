package main

import (
	"fmt"
	"os"

	"github.com/jhrick/bop/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "bop: %v\n", err)
		os.Exit(1)
	}
}
