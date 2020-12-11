package version

import "fmt"

var (
	version string
	commit  string
)

// String returns the version.
func String() string {
	return fmt.Sprintf("%s-%s", version, commit)
}
