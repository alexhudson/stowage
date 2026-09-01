package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStorage(t *testing.T) storage {
	dir := t.TempDir()
	return newStorage(
		filepath.Join(dir, "spec"),
		filepath.Join(dir, "repo"),
	)
}

func TestStorageSaveLoadSpecification(t *testing.T) {
	s := testStorage(t)

	spec := Specification{
		Name:        "roundtrip",
		Image:       "example/image",
		Environment: []string{"FOO=bar"},
	}

	err := s.saveSpecification(&spec)
	assert.NoError(t, err)

	loaded, err := s.loadSpecificationByName("roundtrip")
	assert.NoError(t, err)
	assert.Equal(t, spec.Name, loaded.Name)
	assert.Equal(t, spec.Image, loaded.Image)
	assert.Equal(t, spec.Environment, loaded.Environment)
}

func TestStorageRemoveSpecification(t *testing.T) {
	s := testStorage(t)

	spec := Specification{Name: "removeme"}
	assert.NoError(t, s.saveSpecification(&spec))

	_, err := s.loadSpecificationByName("removeme")
	assert.NoError(t, err)

	assert.NoError(t, s.removeSpecification(&spec))

	_, err = s.loadSpecificationByName("removeme")
	assert.Error(t, err)
}

func TestStorageListSpecifications(t *testing.T) {
	s := testStorage(t)

	assert.Empty(t, s.listSpecifications())

	assert.NoError(t, s.saveSpecification(&Specification{Name: "one"}))
	assert.NoError(t, s.saveSpecification(&Specification{Name: "two"}))

	names := s.listSpecifications()
	assert.ElementsMatch(t, []string{"one", "two"}, names)
}
