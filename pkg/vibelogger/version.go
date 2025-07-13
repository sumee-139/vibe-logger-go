package vibelogger

import "fmt"

// Version information for the vibe-logger-go package
const (
	// Version is the current version of vibe-logger-go
	Version = "1.0.0"

	// VersionMajor is the major version number
	VersionMajor = 1

	// VersionMinor is the minor version number
	VersionMinor = 0

	// VersionPatch is the patch version number
	VersionPatch = 0

	// VersionPrerelease is the pre-release version info (empty for stable releases)
	VersionPrerelease = ""

	// BuildMetadata contains build metadata (empty for releases)
	BuildMetadata = ""
)

// VersionInfo contains detailed version information
type VersionInfo struct {
	Version    string `json:"version"`
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Prerelease string `json:"prerelease,omitempty"`
	BuildMeta  string `json:"build_metadata,omitempty"`
	GoVersion  string `json:"go_version"`
	UserAgent  string `json:"user_agent"`
}

// GetVersion returns the current version string
func GetVersion() string {
	if VersionPrerelease != "" {
		if BuildMetadata != "" {
			return fmt.Sprintf("%s-%s+%s", Version, VersionPrerelease, BuildMetadata)
		}
		return fmt.Sprintf("%s-%s", Version, VersionPrerelease)
	}
	if BuildMetadata != "" {
		return fmt.Sprintf("%s+%s", Version, BuildMetadata)
	}
	return Version
}

// GetVersionInfo returns detailed version information
func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		Version:    GetVersion(),
		Major:      VersionMajor,
		Minor:      VersionMinor,
		Patch:      VersionPatch,
		Prerelease: VersionPrerelease,
		BuildMeta:  BuildMetadata,
		GoVersion:  getGoVersion(),
		UserAgent:  fmt.Sprintf("vibe-logger-go/%s", GetVersion()),
	}
}

// IsStableVersion returns true if this is a stable (non-prerelease) version
func IsStableVersion() bool {
	return VersionPrerelease == ""
}

// CompareVersion compares the current version with another version string
// Returns: -1 if current < other, 0 if equal, 1 if current > other
func CompareVersion(other string) int {
	// Simple string comparison for now
	// For production use, consider using semantic version comparison
	current := GetVersion()
	if current < other {
		return -1
	} else if current > other {
		return 1
	}
	return 0
}

// getGoVersion returns the Go runtime version
func getGoVersion() string {
	return getEnvironment()["go_version"]
}
