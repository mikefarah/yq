package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
  "io/ioutil"
  "os"
  "github.com/codegangsta/cli"
  "strings"
)

func main() {
  app := cli.NewApp()
  app.Name = "yaml"
  app.Usage = "command line tool for reading and writing yaml"
  app.Commands = []cli.Command{
    {
      Name:      "read",
      Aliases:     []string{"r"},
      Usage:     "read <filename> <path>\n\te.g.: yaml read sample.json a.b.c\n\t(default) reads a property from a given yaml file",
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

  var parsed_data map[interface{}]interface{}

  err := yaml.Unmarshal([]byte(raw_data), &parsed_data)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  var path = c.Args()[1]
  var paths = strings.Split(path, ".")

  fmt.Println(read(parsed_data, paths[0], paths[1:len(paths)]))
}

func read(context map[interface{}]interface{}, head string, tail []string) interface{} {
  value := context[head]
  // fmt.Println("read called")

  switch value.(type) {
    case bool, int, string, []interface{}:
      return value
    default: // recurse into map
      return read(value.(map[interface{}]interface{}), tail[0], tail[1:len(tail)])
  }
}
