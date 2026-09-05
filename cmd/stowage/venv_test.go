package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVenvInstallAndRemove(t *testing.T) {
	env, err := newVenv(filepath.Join(t.TempDir(), "myenv"))
	assert.NoError(t, err)

	assert.NoError(t, env.ensureDirs())

	spec := Specification{Name: "tool", Image: "example/tool"}
	assert.NoError(t, env.install(&spec))

	store := env.storage()
	loaded, err := store.loadSpecificationByName("tool")
	assert.NoError(t, err)
	assert.Equal(t, "example/tool", loaded.Image)

	wrapperPath := filepath.Join(env.binDir(), "tool")
	if runtime.GOOS == "windows" {
		wrapperPath += ".exe"
	}
	_, err = os.Stat(wrapperPath)
	assert.NoError(t, err)

	assert.NoError(t, env.remove("tool"))

	_, err = store.loadSpecificationByName("tool")
	assert.Error(t, err)

	_, err = os.Stat(wrapperPath)
	assert.True(t, os.IsNotExist(err))
}

func TestVenvIsolatedFromGlobalStorage(t *testing.T) {
	env, err := newVenv(filepath.Join(t.TempDir(), "myenv"))
	assert.NoError(t, err)
	assert.NoError(t, env.ensureDirs())

	spec := Specification{Name: "tool", Image: "venv/tool"}
	assert.NoError(t, env.install(&spec))

	globalStore := testStorage(t)
	_, err = globalStore.loadSpecificationByName("tool")
	assert.Error(t, err, "installing into a venv must not touch unrelated storage")
}

func TestVenvWriteActivationScripts(t *testing.T) {
	env, err := newVenv(filepath.Join(t.TempDir(), "myenv"))
	assert.NoError(t, err)
	assert.NoError(t, env.ensureDirs())

	assert.NoError(t, env.writeActivationScripts())

	for _, name := range []string{"activate", "Activate.ps1", "activate.bat", "deactivate.bat"} {
		_, err := os.Stat(filepath.Join(env.binDir(), name))
		assert.NoError(t, err, "expected %s to be written", name)
	}
}
