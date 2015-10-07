package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var trimOutput = true
var writeInplace = false
var writeScript = ""

func main() {
	var cmdRead = &cobra.Command{
		Use:     "read [yaml_file] [path]",
		Aliases: []string{"r"},
		Short:   "yaml r sample.yaml a.b.c",
		Example: `
yaml read things.yaml a.b.c
yaml r - a.b.c (reads from stdin)
yaml r things.yaml a.*.c
yaml r things.yaml a.array[0].blah
yaml r things.yaml a.array[*].blah
			`,
		Long: "Outputs the value of the given path in the yaml file to STDOUT",
		Run:  readProperty,
	}

	var cmdWrite = &cobra.Command{
		Use:     "write [yaml_file] [path] [value]",
		Aliases: []string{"w"},
		Short:   "yaml w [--inplace/-i] sample.yaml a.b.c newValueForC",
		Example: `
yaml write things.yaml a.b.c cat
yaml write --inplace things.yaml a.b.c cat
yaml w -i things.yaml a.b.c cat
			`,
		Long: `Updates the yaml file w.r.t the given path and value.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.`,
		Run: writeProperty,
	}
	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdWrite.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")

	var rootCmd = &cobra.Command{Use: "yaml"}
	rootCmd.PersistentFlags().BoolVarP(&trimOutput, "trim", "t", true, "trim yaml output")
	rootCmd.AddCommand(cmdRead, cmdWrite)
	rootCmd.Execute()
}

func readProperty(cmd *cobra.Command, args []string) {
	var parsedData map[interface{}]interface{}

	readYaml(args[0], &parsedData)

	if len(args) == 1 {
		printYaml(parsedData)
		os.Exit(0)
	}

	var paths = parsePath(args[1])

	printYaml(readMap(parsedData, paths[0], paths[1:len(paths)]))
}

func writeProperty(cmd *cobra.Command, args []string) {
	var writeCommands map[string]interface{}
	if writeScript != "" {
		readYaml(writeScript, &writeCommands)
	} else if len(args) < 3 {
		die("Must provide <filename> <path_to_update> <value>")
	} else {
		writeCommands[args[1]] = parseValue(args[2])
	}

	var parsedData map[interface{}]interface{}
	readYaml(args[0], &parsedData)

	for path, value := range writeCommands {
		var paths = parsePath(path)
		write(parsedData, paths[0], paths[1:len(paths)], value)
	}

	if writeInplace {
		ioutil.WriteFile(args[0], []byte(yamlToString(parsedData)), 0644)
	} else {
		printYaml(parsedData)
	}
}
func parseValue(argument string) interface{} {
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
		die("error printing yaml: %v", err)
	}
	outStr := string(out)
	// trim the trailing new line as it's easier for a script to add
	// it in if required than to remove it
	if trimOutput {
		return strings.Trim(outStr, "\n ")
	}
	return outStr
}

func readYaml(filename string, parsedData interface{}) {
	if filename == "" {
		die("Must provide filename")
	}

	var rawData []byte
	if filename == "-" {
		rawData = readStdin()
	} else {
		rawData = readFile(filename)
	}

	err := yaml.Unmarshal([]byte(rawData), parsedData)
	if err != nil {
		die("error: %v", err)
	}
}

func readStdin() []byte {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		die("error reading stdin", err)
	}
	return bytes
}

func readFile(filename string) []byte {
	var rawData, readError = ioutil.ReadFile(filename)
	if readError != nil {
		die("error: %v", readError)
	}
	return rawData
}

func die(message ...interface{}) {
	fmt.Println(message)
	os.Exit(1)
}
