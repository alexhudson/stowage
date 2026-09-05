package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeSpecFile(t *testing.T, dir string, filename string, spec Specification) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, filename), spec.toJSON(), 0644)
	assert.NoError(t, err)
}

func TestRepositoryDirScanBuildsListing(t *testing.T) {
	dir := t.TempDir()

	writeSpecFile(t, dir, "one.json", Specification{Name: "one"})
	writeSpecFile(t, dir, "two.json", Specification{Name: "two"})

	rd := RepositoryDir{Path: dir}
	rd.scan()

	names := make([]string, 0, len(rd.Repo.Listing))
	for _, entry := range rd.Repo.Listing {
		names = append(names, entry.Name)
	}

	assert.ElementsMatch(t, []string{"one", "two"}, names)
}

func TestRepositoryDirScanSkipsRepoFile(t *testing.T) {
	dir := t.TempDir()

	writeSpecFile(t, dir, "one.json", Specification{Name: "one"})
	err := os.WriteFile(filepath.Join(dir, "stowage.json"), []byte(`{"Name":"stowage.json"}`), 0644)
	assert.NoError(t, err)

	rd := RepositoryDir{Path: dir}
	rd.scan()

	assert.Len(t, rd.Repo.Listing, 1)
	assert.Equal(t, "one", rd.Repo.Listing[0].Name)
}
