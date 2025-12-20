package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
)

func TestIsAutomaticOutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"empty format", "", true},
		{"auto format", "auto", true},
		{"short auto format", "a", true},
		{"json format", "json", false},
		{"yaml format", "yaml", false},
		{"xml format", "xml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalFormat := outputFormat
			defer func() { outputFormat = originalFormat }()

			outputFormat = tt.format
			result := isAutomaticOutputFormat()
			if result != tt.expected {
				t.Errorf("isAutomaticOutputFormat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMaybeFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", tempFile.Name(), true},
		{"existing directory", tempDir, false},
		{"non-existent path", "/path/that/does/not/exist", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maybeFile(tt.path)
			if result != tt.expected {
				t.Errorf("maybeFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestProcessArgs(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary .yq file for testing
	tempYqFile, err := os.Create("test.yq")
	if err != nil {
		t.Fatalf("Failed to create temp yq file: %v", err)
	}
	defer os.Remove(tempYqFile.Name())
	if _, err = tempYqFile.WriteString(".a.b"); err != nil {
		t.Fatalf("Failed to write to temp yq file: %v", err)
	}
	tempYqFile.Close()

	tests := []struct {
		name            string
		args            []string
		forceExpression string
		expressionFile  string
		expectedExpr    string
		expectedArgs    []string
		expectError     bool
	}{
		{
			name:            "empty args",
			args:            []string{},
			forceExpression: "",
			expressionFile:  "",
			expectedExpr:    "",
			expectedArgs:    []string{},
			expectError:     false,
		},
		{
			name:            "force expression",
			args:            []string{"file1"},
			forceExpression: ".a.b",
			expressionFile:  "",
			expectedExpr:    ".a.b",
			expectedArgs:    []string{"file1"},
			expectError:     false,
		},
		{
			name:            "expression as first arg",
			args:            []string{".a.b", "file1"},
			forceExpression: "",
			expressionFile:  "",
			expectedExpr:    ".a.b",
			expectedArgs:    []string{"file1"},
			expectError:     false,
		},
		{
			name:            "file as first arg",
			args:            []string{tempFile.Name()},
			forceExpression: "",
			expressionFile:  "",
			expectedExpr:    "",
			expectedArgs:    []string{tempFile.Name()},
			expectError:     false,
		},
		{
			name:            "yq file as first arg",
			args:            []string{tempYqFile.Name(), "things"},
			forceExpression: "",
			expressionFile:  "",
			expectedExpr:    ".a.b",
			expectedArgs:    []string{"things"},
			expectError:     false,
		},
		{
			name:            "dash as first arg",
			args:            []string{"-"},
			forceExpression: "",
			expressionFile:  "",
			expectedExpr:    "",
			expectedArgs:    []string{"-"},
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalForceExpression := forceExpression
			originalExpressionFile := expressionFile
			defer func() {
				forceExpression = originalForceExpression
				expressionFile = originalExpressionFile
			}()

			forceExpression = tt.forceExpression
			expressionFile = tt.expressionFile

			expr, args, err := processArgs(tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("processArgs() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("processArgs() unexpected error: %v", err)
				return
			}

			if expr != tt.expectedExpr {
				t.Errorf("processArgs() expression = %v, want %v", expr, tt.expectedExpr)
			}

			if !stringsEqual(args, tt.expectedArgs) {
				t.Errorf("processArgs() args = %v, want %v", args, tt.expectedArgs)
			}
		})
	}
}

func TestConfigureDecoder(t *testing.T) {
	tests := []struct {
		name             string
		inputFormat      string
		evaluateTogether bool
		expectError      bool
		expectType       string
	}{
		{
			name:             "yaml format",
			inputFormat:      "yaml",
			evaluateTogether: false,
			expectError:      false,
			expectType:       "yamlDecoder",
		},
		{
			name:             "json format",
			inputFormat:      "json",
			evaluateTogether: true,
			expectError:      false,
			expectType:       "jsonDecoder",
		},
		{
			name:             "xml format",
			inputFormat:      "xml",
			evaluateTogether: false,
			expectError:      false,
			expectType:       "xmlDecoder",
		},
		{
			name:             "invalid format",
			inputFormat:      "invalid",
			evaluateTogether: false,
			expectError:      true,
			expectType:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalInputFormat := inputFormat
			defer func() { inputFormat = originalInputFormat }()

			inputFormat = tt.inputFormat

			decoder, err := configureDecoder(tt.evaluateTogether)
			if tt.expectError {
				if err == nil {
					t.Errorf("configureDecoder() expected error but got none")
				}
				if decoder != nil {
					t.Errorf("configureDecoder() expected nil decoder but got %v", decoder)
				}
				return
			}

			if err != nil {
				t.Errorf("configureDecoder() unexpected error: %v", err)
				return
			}

			if decoder == nil {
				t.Errorf("configureDecoder() expected decoder but got nil")
				return
			}

			typeStr := fmt.Sprintf("%T", decoder)
			if !strings.Contains(typeStr, tt.expectType) {
				t.Errorf("configureDecoder() expected type to contain %q but got %q", tt.expectType, typeStr)
			}
		})
	}
}

func TestConfigurePrinterWriter(t *testing.T) {
	yqlib.InitExpressionParser()

	tests := []struct {
		name                string
		splitFileExp        string
		format              *yqlib.Format
		forceColor          bool
		expectError         bool
		expectMulti         bool
		expectColorsEnabled bool
	}{
		{
			name:                "single printer writer",
			splitFileExp:        "",
			format:              &yqlib.Format{},
			forceColor:          false,
			expectError:         false,
			expectMulti:         false,
			expectColorsEnabled: false,
		},
		{
			name:                "multi printer writer with valid expression",
			splitFileExp:        ".a.b",
			format:              &yqlib.Format{},
			forceColor:          true,
			expectError:         false,
			expectMulti:         true,
			expectColorsEnabled: true,
		},
		{
			name:                "multi printer writer with invalid expression",
			splitFileExp:        "[invalid",
			format:              &yqlib.Format{},
			forceColor:          false,
			expectError:         true,
			expectMulti:         false,
			expectColorsEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalSplitFileExp := splitFileExp
			originalForceColor := forceColor
			originalColorsEnabled := colorsEnabled
			defer func() {
				splitFileExp = originalSplitFileExp
				forceColor = originalForceColor
				colorsEnabled = originalColorsEnabled
			}()

			splitFileExp = tt.splitFileExp
			forceColor = tt.forceColor
			colorsEnabled = false // Reset to test the setting

			writer, err := configurePrinterWriter(tt.format, os.Stdout)
			if tt.expectError {
				if err == nil {
					t.Errorf("configurePrinterWriter() expected error but got none")
				}
				if writer != nil {
					t.Errorf("configurePrinterWriter() expected nil writer but got %v", writer)
				}
				return
			}

			if err != nil {
				t.Errorf("configurePrinterWriter() unexpected error: %v", err)
				return
			}

			if writer == nil {
				t.Errorf("configurePrinterWriter() expected writer but got nil")
				return
			}

			// Explicitly check colorsEnabled
			if colorsEnabled != tt.expectColorsEnabled {
				t.Errorf("configurePrinterWriter() colorsEnabled = %v, want %v", colorsEnabled, tt.expectColorsEnabled)
			}

			// Check the type of the returned writer
			writerType := fmt.Sprintf("%T", writer)
			if tt.expectMulti {
				if !strings.Contains(writerType, "multiPrintWriter") {
					t.Errorf("configurePrinterWriter() expected multiPrintWriter but got %s", writerType)
				}
			} else {
				if !strings.Contains(writerType, "singlePrinterWriter") {
					t.Errorf("configurePrinterWriter() expected singlePrinterWriter but got %s", writerType)
				}
			}
		})
	}
}

func TestConfigureEncoder(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		expectError  bool
		expectType   string
	}{
		{
			name:         "yaml format",
			outputFormat: "yaml",
			expectError:  false,
			expectType:   "yamlEncoder",
		},
		{
			name:         "json format",
			outputFormat: "json",
			expectError:  false,
			expectType:   "jsonEncoder",
		},
		{
			name:         "xml format",
			outputFormat: "xml",
			expectError:  false,
			expectType:   "xmlEncoder",
		},
		{
			name:         "properties format",
			outputFormat: "properties",
			expectError:  false,
			expectType:   "propertiesEncoder",
		},
		{
			name:         "invalid format",
			outputFormat: "invalid",
			expectError:  true,
			expectType:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalOutputFormat := outputFormat
			originalIndent := indent
			originalUnwrapScalar := unwrapScalar
			originalColorsEnabled := colorsEnabled
			originalNoDocSeparators := noDocSeparators
			defer func() {
				outputFormat = originalOutputFormat
				indent = originalIndent
				unwrapScalar = originalUnwrapScalar
				colorsEnabled = originalColorsEnabled
				noDocSeparators = originalNoDocSeparators
			}()

			outputFormat = tt.outputFormat
			indent = 2
			unwrapScalar = false
			colorsEnabled = false
			noDocSeparators = false

			encoder, err := configureEncoder()
			if tt.expectError {
				if err == nil {
					t.Errorf("configureEncoder() expected error but got none")
				}
				if encoder != nil {
					t.Errorf("configureEncoder() expected nil encoder but got %v", encoder)
				}
				return
			}

			if err != nil {
				t.Errorf("configureEncoder() unexpected error: %v", err)
				return
			}

			if encoder == nil {
				t.Errorf("configureEncoder() expected encoder but got nil")
				return
			}

			typeStr := fmt.Sprintf("%T", encoder)
			if !strings.Contains(typeStr, tt.expectType) {
				t.Errorf("configureEncoder() expected type to contain %q but got %q", tt.expectType, typeStr)
			}
		})
	}
}

