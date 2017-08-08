package main

import (
	"fmt"
	"github.com/op/go-logging"
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
var inputJSON = false
var outputToJSON = false
var verbose = false
var log = logging.MustGetLogger("yaml")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
)
var backend = logging.AddModuleLevel(
	logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

func main() {
	backend.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend)

	var cmdRead = createReadCmd()
	var cmdWrite = createWriteCmd()
	var cmdNew = createNewCmd()

	var rootCmd = &cobra.Command{Use: "yaml"}
	rootCmd.PersistentFlags().BoolVarP(&trimOutput, "trim", "t", true, "trim yaml output")
	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "output as json")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.AddCommand(cmdRead, cmdWrite, cmdNew)
	rootCmd.Execute()
}

func createReadCmd() *cobra.Command {
	return &cobra.Command{
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
}

func createWriteCmd() *cobra.Command {
	var cmdWrite = &cobra.Command{
		Use:     "write [yaml_file] [path] [value]",
		Aliases: []string{"w"},
		Short:   "yaml w [--inplace/-i] [--script/-s script_file] sample.yaml a.b.c newValueForC",
		Example: `
yaml write things.yaml a.b.c cat
yaml write --inplace things.yaml a.b.c cat
yaml w -i things.yaml a.b.c cat
yaml w --script update_script.yaml things.yaml
yaml w -i -s update_script.yaml things.yaml
      `,
		Long: `Updates the yaml file w.r.t the given path and value.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

Update Scripts:
Note that you can give an update script to perform more sophisticated updated. Update script
format is a yaml map where the key is the path and the value is..well the value. e.g.:
---
a.b.c: true,
a.b.e:
  - name: bob
`,
		Run: writeProperty,
	}
	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdWrite.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	return cmdWrite
}

func createNewCmd() *cobra.Command {
	var cmdNew = &cobra.Command{
		Use:     "new [path] [value]",
		Aliases: []string{"n"},
		Short:   "yaml n [--script/-s script_file] a.b.c newValueForC",
		Example: `
yaml new a.b.c cat
yaml n a.b.c cat
yaml n --script create_script.yaml
      `,
		Long: `Creates a new yaml w.r.t the given path and value.
Outputs to STDOUT

Create Scripts:
Note that you can give a create script to perform more sophisticated yaml. This follows the same format as the update script.
`,
		Run: newProperty,
	}
	cmdNew.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	return cmdNew
}

func readProperty(cmd *cobra.Command, args []string) {
	if verbose {
		backend.SetLevel(logging.DEBUG, "")
	}
	print(read(args))
}

func read(args []string) interface{} {

	var parsedData yaml.MapSlice
	var path = ""
	if len(args) > 1 {
		path = args[1]
	}
	err := readData(args[0], &parsedData, inputJSON)
	if err != nil {
		var generalData interface{}
		readDataOrDie(args[0], &generalData, inputJSON)
		item := yaml.MapItem{Key: "thing", Value: generalData}
		parsedData = yaml.MapSlice{item}
		path = "thing." + path
	}

	if path == "" {
		return parsedData
	}

	var paths = parsePath(path)

	return readMap(parsedData, paths[0], paths[1:len(paths)])
}

func newProperty(cmd *cobra.Command, args []string) {
	if verbose {
		backend.SetLevel(logging.DEBUG, "")
	}
	updatedData := newYaml(args)
	print(updatedData)
}

func newYaml(args []string) interface{} {
	var writeCommands yaml.MapSlice
	if writeScript != "" {
		readDataOrDie(writeScript, &writeCommands, false)
	} else if len(args) < 2 {
		die("Must provide <path_to_update> <value>")
	} else {
		writeCommands = make(yaml.MapSlice, 1)
		writeCommands[0] = yaml.MapItem{Key: args[0], Value: parseValue(args[1])}
	}

	var parsedData yaml.MapSlice
	var prependCommand = ""
	var isArray = strings.HasPrefix(writeCommands[0].Key.(string), "[")
	if isArray {
		item := yaml.MapItem{Key: "thing", Value: make(yaml.MapSlice, 0)}
		parsedData = yaml.MapSlice{item}
		prependCommand = "thing"
	} else {
		parsedData = make(yaml.MapSlice, 0)
	}

	return updateParsedData(parsedData, writeCommands, prependCommand)
}

func writeProperty(cmd *cobra.Command, args []string) {
	if verbose {
		backend.SetLevel(logging.DEBUG, "")
	}
	updatedData := updateYaml(args)
	if writeInplace {
		ioutil.WriteFile(args[0], []byte(yamlToString(updatedData)), 0644)
	} else {
		print(updatedData)
	}
}

func updateParsedData(parsedData yaml.MapSlice, writeCommands yaml.MapSlice, prependCommand string) interface{} {
	var prefix = ""
	if prependCommand != "" {
		prefix = prependCommand + "."
	}
	for _, entry := range writeCommands {
		path := prefix + entry.Key.(string)
		value := entry.Value
		var paths = parsePath(path)
		parsedData = writeMap(parsedData, paths, value)
	}
	if prependCommand != "" {
		return readMap(parsedData, prependCommand, make([]string, 0))
	}
	return parsedData
}

func updateYaml(args []string) interface{} {
	var writeCommands yaml.MapSlice
	var prependCommand = ""
	if writeScript != "" {
		readDataOrDie(writeScript, &writeCommands, false)
	} else if len(args) < 3 {
		die("Must provide <filename> <path_to_update> <value>")
	} else {
		writeCommands = make(yaml.MapSlice, 1)
		writeCommands[0] = yaml.MapItem{Key: args[1], Value: parseValue(args[2])}
	}

	var parsedData yaml.MapSlice
	err := readData(args[0], &parsedData, inputJSON)
	if err != nil {
		var generalData interface{}
		readDataOrDie(args[0], &generalData, inputJSON)
		item := yaml.MapItem{Key: "thing", Value: generalData}
		parsedData = yaml.MapSlice{item}
		prependCommand = "thing"
	}

	return updateParsedData(parsedData, writeCommands, prependCommand)
}

func parseValue(argument string) interface{} {
	var value, err interface{}
	var inQuotes = len(argument) > 0 && argument[0] == '"'
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

func print(context interface{}) {
	var out string
	if outputToJSON {
		out = jsonToString(context)
	} else {
		out = yamlToString(context)
	}
	fmt.Println(out)
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

func readDataOrDie(filename string, parsedData interface{}, readAsJSON bool) {
	err := readData(filename, parsedData, readAsJSON)
	if err != nil {
		die("error parsing data: ", err)
	}
}

func readData(filename string, parsedData interface{}, readAsJSON bool) error {
	if filename == "" {
		die("Must provide filename")
	}

	var rawData []byte
	if filename == "-" {
		rawData = readStdin()
	} else {
		rawData = readFile(filename)
	}

	return yaml.Unmarshal([]byte(rawData), parsedData)
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
