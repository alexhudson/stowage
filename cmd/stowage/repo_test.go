package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoAddSpecificationAppendsNew(t *testing.T) {
	repo := Repository{Name: "myrepo"}

	repo.addSpecification(&Specification{Name: "one"})
	repo.addSpecification(&Specification{Name: "two"})

	assert.Len(t, repo.Listing, 2)
	assert.Equal(t, "one", repo.Listing[0].Name)
	assert.Equal(t, "two", repo.Listing[1].Name)
}

func TestRepoAddSpecificationReplacesExisting(t *testing.T) {
	repo := Repository{
		Name: "myrepo",
		Listing: []RepositoryEntry{
			{Name: "one", Description: "original"},
		},
	}

	repo.addSpecification(&Specification{Name: "one"})

	assert.Len(t, repo.Listing, 1)
	assert.Equal(t, "", repo.Listing[0].Description, "replaced entry should not keep the old description")
}

func TestRepoSearchMatchesEntryName(t *testing.T) {
	repo := Repository{
		Name: "myrepo",
		Listing: []RepositoryEntry{
			{Name: "hello-world", Description: "says hello"},
			{Name: "other", Description: "does something else"},
		},
	}

	result := repo.search("hello")

	assert.Len(t, result, 1)
	assert.Equal(t, "hello-world", result[0].Name)
}

func TestRepoSearchMatchesDescription(t *testing.T) {
	repo := Repository{
		Name: "myrepo",
		Listing: []RepositoryEntry{
			{Name: "tool", Description: "a very special tool"},
			{Name: "other", Description: "boring"},
		},
	}

	result := repo.search("special")

	assert.Len(t, result, 1)
	assert.Equal(t, "tool", result[0].Name)
}

func TestRepoSearchMatchesRepoName(t *testing.T) {
	repo := Repository{
		Name: "azure-tools",
		Listing: []RepositoryEntry{
			{Name: "cli", Description: "the cli"},
			{Name: "portal", Description: "the portal"},
		},
	}

	result := repo.search("azure")

	assert.Len(t, result, 2, "matching the repo name should include every entry")
}