func TestInitCommand(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary split file
	tempSplitFile, err := os.CreateTemp("", "split")
	if err != nil {
		t.Fatalf("Failed to create temp split file: %v", err)
	}
	defer os.Remove(tempSplitFile.Name())
	if _, err = tempSplitFile.WriteString(".a.b"); err != nil {
		t.Fatalf("Failed to write to temp split file: %v", err)
	}
	tempSplitFile.Close()

	tests := []struct {
		name             string
		args             []string
		writeInplace     bool
		frontMatter      string
		nullInput        bool
		splitFileExpFile string
		splitFileExp     string
		outputToJSON     bool
		expectError      bool
		errorContains    string
		expectExpr       string
		expectArgs       []string
	}{
		{
			name:         "basic command",
			args:         []string{tempFile.Name()},
			writeInplace: false,
			frontMatter:  "",
			nullInput:    false,
			expectError:  false,
			expectExpr:   "",
			expectArgs:   []string{tempFile.Name()},
		},
		{
			name:          "write inplace with no args",
			args:          []string{},
			writeInplace:  true,
			frontMatter:   "",
			nullInput:     false,
			expectError:   true,
			errorContains: "write in place flag only applicable when giving an expression and at least one file",
		},
		{
			name:             "split file expression from file",
			args:             []string{tempFile.Name()},
			writeInplace:     false,
			frontMatter:      "",
			nullInput:        false,
			splitFileExpFile: tempSplitFile.Name(),
			expectError:      false,
			expectExpr:       "",
			expectArgs:       []string{tempFile.Name()},
		},
		{
			name:         "output to JSON",
			args:         []string{tempFile.Name()},
			writeInplace: false,
			frontMatter:  "",
			nullInput:    false,
			outputToJSON: true,
			expectError:  false,
			expectExpr:   "",
			expectArgs:   []string{tempFile.Name()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalWriteInplace := writeInplace
			originalFrontMatter := frontMatter
			originalNullInput := nullInput
			originalSplitFileExpFile := splitFileExpFile
			originalSplitFileExp := splitFileExp
			originalOutputToJSON := outputToJSON
			originalInputFormat := inputFormat
			originalOutputFormat := outputFormat
			originalForceColor := forceColor
			originalForceNoColor := forceNoColor
			originalColorsEnabled := colorsEnabled
			defer func() {
				writeInplace = originalWriteInplace
				frontMatter = originalFrontMatter
				nullInput = originalNullInput
				splitFileExpFile = originalSplitFileExpFile
				splitFileExp = originalSplitFileExp
				outputToJSON = originalOutputToJSON
				inputFormat = originalInputFormat
				outputFormat = originalOutputFormat
				forceColor = originalForceColor
				forceNoColor = originalForceNoColor
				colorsEnabled = originalColorsEnabled
			}()

			writeInplace = tt.writeInplace
			frontMatter = tt.frontMatter
			nullInput = tt.nullInput
			splitFileExpFile = tt.splitFileExpFile
			splitFileExp = tt.splitFileExp
			outputToJSON = tt.outputToJSON
			inputFormat = "auto"
			outputFormat = "auto"
			forceColor = false
			forceNoColor = false
			colorsEnabled = false

			cmd := &cobra.Command{}
			expr, args, err := initCommand(cmd, tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("initCommand() expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("initCommand() error '%v' does not contain '%v'", err.Error(), tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("initCommand() unexpected error: %v", err)
				return
			}

			if expr != tt.expectExpr {
				t.Errorf("initCommand() expr = %v, want %v", expr, tt.expectExpr)
			}
			if !stringsEqual(args, tt.expectArgs) {
				t.Errorf("initCommand() args = %v, want %v", args, tt.expectArgs)
			}
		})
	}
}

func TestProcessArgsWithExpressionFile(t *testing.T) {
	// Create a temporary .yq file with Windows line endings
	tempYqFile, err := os.CreateTemp("", "test.yq")
	if err != nil {
		t.Fatalf("Failed to create temp yq file: %v", err)
	}
	defer os.Remove(tempYqFile.Name())
	if _, err = tempYqFile.WriteString(".a.b\r\n.c.d"); err != nil {
		t.Fatalf("Failed to write to temp yq file: %v", err)
	}
	tempYqFile.Close()

	// Save original values
	originalExpressionFile := expressionFile
	defer func() { expressionFile = originalExpressionFile }()

	expressionFile = tempYqFile.Name()

	expr, args, err := processArgs([]string{"file1"})
	if err != nil {
		t.Errorf("processArgs() unexpected error: %v", err)
		return
	}

	expectedExpr := ".a.b\n.c.d" // Should convert \r\n to \n
	if expr != expectedExpr {
		t.Errorf("processArgs() expression = %v, want %v", expr, expectedExpr)
	}

	expectedArgs := []string{"file1"}
	if !stringsEqual(args, expectedArgs) {
		t.Errorf("processArgs() args = %v, want %v", args, expectedArgs)
	}
}

func TestProcessArgsWithNonExistentExpressionFile(t *testing.T) {
	// Save original values
	originalExpressionFile := expressionFile
	defer func() { expressionFile = originalExpressionFile }()

	expressionFile = "/path/that/does/not/exist"

	expr, args, err := processArgs([]string{"file1"})
	if err == nil {
		t.Errorf("processArgs() expected error but got none")
	}
	if expr != "" {
		t.Errorf("processArgs() expected empty expression but got %v", expr)
	}
	if args != nil {
		t.Errorf("processArgs() expected nil args but got %v", args)
	}
}

func TestInitCommandWithInvalidOutputFormat(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Save original values
	originalInputFormat := inputFormat
	originalOutputFormat := outputFormat
	defer func() {
		inputFormat = originalInputFormat
		outputFormat = originalOutputFormat
	}()

	inputFormat = "auto"
	outputFormat = "invalid"

	cmd := &cobra.Command{}
	expr, args, err := initCommand(cmd, []string{tempFile.Name()})
	if err == nil {
		t.Errorf("initCommand() expected error but got none")
	}
	if expr != "" {
		t.Errorf("initCommand() expected empty expression but got %v", expr)
	}
	if args != nil {
		t.Errorf("initCommand() expected nil args but got %v", args)
	}
}

func TestInitCommandWithUnknownInputFormat(t *testing.T) {
	// Create a temporary file with unknown extension
	tempFile, err := os.CreateTemp("", "test.unknown")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Save original values
	originalInputFormat := inputFormat
	originalOutputFormat := outputFormat
	defer func() {
		inputFormat = originalInputFormat
		outputFormat = originalOutputFormat
	}()

	inputFormat = "auto"
	outputFormat = "auto"

	cmd := &cobra.Command{}
	expr, args, err := initCommand(cmd, []string{tempFile.Name()})
	if err != nil {
		t.Errorf("initCommand() unexpected error: %v", err)
		return
	}

	// expr can be empty when no expression is provided
	_ = expr
	if args == nil {
		t.Errorf("initCommand() expected non-nil args")
	}
}

func TestConfigurePrinterWriterWithInvalidSplitExpression(t *testing.T) {
	// Save original value
	originalSplitFileExp := splitFileExp
	defer func() { splitFileExp = originalSplitFileExp }()

	splitFileExp = "[invalid expression"

	writer, err := configurePrinterWriter(&yqlib.Format{}, os.Stdout)
	if err == nil {
		t.Errorf("configurePrinterWriter() expected error but got none")
	}
	if writer != nil {
		t.Errorf("configurePrinterWriter() expected nil writer but got %v", writer)
	}
	if err != nil && !strings.Contains(err.Error(), "bad split document expression") {
		t.Errorf("configurePrinterWriter() error '%v' does not contain expected message", err.Error())
	}
}

func TestMaybeFileWithDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	result := maybeFile(tempDir)
	if result {
		t.Errorf("maybeFile(%q) = %v, want false", tempDir, result)
	}
}

