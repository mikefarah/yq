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

	var path = c.Args()[1]
	var paths = strings.Split(path, ".")

	printYaml(readMap(parsedData, paths[0], paths[1:len(paths)]), c.Bool("trim"))
}

func writeProperty(c *cli.Context) {
	var parsedData map[interface{}]interface{}
	readYaml(c, &parsedData)

	if len(c.Args()) < 3 {
		log.Fatalf("Must provide <filename> <path_to_update> <value>")
	}

	var forceString bool
	if len(c.Args()) == 4 {
		forceString = true
	}

	var path = c.Args()[1]
	var paths = strings.Split(path, ".")

	write(parsedData, paths[0], paths[1:len(paths)], getValue(c.Args()[2], forceString))

	printYaml(parsedData, c.Bool("trim"))
}

func getValue(argument string, forceString bool) interface{} {
	var value, err interface{}

	if !forceString {
		value, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			return value
		}
		value, err = strconv.ParseBool(argument)
		if err == nil {
			return value
		}
	}
	return argument
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

func write(context map[interface{}]interface{}, head string, tail []string, value interface{}) {
	// e.g. if updating a.b.c, we need to get the 'b' map...
	toUpdate := readMap(context, head, tail[0:len(tail)-1]).(map[interface{}]interface{})
	//  and then set the 'c' key.
	key := (tail[len(tail)-1])
	toUpdate[key] = value
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
