package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "yaml"
	app.Usage = "command line tool for reading and writing yaml"
	app.Commands = []cli.Command{
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "read <filename> <path>\n\te.g.: yaml read sample.json a.b.c\n\t(default) reads a property from a given yaml file",
			Action:  readProperty,
		},
	}
	app.Action = readProperty
	app.Run(os.Args)
}

func readProperty(c *cli.Context) {
	var parsedData map[interface{}]interface{}
	readYaml(c, &parsedData)

	var path = c.Args()[1]
	var paths = strings.Split(path, ".")

	printYaml(readMap(parsedData, paths[0], paths[1:len(paths)]))
}

func printYaml(context interface{}) {
  out, err := yaml.Marshal(context)
  if err != nil {
    log.Fatalf("error printing yaml: %v", err)
  }
  fmt.Println(string(out))
}

func readYaml(c *cli.Context, parsedData *map[interface{}]interface{}) {
	if len(c.Args()) == 0 {
		log.Fatalf("Must provide filename")
	}
	var rawData = readFile(c.Args()[0])

	if len(c.Args()) == 1 {
		fmt.Println(string(rawData[:]))
		os.Exit(0)
	}

	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func readFile(filename string) []byte {
	var rawData, readError = ioutil.ReadFile(filename)
	if readError != nil {
		log.Fatalf("error: %v", readError)
	}
	return rawData
}

func readMap(context map[interface{}]interface{}, head string, tail []string) interface{} {
	value := context[head]
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:len(tail)])
	}
	return value
}

func recurse(value interface{}, head string, tail []string) interface{} {
	switch value.(type) {
	case []interface{}:
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			log.Fatalf("Error accessing array: %v", err)
		}
		return readArray(value.([]interface{}), index, tail)
	default:
		return readMap(value.(map[interface{}]interface{}), head, tail)
	}
}

func readArray(array []interface{}, head int64, tail []string) interface{} {
	value := array[head]
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:len(tail)])
	}
	return value
}