func TestProcessStdInArgsWithDash(t *testing.T) {
	args := []string{"-", "file1"}
	result := processStdInArgs(args)
	if !stringsEqual(result, args) {
		t.Errorf("processStdInArgs() = %v, want %v", result, args)
	}
}

func TestProcessArgsWithYqFileExtension(t *testing.T) {
	tempYqFile, err := os.Create("test.yq")
	if err != nil {
		t.Fatalf("Failed to create temp yq file: %v", err)
	}
	defer os.Remove(tempYqFile.Name())
	if _, err = tempYqFile.WriteString(".a.b"); err != nil {
		t.Fatalf("Failed to write to temp yq file: %v", err)
	}
	tempYqFile.Close()

	// Save original values
	originalExpressionFile := expressionFile
	originalForceExpression := forceExpression
	defer func() {
		expressionFile = originalExpressionFile
		forceExpression = originalForceExpression
	}()

	// Reset expressionFile to empty to test the auto-detection
	expressionFile = ""
	forceExpression = ""

	// Debug: check the conditions manually
	t.Logf("expressionFile: %q", expressionFile)
	t.Logf("forceExpression: %q", forceExpression)
	t.Logf("tempYqFile.Name(): %q", tempYqFile.Name())
	t.Logf("strings.HasSuffix(tempYqFile.Name(), '.yq'): %v", strings.HasSuffix(tempYqFile.Name(), ".yq"))
	t.Logf("maybeFile(tempYqFile.Name()): %v", maybeFile(tempYqFile.Name()))

	// Test with only the yq file as argument (should be treated as expression file)
	expr, args, err := processArgs([]string{tempYqFile.Name()})
	if err != nil {
		t.Errorf("processArgs() unexpected error: %v", err)
		return
	}

	if expr != ".a.b" {
		t.Errorf("processArgs() expression = %v, want .a.b", expr)
	}

	expectedArgs := []string{}
	if !stringsEqual(args, expectedArgs) {
		t.Errorf("processArgs() args = %v, want %v", args, expectedArgs)
	}
}

