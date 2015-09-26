package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
  "io/ioutil"
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "yaml"
  app.Usage = "command line tool for reading and writing yaml"
  app.Commands = []cli.Command{
    {
      Name:      "read",
      Aliases:     []string{"r"},
      Usage:     "(default) reads a property from a given yaml file",
      Action:   read_file,
    },
  }
  app.Action = read_file
  app.Run(os.Args)
}

func read_file(c *cli.Context) {
  if len(c.Args()) == 0 {
    log.Fatalf("Must provide filename")
  }
  var filename = c.Args()[0]
  var raw_data, read_error = ioutil.ReadFile(filename)

  if read_error != nil {
    log.Fatalf("error: %v", read_error)
  }

  var parsed_data interface{}

  err := yaml.Unmarshal([]byte(raw_data), &parsed_data)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  fmt.Println(parsed_data)
}
