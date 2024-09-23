package vcs

import (
	"fmt"
	"runtime/debug"
)

// Version retrieves the version of the current build.
// The function reads the build information to get two value pairs: "vcs.revision", "vcs.modified".
// If "vcs.modified" value is "true", it returns the revision string with "-dirty" suffix.
// Otherwise, it just returns the revision string.
func Version() string {
	var revision string
	var modified bool

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
