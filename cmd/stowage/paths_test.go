package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnixInstallPrefixDefaultsToUserLocal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unixInstallPrefix only applies on non-Windows platforms")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root - unixInstallPrefix falls back to /usr/local")
	}

	home, err := os.UserHomeDir()
	assert.NoError(t, err)

	assert.Equal(t, filepath.Join(home, ".local"), unixInstallPrefix())
}
