package main

import (
	"io"
	"testing"

	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

// TestGoInstallCompatibility ensures the module can be zipped for go install.
// This is an integration test that uses the same zip.CreateFromDir function
// that go install uses internally. If this test fails, go install will fail.
// See: https://github.com/mikefarah/yq/issues/2587
func TestGoInstallCompatibility(t *testing.T) {
	mod := module.Version{
		Path:    "github.com/mikefarah/yq/v4",
		Version: "v4.0.0", // the actual version doesn't matter for validation
	}

	if err := zip.CreateFromDir(io.Discard, mod, "."); err != nil {
		t.Fatalf("Module cannot be zipped for go install: %v", err)
	}
}
