package resbeat

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestLogging_CreateLogger(t *testing.T) {
	appVer := "1.0.0"
	commitSha := "abcd"

	versionStr := GetVersion(appVer, commitSha)

	assert.Contains(t, versionStr, appVer)
	assert.Contains(t, versionStr, commitSha)
	assert.Contains(t, versionStr, runtime.Version())
}
