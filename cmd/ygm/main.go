package main

import (
	"os"

	"github.com/CromulentConsulting/ygm-cli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
