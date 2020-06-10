package main

import (
	"os"

	command "github.com/mikefarah/yq/v3/cmd"
)

func main() {
	cmd := command.New()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
