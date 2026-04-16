package version

import (
	"fmt"
	"runtime"
)

// Info holds build-time version metadata.
type Info struct {
	Version   string
	Commit    string
	BuildDate string
	GoVersion string
}

// These are set via -ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// Get returns the current version Info.
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("vaultline %s (commit=%s built=%s %s)",
		i.Version, i.Commit, i.BuildDate, i.GoVersion)
}
