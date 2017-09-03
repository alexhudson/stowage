package main

import (
	"encoding/json"
	"strings"
)

// RepositoryEntry represents a command we can install
type RepositoryEntry struct {
	Name        string
	Description string
	Author      string
}

// Repository represents a local file repository
type Repository struct {
	Name        string
	URI         string
	Description string
	Listing     []RepositoryEntry
}

func (r *Repository) fromJSON(byt []byte) error {
	if err := json.Unmarshal(byt, &r); err != nil {
		panic(err)
	}
	return nil
}

func (r *Repository) toJSON() []byte {
	ret, _ := json.MarshalIndent(r, "", "\t")
	return ret
}

func (r *Repository) addSpecification(s *Specification) {
	entry := RepositoryEntry{
		Name: s.Name,
	}
	replaced := false
	for i, spec := range r.Listing {
		if spec.Name == entry.Name {
			// this spec is already listed, so we're going to update it instead
			r.Listing[i] = entry
			replaced = true
		}
	}
	if !replaced {
		r.Listing = append(r.Listing, entry)
	}
}

func (r *Repository) getURLForSpec(name string) string {
	return r.URI + "/" + name + ".json"
}

func (r *Repository) search(term string) []RepositoryEntry {
	result := make([]RepositoryEntry, 0)

	var include bool
	includeDefault := false
	if strings.Index(r.Name, term) > -1 {
		includeDefault = true
	}

	for _, entry := range r.Listing {
		include = includeDefault

		if strings.Index(entry.Name, term) > -1 {
			include = true
		}
		if strings.Index(entry.Description, term) > -1 {
			include = true
		}

		if include {
			result = append(result, entry)
		}
	}

	return result
}
