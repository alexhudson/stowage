package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVenvManifestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), defaultVenvManifestName)

	manifest := venvManifest{}
	manifest.upsert(venvManifestEntry{Name: "python", Request: "python:3.11"})

	assert.NoError(t, saveVenvManifest(path, manifest))

	loaded, err := loadVenvManifest(path)
	assert.NoError(t, err)
	assert.Len(t, loaded.Entries, 1)
	assert.Equal(t, "python", loaded.Entries[0].Name)
	assert.Equal(t, "python:3.11", loaded.Entries[0].Request)
}

func TestVenvManifestLoadMissing(t *testing.T) {
	_, err := loadVenvManifest(filepath.Join(t.TempDir(), "missing.json"))
	assert.Error(t, err)
}

func TestVenvManifestUpsertReplacesExistingByName(t *testing.T) {
	manifest := venvManifest{}
	manifest.upsert(venvManifestEntry{Name: "python", Request: "python:3.10"})
	manifest.upsert(venvManifestEntry{Name: "python", Request: "python:3.11"})

	assert.Len(t, manifest.Entries, 1)
	assert.Equal(t, "python:3.11", manifest.Entries[0].Request)
}

func TestVenvManifestRemoveByName(t *testing.T) {
	manifest := venvManifest{}
	manifest.upsert(venvManifestEntry{Name: "python", Request: "python:3.11"})
	manifest.upsert(venvManifestEntry{Name: "node", Request: "node:18"})

	manifest.removeByName("python")

	assert.Len(t, manifest.Entries, 1)
	assert.Equal(t, "node", manifest.Entries[0].Name)
}
