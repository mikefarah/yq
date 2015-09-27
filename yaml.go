package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
  "io/ioutil"
  "os"
  "github.com/codegangsta/cli"
  "strings"
  "strconv"
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
      Action:   read_property,
    },
  }
  app.Action = read_property
  app.Run(os.Args)
}

func read_property(c *cli.Context) {
  var parsed_data map[interface{}]interface{}
  read_yaml(c, &parsed_data)

  var path = c.Args()[1]
  var paths = strings.Split(path, ".")

  fmt.Println(read_map(parsed_data, paths[0], paths[1:len(paths)]))
}

func read_yaml(c *cli.Context, parsed_data *map[interface{}]interface{}) {
  if len(c.Args()) == 0 {
    log.Fatalf("Must provide filename")
  }
  var raw_data = read_file(c.Args()[0])

  if len(c.Args()) == 1 {
    fmt.Println(string(raw_data[:]))
    os.Exit(0)
  }

  err := yaml.Unmarshal([]byte(raw_data), &parsed_data)
  if err != nil {
    log.Fatalf("error: %v", err)
  }
}

func read_file(filename string) []byte {
  var raw_data, read_error = ioutil.ReadFile(filename)
  if read_error != nil {
    log.Fatalf("error: %v", read_error)
  }
  return raw_data
}

func read_map(context map[interface{}]interface{}, head string, tail []string) interface{} {
  value := context[head]
  if (len(tail) > 0) {
    return recurse(value, tail[0], tail[1:len(tail)])
  } else {
    return value
  }
}

func recurse(value interface{}, head string, tail []string) interface{} {
  switch value.(type) {
      case []interface {}:
        index, err := strconv.ParseInt(head, 10, 64)
        if err != nil {
          log.Fatalf("Error accessing array: %v", err)
        }
        return read_array(value.([]interface {}), index, tail)
      default:
        return read_map(value.(map[interface{}]interface{}), head, tail)
    }
}

func read_array(array []interface {}, head int64, tail[]string) interface{} {
  value := array[head]
  if (len(tail) > 0) {
    return recurse(value, tail[0], tail[1:len(tail)])
  } else {
    return value
  }
}

