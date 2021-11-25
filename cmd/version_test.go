package cmd

import "testing"

func TestGetVersionDisplay(t *testing.T) {
	expectedVersion := ProductName + " (https://github.com/mikefarah/yq/) version " + Version
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
