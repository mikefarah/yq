package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var trimOutput = true
var writeInplace = false
var writeScript = ""
var outputToJSON = false
var overwriteFlag = false
var verbose = false
var version = false
var log = logging.MustGetLogger("yq")

func main() {
	cmd := newCommandCLI()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func newCommandCLI() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "yq",
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			cmd.Println(cmd.UsageString())

			return nil
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var format = logging.MustStringFormatter(
				`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
			)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			if verbose {
				backend.SetLevel(logging.DEBUG, "")
			} else {
				backend.SetLevel(logging.ERROR, "")
			}

			logging.SetBackend(backend)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&trimOutput, "trim", "t", true, "trim yaml output")
	rootCmd.PersistentFlags().BoolVarP(&outputToJSON, "tojson", "j", false, "output as json")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")

	rootCmd.AddCommand(createReadCmd(), createWriteCmd(), createNewCmd(), createMergeCmd())
	rootCmd.SetOutput(os.Stdout)

	return rootCmd
}

func createReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "read [yaml_file] [path]",
		Aliases: []string{"r"},
		Short:   "yq r sample.yaml a.b.c",
		Example: `
yq read things.yaml a.b.c
yq r - a.b.c (reads from stdin)
yq r things.yaml a.*.c
yq r things.yaml a.array[0].blah
yq r things.yaml a.array[*].blah
      `,
		Long: "Outputs the value of the given path in the yaml file to STDOUT",
		RunE: readProperty,
	}
}

func createWriteCmd() *cobra.Command {
	var cmdWrite = &cobra.Command{
		Use:     "write [yaml_file] [path] [value]",
		Aliases: []string{"w"},
		Short:   "yq w [--inplace/-i] [--script/-s script_file] sample.yaml a.b.c newValueForC",
		Example: `
yq write things.yaml a.b.c cat
yq write --inplace things.yaml a.b.c cat
yq w -i things.yaml a.b.c cat
yq w --script update_script.yaml things.yaml
yq w -i -s update_script.yaml things.yaml
yq w things.yaml a.b.d[+] foo
yq w things.yaml a.b.d[+] foo
      `,
		Long: `Updates the yaml file w.r.t the given path and value.
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

Append value to array adds the value to the end of array.

Update Scripts:
Note that you can give an update script to perform more sophisticated updated. Update script
format is a yaml map where the key is the path and the value is..well the value. e.g.:
---
a.b.c: true,
a.b.e:
  - name: bob
`,
		RunE: writeProperty,
	}
	cmdWrite.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdWrite.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	return cmdWrite
}

func createNewCmd() *cobra.Command {
	var cmdNew = &cobra.Command{
		Use:     "new [path] [value]",
		Aliases: []string{"n"},
		Short:   "yq n [--script/-s script_file] a.b.c newValueForC",
		Example: `
yq new a.b.c cat
yq n a.b.c cat
yq n --script create_script.yaml
      `,
		Long: `Creates a new yaml w.r.t the given path and value.
Outputs to STDOUT

Create Scripts:
Note that you can give a create script to perform more sophisticated yaml. This follows the same format as the update script.
`,
		RunE: newProperty,
	}
	cmdNew.PersistentFlags().StringVarP(&writeScript, "script", "s", "", "yaml script for updating yaml")
	return cmdNew
}

func createMergeCmd() *cobra.Command {
	var cmdMerge = &cobra.Command{
		Use:     "merge [initial_yaml_file] [additional_yaml_file]...",
		Aliases: []string{"m"},
		Short:   "yq m [--inplace/-i] [--overwrite/-x] sample.yaml sample2.yaml",
		Example: `
yq merge things.yaml other.yaml
yq merge --inplace things.yaml other.yaml
yq m -i things.yaml other.yaml
yq m --overwrite things.yaml other.yaml
yq m -i -x things.yaml other.yaml
      `,
		Long: `Updates the yaml file by adding/updating the path(s) and value(s) from additional yaml file(s).
Outputs to STDOUT unless the inplace flag is used, in which case the file is updated instead.

If overwrite flag is set then existing values will be overwritten using the values from each additional yaml file.
`,
		RunE: mergeProperties,
	}
	cmdMerge.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace")
	cmdMerge.PersistentFlags().BoolVarP(&overwriteFlag, "overwrite", "x", false, "update the yaml file by overwriting existing values")
	return cmdMerge
}

func readProperty(cmd *cobra.Command, args []string) error {
	data, err := read(args)
	if err != nil {
		return err
	}
	dataStr, err := toString(data)
	if err != nil {
		return err
	}
	cmd.Println(dataStr)
	return nil
}

var regexpSubdocPath = regexp.MustCompile(`^@(\d+)\.(.*)$`)

func read(args []string) (interface{}, error) {
	var (
		parsedData    yaml.MapSlice
		path          = ""
		multidoc      bool
		docNumber     int
		matches       []string
		err           error
		childDocument string
	)

	if len(args) < 1 {
		return nil, errors.New("Must provide filename")
	} else if len(args) > 1 {
		path = args[1]
		multidoc = regexpSubdocPath.MatchString(path)
	}

	if multidoc {
		matches = regexpSubdocPath.FindStringSubmatch(path)
		// matches[0] is the original path
		// matches[1] is the document number
		// matches[2] is the path within that document
		if len(matches) == 3 {
			docNumber, err = strconv.Atoi(matches[1])
			if err == nil {
				path = matches[2]
			}
		}
		childDocument, err = readSubDocument(args[0], docNumber)
		if err != nil {
			return nil, err
		}
		tmp, err := ioutil.TempFile("", "")
		defer os.Remove(tmp.Name())
		if err != nil {
			return nil, err
		}
		n, err := tmp.WriteString(childDocument)
		if err != nil {
			return nil, err
		}
		if n != len(childDocument) {
			return nil, fmt.Errorf("multidoc tmp file expected to write %v bytes but wrote %v bytes", len(childDocument), n)
		}

		args[0] = tmp.Name()

	}

	// If we're asking for any document beyond the first,
	// we need to provide the yaml package with that document

	if err := readData(args[0], &parsedData); err != nil {
		var generalData interface{}
		if err = readData(args[0], &generalData); err != nil {
			return nil, err
		}
		item := yaml.MapItem{Key: "thing", Value: generalData}
		parsedData = yaml.MapSlice{item}
		path = "thing." + path
	}

	if parsedData != nil && parsedData[0].Key == nil {
		var interfaceData []map[interface{}]interface{}
		if err := readData(args[0], &interfaceData); err == nil {
			var listMap []yaml.MapSlice
			for _, item := range interfaceData {
				listMap = append(listMap, mapToMapSlice(item))
			}
			return readYamlArray(listMap, path)
		}
	}

	if path == "" {
		return parsedData, nil
	}

	var paths = parsePath(path)

	return readMap(parsedData, paths[0], paths[1:])
}

func readYamlArray(listMap []yaml.MapSlice, path string) (interface{}, error) {
	if path == "" {
		return listMap, nil
	}

	var paths = parsePath(path)

	if paths[0] == "*" {
		if len(paths[1:]) == 0 {
			return listMap, nil
		}
		var results []interface{}
		for _, m := range listMap {
			value, err := readMap(m, paths[1], paths[2:])
			if err != nil {
				return nil, err
			}
			results = append(results, value)
		}
		return results, nil
	}

	index, err := strconv.ParseInt(paths[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error accessing array: %v", err)
	}
	if len(paths[1:]) == 0 {
		return listMap[index], nil
	}
	return readMap(listMap[index], paths[1], paths[2:])
}

func newProperty(cmd *cobra.Command, args []string) error {
	updatedData, err := newYaml(args)
	if err != nil {
		return err
	}
	dataStr, err := toString(updatedData)
	if err != nil {
		return err
	}
	cmd.Println(dataStr)
	return nil
}

func newYaml(args []string) (interface{}, error) {
	var writeCommands yaml.MapSlice
	if writeScript != "" {
		if err := readData(writeScript, &writeCommands); err != nil {
			return nil, err
		}
	} else if len(args) < 2 {
		return nil, errors.New("Must provide <path_to_update> <value>")
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

func writeProperty(cmd *cobra.Command, args []string) error {
	updatedData, err := updateYaml(args)
	if err != nil {
		return err
	}
	return write(cmd, args[0], updatedData)
}

func write(cmd *cobra.Command, filename string, updatedData interface{}) error {
	if writeInplace {
		dataStr, err := yamlToString(updatedData)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filename, []byte(dataStr), 0644)
	}
	dataStr, err := toString(updatedData)
	if err != nil {
		return err
	}
	cmd.Println(dataStr)
	return nil
}

func mergeProperties(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Must provide at least 2 yaml files")
	}

	updatedData, err := mergeYaml(args)
	if err != nil {
		return err
	}
	return write(cmd, args[0], updatedData)
}

func mergeYaml(args []string) (interface{}, error) {
	var updatedData map[interface{}]interface{}

	for _, f := range args {
		var parsedData map[interface{}]interface{}
		if err := readData(f, &parsedData); err != nil {
			return nil, err
		}
		if err := merge(&updatedData, parsedData, overwriteFlag); err != nil {
			return nil, err
		}
	}

	return mapToMapSlice(updatedData), nil
}

func updateParsedData(parsedData yaml.MapSlice, writeCommands yaml.MapSlice, prependCommand string) (interface{}, error) {
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
	return parsedData, nil
}

func updateYaml(args []string) (interface{}, error) {
	var writeCommands yaml.MapSlice
	var prependCommand = ""
	if writeScript != "" {
		if err := readData(writeScript, &writeCommands); err != nil {
			return nil, err
		}
	} else if len(args) < 3 {
		return nil, errors.New("Must provide <filename> <path_to_update> <value>")
	} else {
		writeCommands = make(yaml.MapSlice, 1)
		writeCommands[0] = yaml.MapItem{Key: args[1], Value: parseValue(args[2])}
	}

	var parsedData yaml.MapSlice
	if err := readData(args[0], &parsedData); err != nil {
		var generalData interface{}
		if err = readData(args[0], &generalData); err != nil {
			return nil, err
		}
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
		if argument == "[]" {
			return make([]interface{}, 0)
		}
		return argument
	}
	return argument[1 : len(argument)-1]
}

func toString(context interface{}) (string, error) {
	if outputToJSON {
		return jsonToString(context)
	}
	return yamlToString(context)
}

func yamlToString(context interface{}) (string, error) {
	switch context.(type) {
	case string:
		return context.(string), nil
	default:
		return marshalContext(context)
	}
}

func marshalContext(context interface{}) (string, error) {
	out, err := yaml.Marshal(context)

	if err != nil {
		return "", fmt.Errorf("error printing yaml: %v", err)
	}

	outStr := string(out)
	// trim the trailing new line as it's easier for a script to add
	// it in if required than to remove it
	if trimOutput {
		return strings.Trim(outStr, "\n "), nil
	}
	return outStr, nil
}

func readData(filename string, parsedData interface{}) error {
	if filename == "" {
		return errors.New("Must provide filename")
	}

	var rawData []byte
	var err error
	if filename == "-" {
		rawData, err = ioutil.ReadAll(os.Stdin)
	} else {
		rawData, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		return err
	}

	return yaml.Unmarshal(rawData, parsedData)
}

var regexpDocumentSeparator = regexp.MustCompile(`(?m)^---$`)

// readSubDocument takes a filename and an index.
// It looks for ^---$ markers that separate yaml documents and splits the file on those.
// It takes the zero-indexed nth document, writes it to a temporary file, and return a handle to that file.
// Is trailing whitespace legal yaml?
func readSubDocument(filename string, n int) (string, error) {
	allDocs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	subDocs := regexpDocumentSeparator.Split(string(allDocs), -1)
	if len(subDocs) == 0 {
		return "", errors.New("document was split into zero subdocuments")
	}
	if subDocs[0] == "" {
		subDocs = subDocs[1:]
	}
	if len(subDocs) <= n {
		return "", fmt.Errorf("requested document %v out of total %v documents", n, len(subDocs))
	}
	return strings.TrimSpace(subDocs[n]), nil
}
