package yqlib

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
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
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if filepath.IsAbs(name) {
		absoluteName, err := filepath.Abs(name)
		if err != nil {
			return nil, err
		}
		name, err = filepath.Rel(workingDirectory, absoluteName)
		if err != nil {
			return nil, sp.rejectPath(name, err)
		}
	}
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
	// Recheck after creating directories in case a path component resolved through a symlink.
	if err := rejectSymlinkDirectories(root, dir); err != nil {
		return nil, sp.rejectPath(name, err)
	}

	info, err := root.Lstat(name)
	mode := os.FileMode(0666)
	preservePermissions := err == nil && info.Mode().IsRegular()
	if preservePermissions {
		mode = info.Mode().Perm()
	} else if err == nil && info.Mode()&os.ModeSymlink == 0 {
		return nil, sp.rejectPath(name, fmt.Errorf("destination is not a regular file or symbolic link"))
	}
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	var source *os.File
	if preservePermissions {
		source, err = openSplitSource(root, name)
		if err != nil {
			return nil, sp.rejectPath(name, fmt.Errorf("could not preserve destination access-control metadata: %w", err))
		}
		defer source.Close()
		sourceInfo, err := source.Stat()
		if err != nil {
			return nil, sp.rejectPath(name, fmt.Errorf("could not preserve destination access-control metadata: %w", err))
		}
		if !sourceInfo.Mode().IsRegular() || !os.SameFile(info, sourceInfo) {
			return nil, sp.rejectPath(name, fmt.Errorf("destination changed while preserving access-control metadata"))
		}
		mode = sourceInfo.Mode().Perm()
	}

	// Atomically replace the destination with a new inode instead of following symlinks or truncating hard links.
	for range 10 {
		randomName := make([]byte, 16)
		if _, err := rand.Read(randomName); err != nil {
			return nil, err
		}
		temporaryName := filepath.Join(filepath.Dir(name), ".yq-split-"+hex.EncodeToString(randomName))
		file, err := root.OpenFile(temporaryName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if source != nil {
			if err := copySplitAccessControl(source, file); err != nil {
				_ = file.Close()
				_ = root.Remove(temporaryName)
				return nil, sp.rejectPath(name, fmt.Errorf("could not preserve destination access-control metadata: %w", err))
			}
		}
		if err := replaceSplitFile(root, file, temporaryName, name, mode, preservePermissions); err != nil {
			_ = file.Close()
			_ = root.Remove(temporaryName)
			return nil, err
		}
		return file, nil
	}
	return nil, fmt.Errorf("could not create a temporary split-exp output file")
}

func replaceSplitFile(root *os.Root, file *os.File, temporaryName, name string, mode os.FileMode, preservePermissions bool) error {
	if preservePermissions {
		// The open file handle remains bound to the new inode even if its path is replaced.
		if err := file.Chmod(mode); err != nil {
			return err
		}
	}
	return root.Rename(temporaryName, name)
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
