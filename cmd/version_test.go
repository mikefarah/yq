package cmd

import (
	"strings"
	"testing"
)

func TestGetVersionDisplay(t *testing.T) {
	var expectedVersion = ProductName + " (https://github.com/mikefarah/yq/) version " + Version
	if VersionPrerelease != "" {
		expectedVersion = expectedVersion + "-" + VersionPrerelease
	}
	expectedVersion = expectedVersion + "\n"
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Display Version",
			want: expectedVersion,
		},
	}
	for _, tt := range tests {
		if got := GetVersionDisplay(); got != tt.want {
			t.Errorf("%q. GetVersionDisplay() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getHumanVersion(t *testing.T) {
	// Save original values
	origGitDescribe := GitDescribe
	origGitCommit := GitCommit
	origVersionPrerelease := VersionPrerelease

	// Restore after test
	defer func() {
		GitDescribe = origGitDescribe
		GitCommit = origGitCommit
		VersionPrerelease = origVersionPrerelease
	}()

	GitDescribe = "e42813d"
	GitCommit = "e42813d+CHANGES"
	var wanted string
	if VersionPrerelease == "" {
		wanted = GitDescribe
	} else {
		wanted = "e42813d-" + VersionPrerelease + " (e42813d+CHANGES)"
	}

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Git Variables defined",
			want: wanted,
		},
	}
	for _, tt := range tests {
		if got := getHumanVersion(); got != tt.want {
			t.Errorf("%q. getHumanVersion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getHumanVersion_NoGitDescribe(t *testing.T) {
	// Save original values
	origGitDescribe := GitDescribe
	origGitCommit := GitCommit
	origVersionPrerelease := VersionPrerelease

	// Restore after test
	defer func() {
		GitDescribe = origGitDescribe
		GitCommit = origGitCommit
		VersionPrerelease = origVersionPrerelease
	}()

	GitDescribe = ""
	GitCommit = ""
	VersionPrerelease = ""

	got := getHumanVersion()
	if got != Version {
		t.Errorf("getHumanVersion() = %v, want %v", got, Version)
	}
}

func Test_getHumanVersion_WithPrerelease(t *testing.T) {
	// Save original values
	origGitDescribe := GitDescribe
	origGitCommit := GitCommit
	origVersionPrerelease := VersionPrerelease

	// Restore after test
	defer func() {
		GitDescribe = origGitDescribe
		GitCommit = origGitCommit
		VersionPrerelease = origVersionPrerelease
	}()

	GitDescribe = ""
	GitCommit = "abc123"
	VersionPrerelease = "beta"

	got := getHumanVersion()
	expected := Version + "-beta (abc123)"
	if got != expected {
		t.Errorf("getHumanVersion() = %v, want %v", got, expected)
	}
}

func Test_getHumanVersion_PrereleaseInVersion(t *testing.T) {
	// Save original values
	origGitDescribe := GitDescribe
	origGitCommit := GitCommit
	origVersionPrerelease := VersionPrerelease

	// Restore after test
	defer func() {
		GitDescribe = origGitDescribe
		GitCommit = origGitCommit
		VersionPrerelease = origVersionPrerelease
	}()

	GitDescribe = "v1.2.3-rc1"
	GitCommit = "xyz789"
	VersionPrerelease = "rc1"

	got := getHumanVersion()
	// Should not duplicate "rc1" since it's already in GitDescribe
	expected := "v1.2.3-rc1 (xyz789)"
	if got != expected {
		t.Errorf("getHumanVersion() = %v, want %v", got, expected)
	}
}

func Test_getHumanVersion_StripSingleQuotes(t *testing.T) {
	// Save original values
	origGitDescribe := GitDescribe
	origGitCommit := GitCommit
	origVersionPrerelease := VersionPrerelease

	// Restore after test
	defer func() {
		GitDescribe = origGitDescribe
		GitCommit = origGitCommit
		VersionPrerelease = origVersionPrerelease
	}()

	GitDescribe = "'v1.2.3'"
	GitCommit = "'commit123'"
	VersionPrerelease = ""

	got := getHumanVersion()
	// Should strip single quotes
	if strings.Contains(got, "'") {
		t.Errorf("getHumanVersion() = %v, should not contain single quotes", got)
	}
	expected := "v1.2.3"
	if got != expected {
		t.Errorf("getHumanVersion() = %v, want %v", got, expected)
	}
}

func TestProductName(t *testing.T) {
	if ProductName != "yq" {
		t.Errorf("ProductName = %v, want yq", ProductName)
	}
}

func TestVersionIsSet(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if !strings.HasPrefix(Version, "v") {
		t.Errorf("Version %v should start with 'v'", Version)
	}
}
