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
			Usage:   "read <filename> <path>\n\te.g.: yaml read sample.json a.b.c\n\t(default) reads a property from a given yaml file\n",
			Action:  readProperty,
		},
		{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "write <filename> <path> <value>\n\te.g.: yaml write sample.json a.b.c 5\n\tupdates a property from a given yaml file, outputs to stdout\n",
			Action:  writeProperty,
		},
	}
	app.Action = readProperty
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "trim, t",
			Value: "true",
			Usage: "trim output",
		},
	}
	app.Run(os.Args)
}

func readProperty(c *cli.Context) {

	var parsedData map[interface{}]interface{}

	readYaml(c, &parsedData)

	if len(c.Args()) == 1 {
		printYaml(parsedData, c.Bool("trim"))
		os.Exit(0)
	}

	var paths = parsePath(c.Args()[1])

	printYaml(readMap(parsedData, paths[0], paths[1:len(paths)]), c.Bool("trim"))
}

func writeProperty(c *cli.Context) {
	var parsedData map[interface{}]interface{}
	readYaml(c, &parsedData)

	if len(c.Args()) < 3 {
		log.Fatalf("Must provide <filename> <path_to_update> <value>")
	}

	var paths = parsePath(c.Args()[1])

	write(parsedData, paths[0], paths[1:len(paths)], getValue(c.Args()[2]))

	printYaml(parsedData, c.Bool("trim"))
}

func getValue(argument string) interface{} {
	var value, err interface{}
	var inQuotes = argument[0] == '"'
	if !inQuotes {
		value, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			return value
		}
		value, err = strconv.ParseBool(argument)
		if err == nil {
			return value
		}
		return argument
	}
	return argument[1 : len(argument)-1]
}

func printYaml(context interface{}, trim bool) {
	out, err := yaml.Marshal(context)
	if err != nil {
		log.Fatalf("error printing yaml: %v", err)
	}
	outStr := string(out)
	if trim {
		outStr = strings.Trim(outStr, "\n ")
	}
	fmt.Println(outStr)
}

func readYaml(c *cli.Context, parsedData *map[interface{}]interface{}) {
	if len(c.Args()) == 0 {
		log.Fatalf("Must provide filename")
	}
	var rawData = readFile(c.Args()[0])

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
