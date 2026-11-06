package yqlib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PrinterWriter interface {
	GetWriter(node *CandidateNode) (*bufio.Writer, error)
}

type singlePrinterWriter struct {
	bufferedWriter *bufio.Writer
}

func NewSinglePrinterWriter(writer io.Writer) PrinterWriter {
	return &singlePrinterWriter{
		bufferedWriter: bufio.NewWriter(writer),
	}
}

func (sp *singlePrinterWriter) GetWriter(_ *CandidateNode) (*bufio.Writer, error) {
	return sp.bufferedWriter, nil
}

type multiPrintWriter struct {
	treeNavigator  DataTreeNavigator
	nameExpression *ExpressionNode
	extension      string
	index          int
}

func NewMultiPrinterWriter(expression *ExpressionNode, format *Format) PrinterWriter {
	extension := "yml"

	switch format {
	case JSONFormat:
		extension = "json"
	case PropertiesFormat:
		extension = "properties"
	}

	return &multiPrintWriter{
		nameExpression: expression,
		extension:      extension,
		treeNavigator:  NewDataTreeNavigator(),
		index:          0,
	}
}

func (sp *multiPrintWriter) GetWriter(node *CandidateNode) (*bufio.Writer, error) {
	name := ""

	indexVariableNode := CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: fmt.Sprintf("%v", sp.index)}

	context := Context{MatchingNodes: node.AsList()}
	context.SetVariable("index", indexVariableNode.AsList())
	result, err := sp.treeNavigator.GetMatchingNodes(context, sp.nameExpression)
	if err != nil {
		return nil, err
	}
	if result.MatchingNodes.Len() > 0 {
		name = result.MatchingNodes.Front().Value.(*CandidateNode).Value
	}
	var extensionRegexp = regexp.MustCompile(`\.[a-zA-Z0-9]+$`)
	if !extensionRegexp.MatchString(name) {
		name = fmt.Sprintf("%v.%v", name, sp.extension)
	}

	var f *os.File
	if ConfiguredSecurityPreferences.AllowWriteAnyPath {
		err = os.MkdirAll(filepath.Dir(name), 0750)
		if err != nil {
			return nil, err
		}
		f, err = os.Create(name)
	} else {
		f, err = sp.createLocalFile(name)
	}

	if err != nil {
		return nil, err
	}
	sp.index = sp.index + 1

	return bufio.NewWriter(f), nil

}

func (sp *multiPrintWriter) createLocalFile(name string) (*os.File, error) {
	name = filepath.Clean(name)
	if !filepath.IsLocal(name) {
		return nil, sp.rejectPath(name, fmt.Errorf("path must stay within the current directory"))
	}

	root, err := os.OpenRoot(".")
	if err != nil {
		return nil, err
	}
	defer root.Close()

	dir := filepath.Dir(name)
	if err := rejectSymlinkDirectories(root, dir); err != nil {
		return nil, sp.rejectPath(name, err)
	}
	if err := root.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}
	if err := rejectSymlinkDirectories(root, dir); err != nil {
		return nil, sp.rejectPath(name, err)
	}

	info, err := root.Lstat(name)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		if err := root.Remove(name); err != nil {
			return nil, err
		}
		return root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	}
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
}

func rejectSymlinkDirectories(root *os.Root, dir string) error {
	current := ""
	for _, part := range strings.Split(dir, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := root.Lstat(current)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("directory %q is a symbolic link", current)
		}
		if !info.IsDir() {
			return fmt.Errorf("path component %q is not a directory", current)
		}
	}
	return nil
}

func (sp *multiPrintWriter) rejectPath(name string, reason error) error {
	message := fmt.Sprintf("split-exp output path %q was rejected: %v; use --security-allow-write-any-path to override", name, reason)
	GetLogger().Warning(message)
	return fmt.Errorf("%s", message)
}