func TestConfigureEncoderWithYamlFormat(t *testing.T) {
	// Save original values
	originalOutputFormat := outputFormat
	originalIndent := indent
	originalUnwrapScalar := unwrapScalar
	originalColorsEnabled := colorsEnabled
	originalNoDocSeparators := noDocSeparators
	defer func() {
		outputFormat = originalOutputFormat
		indent = originalIndent
		unwrapScalar = originalUnwrapScalar
		colorsEnabled = originalColorsEnabled
		noDocSeparators = originalNoDocSeparators
	}()

	outputFormat = "yaml"
	indent = 4
	unwrapScalar = true
	colorsEnabled = true
	noDocSeparators = true

	encoder, err := configureEncoder()
	if err != nil {
		t.Errorf("configureEncoder() unexpected error: %v", err)
		return
	}

	if encoder == nil {
		t.Errorf("configureEncoder() expected encoder but got nil")
	}
}

func TestConfigureEncoderWithPropertiesFormat(t *testing.T) {
	// Save original values
	originalOutputFormat := outputFormat
	originalIndent := indent
	originalUnwrapScalar := unwrapScalar
	originalColorsEnabled := colorsEnabled
	originalNoDocSeparators := noDocSeparators
	defer func() {
		outputFormat = originalOutputFormat
		indent = originalIndent
		unwrapScalar = originalUnwrapScalar
		colorsEnabled = originalColorsEnabled
		noDocSeparators = originalNoDocSeparators
	}()

	outputFormat = "properties"
	indent = 2
	unwrapScalar = false
	colorsEnabled = false
	noDocSeparators = false

	encoder, err := configureEncoder()
	if err != nil {
		t.Errorf("configureEncoder() unexpected error: %v", err)
		return
	}

	if encoder == nil {
		t.Errorf("configureEncoder() expected encoder but got nil")
	}
}

