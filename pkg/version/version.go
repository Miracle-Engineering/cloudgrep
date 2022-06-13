package version

import (
	"flag"
	"runtime"
)

var (
	// Version is the current application version
	Version = "dev"

	// GitCommit contains the Git commit SHA for the binary
	GitCommit string

	// BuildTime contains the binary build time
	BuildTime string

	// GoVersion contains the build time Go version
	GoVersion string
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"gitCommit,omitempty"`
	// BuildTime is the state of the git tree.
	BuildTime string `json:"buildTime,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"goVersion,omitempty"`
}

// Get returns build info
func Get() BuildInfo {
	v := BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
	}

	// Strip out GoVersion during a test run for consistent test output
	if flag.Lookup("test.v") != nil {
		v.GoVersion = "testing"
	}
	return v
}

func IsDev() bool {
	return Version == "dev"
}
