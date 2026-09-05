package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageSaveLoadRepository(t *testing.T) {
	s := testStorage(t)

	repo := Repository{
		Name:        "myrepo",
		URI:         "https://example.com/repo",
		Description: "example repo",
		Listing: []RepositoryEntry{
			{Name: "tool", Description: "a tool"},
		},
	}

	err := s.saveRepositoryByName(&repo, repo.Name)
	assert.NoError(t, err)

	loaded, err := s.loadRepositoryByName("myrepo")
	assert.NoError(t, err)
	assert.Equal(t, repo.Name, loaded.Name)
	assert.Equal(t, repo.URI, loaded.URI)
	assert.Equal(t, repo.Listing, loaded.Listing)
}

func TestStorageListRepositories(t *testing.T) {
	s := testStorage(t)

	assert.Empty(t, s.listRepositories())

	assert.NoError(t, s.saveRepositoryByName(&Repository{Name: "one"}, "one"))
	assert.NoError(t, s.saveRepositoryByName(&Repository{Name: "two"}, "two"))

	names := s.listRepositories()
	assert.ElementsMatch(t, []string{"one", "two"}, names)
}