// Mock boolFlag for testing
type mockBoolFlag struct {
	explicitlySet bool
	value         bool
}

func (f *mockBoolFlag) IsExplicitlySet() bool {
	return f.explicitlySet
}

func (f *mockBoolFlag) IsSet() bool {
	return f.value
}

func (f *mockBoolFlag) String() string {
	return "mock"
}

func (f *mockBoolFlag) Set(_ string) error {
	return nil
}

func (f *mockBoolFlag) Type() string {
	return "bool"
}

// Helper function to compare string slices
func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSetupColors(t *testing.T) {
	tests := []struct {
		name         string
		forceColor   bool
		forceNoColor bool
		expectColors bool
	}{
		{
			name:         "force colour enabled",
			forceColor:   true,
			forceNoColor: false,
			expectColors: true,
		},
		{
			name:         "force no colour enabled",
			forceColor:   false,
			forceNoColor: true,
			expectColors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalForceColor := forceColor
			originalForceNoColor := forceNoColor
			originalColorsEnabled := colorsEnabled
			defer func() {
				forceColor = originalForceColor
				forceNoColor = originalForceNoColor
				colorsEnabled = originalColorsEnabled
			}()

			forceColor = tt.forceColor
			forceNoColor = tt.forceNoColor
			colorsEnabled = false // Reset to test the setting

			setupColors()

			if colorsEnabled != tt.expectColors {
				t.Errorf("setupColors() colorsEnabled = %v, want %v", colorsEnabled, tt.expectColors)
			}
		})
	}
}

