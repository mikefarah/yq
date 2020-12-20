package cmd

// import (
// 	"strings"
// 	"testing"

// 	"github.com/mikefarah/yq/v3/test"
// 	"github.com/spf13/cobra"
// )

// func getRootCommand() *cobra.Command {
// 	return New()
// }

// func TestRootCmd(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	if !strings.Contains(result.Output, "Usage:") {
// 		t.Error("Expected usage message to be printed out, but the usage message was not found.")
// 	}
// }

// func TestRootCmd_Help(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "--help")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	if !strings.Contains(result.Output, "yq is a lightweight and portable command-line YAML processor. It aims to be the jq or sed of yaml files.") {
// 		t.Error("Expected usage message to be printed out, but the usage message was not found.")
// 	}
// }

// func TestRootCmd_VerboseLong(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "--verbose")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	if !verbose {
// 		t.Error("Expected verbose to be true")
// 	}
// }

// func TestRootCmd_VerboseShort(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "-v")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	if !verbose {
// 		t.Error("Expected verbose to be true")
// 	}
// }

// func TestRootCmd_VersionShort(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "-V")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	if !strings.Contains(result.Output, "yq version") {
// 		t.Error("expected version message to be printed out, but the message was not found.")
// 	}
// }

// func TestRootCmd_VersionLong(t *testing.T) {
// 	cmd := getRootCommand()
// 	result := test.RunCmd(cmd, "--version")
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}
// 	if !strings.Contains(result.Output, "yq version") {
// 		t.Error("expected version message to be printed out, but the message was not found.")
// 	}
// }
