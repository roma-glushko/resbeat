package resbeat

import (
	"fmt"
	"runtime"
)

func GetVersion(version string, commitSha string) string {
	return fmt.Sprintf("%s (commit: %s, %s)", version, commitSha, runtime.Version())
}
