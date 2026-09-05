package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// defaultVenvManifestName is the project-local file `stowage venv` reads to
// bootstrap an environment; `stowage install`/`install-spec`/`uninstall`
// keep it in sync whenever they act on an active environment, so
// re-creating the environment elsewhere (`stowage venv <dir>`) reproduces
// it.
const defaultVenvManifestName = "stowage-env.json"

type venvManifestEntry struct {
	Name        string `json:"name"`
	Request     string `json:"request,omitempty"`
	SpecRequest string `json:"specRequest,omitempty"`
	Entrypoint  string `json:"entrypoint,omitempty"`
}

type venvManifest struct {
	Entries []venvManifestEntry `json:"entries"`
}

func loadVenvManifest(path string) (venvManifest, error) {
	var manifest venvManifest

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest, err
	}

	return manifest, nil
}

func saveVenvManifest(path string, manifest venvManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// upsert adds entry, replacing any existing entry with the same resolved
// command name.
func (m *venvManifest) upsert(entry venvManifestEntry) {
	for i, existing := range m.Entries {
		if existing.Name == entry.Name {
			m.Entries[i] = entry
			return
		}
	}
	m.Entries = append(m.Entries, entry)
}

func (m *venvManifest) removeByName(name string) {
	result := make([]venvManifestEntry, 0, len(m.Entries))
	for _, entry := range m.Entries {
		if entry.Name != name {
			result = append(result, entry)
		}
	}
	m.Entries = result
}

// updateVenvManifest loads the project's manifest (starting a fresh one if
// none exists yet), applies mutate, and saves it back.
func updateVenvManifest(mutate func(*venvManifest)) {
	manifest, _ := loadVenvManifest(defaultVenvManifestName)
	mutate(&manifest)
	if err := saveVenvManifest(defaultVenvManifestName, manifest); err != nil {
		fmt.Printf("WARNING: couldn't update %s: %s\n", defaultVenvManifestName, err)
	}
}
