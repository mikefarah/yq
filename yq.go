package main

import (
	"os"

	command "github.com/mikefarah/yq/v4/cmd"
	"github.com/mikefarah/yq/v4/pkg/properties"
	"github.com/mikefarah/yq/v4/pkg/xml"
)

func main() {
	properties.RegisterPropertiesFormat()
	xml.RegisterXmlFormat()

	cmd := command.New()

	args := os.Args[1:]

	_, _, err := cmd.Find(args)
	if err != nil && args[0] != "__complete" {
		// default command when nothing matches...
		newArgs := []string{"eval"}
		cmd.SetArgs(append(newArgs, os.Args[1:]...))

	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
