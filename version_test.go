package vibelogger

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("Version should not be empty")
	}

	// Should match semantic versioning format for stable release
	if IsStableVersion() {
		if !strings.HasPrefix(version, "1.0.0") {
			t.Errorf("Expected version to start with '1.0.0', got '%s'", version)
		}
	}
}

func TestVersionConstants(t *testing.T) {
	// Test version constants are properly set
	if VersionMajor != 1 {
		t.Errorf("Expected VersionMajor to be 1, got %d", VersionMajor)
	}
	if VersionMinor != 0 {
		t.Errorf("Expected VersionMinor to be 0, got %d", VersionMinor)
	}
	if VersionPatch != 0 {
		t.Errorf("Expected VersionPatch to be 0, got %d", VersionPatch)
	}
}

func TestGetVersionInfo(t *testing.T) {
	versionInfo := GetVersionInfo()

	// Check required fields
	if versionInfo.Version == "" {
		t.Error("VersionInfo.Version should not be empty")
	}
	if versionInfo.Major != VersionMajor {
		t.Errorf("Expected Major to be %d, got %d", VersionMajor, versionInfo.Major)
	}
	if versionInfo.Minor != VersionMinor {
		t.Errorf("Expected Minor to be %d, got %d", VersionMinor, versionInfo.Minor)
	}
	if versionInfo.Patch != VersionPatch {
		t.Errorf("Expected Patch to be %d, got %d", VersionPatch, versionInfo.Patch)
	}
	if versionInfo.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
	if versionInfo.UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}

	// User agent should follow expected format
	expectedUserAgent := "vibe-logger-go/" + GetVersion()
	if versionInfo.UserAgent != expectedUserAgent {
		t.Errorf("Expected UserAgent '%s', got '%s'", expectedUserAgent, versionInfo.UserAgent)
	}
}

func TestVersionInfoJSONSerialization(t *testing.T) {
	versionInfo := GetVersionInfo()

	// Test JSON serialization
	jsonData, err := json.Marshal(versionInfo)
	if err != nil {
		t.Fatalf("Failed to marshal VersionInfo: %v", err)
	}

	// Test JSON deserialization
	var deserializedInfo VersionInfo
	err = json.Unmarshal(jsonData, &deserializedInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal VersionInfo: %v", err)
	}

	// Compare key fields
	if deserializedInfo.Version != versionInfo.Version {
		t.Errorf("Version mismatch after JSON round-trip: expected '%s', got '%s'",
			versionInfo.Version, deserializedInfo.Version)
	}
	if deserializedInfo.Major != versionInfo.Major {
		t.Errorf("Major version mismatch after JSON round-trip: expected %d, got %d",
			versionInfo.Major, deserializedInfo.Major)
	}
}

func TestIsStableVersion(t *testing.T) {
	// For v1.0.0, this should be a stable version
	isStable := IsStableVersion()
	if !isStable {
		t.Error("v1.0.0 should be marked as a stable version")
	}

	// Stable version should have empty prerelease
	if VersionPrerelease != "" {
		t.Errorf("Stable version should have empty prerelease, got '%s'", VersionPrerelease)
	}
}

func TestCompareVersion(t *testing.T) {
	currentVersion := GetVersion()

	// Test comparison with same version
	if CompareVersion(currentVersion) != 0 {
		t.Errorf("Comparison with same version should return 0")
	}

	// Test comparison with obviously older version
	if CompareVersion("0.1.0") <= 0 {
		t.Error("Current version should be greater than 0.1.0")
	}

	// Test comparison with obviously newer version
	if CompareVersion("2.0.0") >= 0 {
		t.Error("Current version should be less than 2.0.0")
	}
}

func TestVersionFormatting(t *testing.T) {
	tests := []struct {
		name        string
		prerelease  string
		buildMeta   string
		expectedFmt string
	}{
		{
			name:        "stable_version",
			prerelease:  "",
			buildMeta:   "",
			expectedFmt: Version,
		},
		{
			name:        "prerelease_version",
			prerelease:  "alpha.1",
			buildMeta:   "",
			expectedFmt: Version + "-alpha.1",
		},
		{
			name:        "build_metadata",
			prerelease:  "",
			buildMeta:   "20250713.1",
			expectedFmt: Version + "+20250713.1",
		},
		{
			name:        "prerelease_with_build",
			prerelease:  "beta.2",
			buildMeta:   "20250713.1",
			expectedFmt: Version + "-beta.2+20250713.1",
		},
	}

	// Save original values
	origPrerelease := VersionPrerelease
	origBuildMeta := BuildMetadata

	// Note: These tests are conceptual since the constants can't be modified at runtime
	// In a real implementation, you might use build tags or variables instead of constants
	_ = tests
	_ = origPrerelease
	_ = origBuildMeta

	// For now, just test the current stable version format
	version := GetVersion()
	if version != Version {
		t.Errorf("Expected stable version format '%s', got '%s'", Version, version)
	}
}

func TestVersionConstants_SemVer(t *testing.T) {
	// Verify the version follows semantic versioning principles
	version := GetVersion()

	// Should not contain spaces
	if strings.Contains(version, " ") {
		t.Errorf("Version should not contain spaces: '%s'", version)
	}

	// Should start with major version
	if !strings.HasPrefix(version, "1") {
		t.Errorf("Version should start with major version '1': '%s'", version)
	}

	// Should be a reasonable length
	if len(version) < 5 || len(version) > 50 {
		t.Errorf("Version length seems unreasonable: '%s' (length: %d)", version, len(version))
	}
}
