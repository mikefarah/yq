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
			Usage:   "read <filename> <path>\n\te.g.: yaml read sample.yaml a.b.c\n\t(default) reads a property from a given yaml file\n",
			Action:  readProperty,
		},
		{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "write <filename> <path> <value>\n\te.g.: yaml write sample.yaml a.b.c 5\n\tupdates a property from a given yaml file, outputs to stdout\n",
			Action:  writeProperty,
		},
		{
			Name:    "write-inplace",
			Aliases: []string{"wi"},
			Usage:   "wi <filename> <path> <value>\n\te.g.: yaml wi sample.yaml a.b.c 5\n\tupdates a property from a given yaml file and saves it to the given filename (sample.yaml)\n",
			Action:  writePropertyInPlace,
		},
	}
	app.Action = readProperty
	app.Run(os.Args)
}

func readProperty(c *cli.Context) {

	var parsedData map[interface{}]interface{}

	readYaml(c, &parsedData)

	if len(c.Args()) == 1 {
		printYaml(parsedData)
		os.Exit(0)
	}

	var paths = parsePath(c.Args()[1])

	printYaml(readMap(parsedData, paths[0], paths[1:len(paths)]))
}

func writeProperty(c *cli.Context) {
	printYaml(updateProperty(c))
}

func writePropertyInPlace(c *cli.Context) {
	updatedYaml := updateProperty(c)
	ioutil.WriteFile(c.Args()[0], []byte(updatedYaml), 0644)
}

func updateProperty(c *cli.Context) string {
	var parsedData map[interface{}]interface{}
	readYaml(c, &parsedData)

	if len(c.Args()) < 3 {
		log.Fatalf("Must provide <filename> <path_to_update> <value>")
	}

	var paths = parsePath(c.Args()[1])

	write(parsedData, paths[0], paths[1:len(paths)], getValue(c.Args()[2]))

	return yamlToString(parsedData)
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

func printYaml(context interface{}) {
	fmt.Println(yamlToString(context))
}

func yamlToString(context interface{}) string {
	out, err := yaml.Marshal(context)
	if err != nil {
		log.Fatalf("error printing yaml: %v", err)
	}
	outStr := string(out)
	// trim the trailing new line as it's easier for a script to add
	// it in if required than to remove it
	return strings.Trim(outStr, "\n ")
}

func readYaml(c *cli.Context, parsedData *map[interface{}]interface{}) {
	if len(c.Args()) == 0 {
		log.Fatalf("Must provide filename")
	}

	var rawData []byte
	fmt.Println("c.Args()[0]", c.Args()[0])
	if( c.Args()[0] == "-") {
		rawData = readStdin()
	} else {
		rawData = readFile(c.Args()[0])
	}

	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func readStdin() []byte {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading stdin", err)
	}
	return bytes
}

func readFile(filename string) []byte {
	var rawData, readError = ioutil.ReadFile(filename)
	if readError != nil {
		log.Fatalf("error: %v", readError)
	}
	return rawData
}