func TestLoadSplitFileExpression(t *testing.T) {
	// Create a temporary file with expression content
	tempFile, err := os.CreateTemp("", "split")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err = tempFile.WriteString(".a.b"); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name             string
		splitFileExpFile string
		expectError      bool
		expectContent    string
	}{
		{
			name:             "load from file",
			splitFileExpFile: tempFile.Name(),
			expectError:      false,
			expectContent:    ".a.b",
		},
		{
			name:             "no file specified",
			splitFileExpFile: "",
			expectError:      false,
			expectContent:    "",
		},
		{
			name:             "non-existent file",
			splitFileExpFile: "/path/that/does/not/exist",
			expectError:      true,
			expectContent:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalSplitFileExpFile := splitFileExpFile
			originalSplitFileExp := splitFileExp
			defer func() {
				splitFileExpFile = originalSplitFileExpFile
				splitFileExp = originalSplitFileExp
			}()

			splitFileExpFile = tt.splitFileExpFile
			splitFileExp = ""

			err := loadSplitFileExpression()
			if tt.expectError {
				if err == nil {
					t.Errorf("loadSplitFileExpression() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("loadSplitFileExpression() unexpected error: %v", err)
				return
			}

			if splitFileExp != tt.expectContent {
				t.Errorf("loadSplitFileExpression() splitFileExp = %v, want %v", splitFileExp, tt.expectContent)
			}
		})
	}
}

