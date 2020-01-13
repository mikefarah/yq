package main

import (
	command "github.com/mikefarah/yq/v3/cmd"
	logging "gopkg.in/op/go-logging.v1"
	"os"
)

func main() {
	cmd := command.New()
	log := logging.MustGetLogger("yq")
	if err := cmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
