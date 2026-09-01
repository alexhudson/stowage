package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallerLoadSpecFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "myspec.json")

	spec := Specification{Name: "myspec", Image: "example/image"}
	err := os.WriteFile(path, spec.toJSON(), 0644)
	assert.NoError(t, err)

	installer := Installer{}
	ok := installer.loadSpecFromFile(path)

	assert.True(t, ok)
	assert.Equal(t, "myspec", installer.Spec.Name)
	assert.Equal(t, "example/image", installer.Spec.Image)
}

func TestInstallerLoadSpecFromFileMissing(t *testing.T) {
	installer := Installer{}
	ok := installer.loadSpecFromFile(filepath.Join(t.TempDir(), "missing.json"))

	assert.False(t, ok)
}

func TestInstallerSetCommandNameAndEntrypoint(t *testing.T) {
	installer := Installer{Spec: &Specification{Name: "original"}}

	installer.setCommandName("renamed")
	installer.setEntrypoint("/bin/renamed")

	assert.Equal(t, "renamed", installer.Spec.Name)
	assert.Equal(t, "/bin/renamed", installer.Spec.Command)
}

func TestInstallerLoadSpecFromURL(t *testing.T) {
	spec := Specification{Name: "fromurl", Image: "example/image"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(spec.toJSON())
	}))
	defer server.Close()

	installer := Installer{}
	ok := installer.loadSpecFromURL(server.URL)

	assert.True(t, ok)
	assert.Equal(t, "fromurl", installer.Spec.Name)
	assert.Equal(t, "example/image", installer.Spec.Image)
}