func TestHandleBackwardsCompatibility(t *testing.T) {
	tests := []struct {
		name          string
		outputToJSON  bool
		initialFormat string
		expectFormat  string
	}{
		{
			name:          "outputToJSON true",
			outputToJSON:  true,
			initialFormat: "yaml",
			expectFormat:  "json",
		},
		{
			name:          "outputToJSON false",
			outputToJSON:  false,
			initialFormat: "yaml",
			expectFormat:  "yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalOutputToJSON := outputToJSON
			originalOutputFormat := outputFormat
			defer func() {
				outputToJSON = originalOutputToJSON
				outputFormat = originalOutputFormat
			}()

			outputToJSON = tt.outputToJSON
			outputFormat = tt.initialFormat

			handleBackwardsCompatibility()

			if outputFormat != tt.expectFormat {
				t.Errorf("handleBackwardsCompatibility() outputFormat = %v, want %v", outputFormat, tt.expectFormat)
			}
		})
	}
}

func TestValidateCommandFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		writeInplace  bool
		frontMatter   string
		splitFileExp  string
		nullInput     bool
		expectError   bool
		errorContains string
	}{
		{
			name:         "valid flags",
			args:         []string{"file.yaml"},
			writeInplace: false,
			frontMatter:  "",
			splitFileExp: "",
			nullInput:    false,
			expectError:  false,
		},
		{
			name:          "write inplace with no args",
			args:          []string{},
			writeInplace:  true,
			frontMatter:   "",
			splitFileExp:  "",
			nullInput:     false,
			expectError:   true,
			errorContains: "write in place flag only applicable when giving an expression and at least one file",
		},
		{
			name:          "write inplace with dash",
			args:          []string{"-"},
			writeInplace:  true,
			frontMatter:   "",
			splitFileExp:  "",
			nullInput:     false,
			expectError:   true,
			errorContains: "write in place flag only applicable when giving an expression and at least one file",
		},
		{
			name:          "front matter with no args",
			args:          []string{},
			writeInplace:  false,
			frontMatter:   "extract",
			splitFileExp:  "",
			nullInput:     false,
			expectError:   true,
			errorContains: "front matter flag only applicable when giving an expression and at least one file",
		},
		{
			name:          "write inplace with split file",
			args:          []string{"file.yaml"},
			writeInplace:  true,
			frontMatter:   "",
			splitFileExp:  ".a.b",
			nullInput:     false,
			expectError:   true,
			errorContains: "write in place cannot be used with split file",
		},
		{
			name:          "null input with args",
			args:          []string{"file.yaml"},
			writeInplace:  false,
			frontMatter:   "",
			splitFileExp:  "",
			nullInput:     true,
			expectError:   true,
			errorContains: "cannot pass files in when using null-input flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalWriteInplace := writeInplace
			originalFrontMatter := frontMatter
			originalSplitFileExp := splitFileExp
			originalNullInput := nullInput
			defer func() {
				writeInplace = originalWriteInplace
				frontMatter = originalFrontMatter
				splitFileExp = originalSplitFileExp
				nullInput = originalNullInput
			}()

			writeInplace = tt.writeInplace
			frontMatter = tt.frontMatter
			splitFileExp = tt.splitFileExp
			nullInput = tt.nullInput

			err := validateCommandFlags(tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateCommandFlags() expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("validateCommandFlags() error '%v' does not contain '%v'", err.Error(), tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("validateCommandFlags() unexpected error: %v", err)
			}
		})
	}
}

func TestConfigureFormats(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		inputFormat  string
		outputFormat string
		expectError  bool
	}{
		{
			name:         "valid formats",
			args:         []string{"file.yaml"},
			inputFormat:  "auto",
			outputFormat: "auto",
			expectError:  false,
		},
		{
			name:         "invalid output format",
			args:         []string{"file.yaml"},
			inputFormat:  "auto",
			outputFormat: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalInputFormat := inputFormat
			originalOutputFormat := outputFormat
			defer func() {
				inputFormat = originalInputFormat
				outputFormat = originalOutputFormat
			}()

			inputFormat = tt.inputFormat
			outputFormat = tt.outputFormat

			err := configureFormats(tt.args)
			if tt.expectError {
				if err == nil {
					t.Errorf("configureFormats() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("configureFormats() unexpected error: %v", err)
			}
		})
	}
}

func TestConfigureInputFormat(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		inputFormat   string
		outputFormat  string
		expectInput   string
		expectOutput  string
	}{
		{
			name:          "auto format with yaml file",
			inputFilename: "file.yaml",
			inputFormat:   "auto",
			outputFormat:  "auto",
			expectInput:   "yaml",
			expectOutput:  "yaml",
		},
		{
			name:          "auto format with json file",
			inputFilename: "file.json",
			inputFormat:   "auto",
			outputFormat:  "auto",
			expectInput:   "json",
			expectOutput:  "json",
		},
		{
			name:          "auto format with unknown file",
			inputFilename: "file.unknown",
			inputFormat:   "auto",
			outputFormat:  "auto",
			expectInput:   "yaml",
			expectOutput:  "yaml",
		},
		{
			name:          "explicit format",
			inputFilename: "file.yaml",
			inputFormat:   "json",
			outputFormat:  "auto",
			expectInput:   "json",
			expectOutput:  "yaml", // backwards compatibility
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalInputFormat := inputFormat
			originalOutputFormat := outputFormat
			defer func() {
				inputFormat = originalInputFormat
				outputFormat = originalOutputFormat
			}()

			inputFormat = tt.inputFormat
			outputFormat = tt.outputFormat

			err := configureInputFormat(tt.inputFilename)
			if err != nil {
				t.Errorf("configureInputFormat() unexpected error: %v", err)
				return
			}

			if inputFormat != tt.expectInput {
				t.Errorf("configureInputFormat() inputFormat = %v, want %v", inputFormat, tt.expectInput)
			}
			if outputFormat != tt.expectOutput {
				t.Errorf("configureInputFormat() outputFormat = %v, want %v", outputFormat, tt.expectOutput)
			}
		})
	}
}

func TestConfigureOutputFormat(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		expectError  bool
		expectUnwrap bool
	}{
		{
			name:         "yaml format",
			outputFormat: "yaml",
			expectError:  false,
			expectUnwrap: true,
		},
		{
			name:         "properties format",
			outputFormat: "properties",
			expectError:  false,
			expectUnwrap: true,
		},
		{
			name:         "json format",
			outputFormat: "json",
			expectError:  false,
			expectUnwrap: false,
		},
		{
			name:         "invalid format",
			outputFormat: "invalid",
			expectError:  true,
			expectUnwrap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalOutputFormat := outputFormat
			originalUnwrapScalar := unwrapScalar
			defer func() {
				outputFormat = originalOutputFormat
				unwrapScalar = originalUnwrapScalar
			}()

			outputFormat = tt.outputFormat
			unwrapScalar = false // Reset to test the setting

			err := configureOutputFormat()
			if tt.expectError {
				if err == nil {
					t.Errorf("configureOutputFormat() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("configureOutputFormat() unexpected error: %v", err)
				return
			}

			if unwrapScalar != tt.expectUnwrap {
				t.Errorf("configureOutputFormat() unwrapScalar = %v, want %v", unwrapScalar, tt.expectUnwrap)
			}
		})
	}
}

func TestConfigureUnwrapScalar(t *testing.T) {
	tests := []struct {
		name          string
		explicitlySet bool
		flagValue     bool
		initialUnwrap bool
		expectUnwrap  bool
	}{
		{
			name:          "flag not explicitly set",
			explicitlySet: false,
			flagValue:     true,
			initialUnwrap: true,
			expectUnwrap:  true, // Should remain unchanged
		},
		{
			name:          "flag explicitly set to true",
			explicitlySet: true,
			flagValue:     true,
			initialUnwrap: false,
			expectUnwrap:  true,
		},
		{
			name:          "flag explicitly set to false",
			explicitlySet: true,
			flagValue:     false,
			initialUnwrap: true,
			expectUnwrap:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			originalUnwrapScalar := unwrapScalar
			originalUnwrapScalarFlag := unwrapScalarFlag
			defer func() {
				unwrapScalar = originalUnwrapScalar
				unwrapScalarFlag = originalUnwrapScalarFlag
			}()

			unwrapScalar = tt.initialUnwrap
			unwrapScalarFlag = &mockBoolFlag{
				explicitlySet: tt.explicitlySet,
				value:         tt.flagValue,
			}

			configureUnwrapScalar()

			if unwrapScalar != tt.expectUnwrap {
				t.Errorf("configureUnwrapScalar() unwrapScalar = %v, want %v", unwrapScalar, tt.expectUnwrap)
			}
		})
	}
}
