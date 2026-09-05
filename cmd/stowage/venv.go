package main

import (
	"os"
	"path/filepath"
)

// venv represents a project-local stowage environment: a directory holding
// its own bin/ (wrapper binaries) and share/stowage/spec (specs for those
// wrappers), independent of any globally installed commands of the same
// name. It's activated by sourcing a generated script that prepends its
// bin/ to PATH, the same way `python -m venv` works.
type venv struct {
	Dir string // absolute path to the environment directory
}

func newVenv(dir string) (venv, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return venv{}, err
	}
	return venv{Dir: absDir}, nil
}

func (v *venv) binDir() string {
	return filepath.Join(v.Dir, "bin")
}

func (v *venv) storage() storage {
	return newStorage(
		filepath.Join(v.Dir, "share", "stowage", "spec"),
		filepath.Join(v.Dir, "share", "stowage", "repo"),
	)
}

func (v *venv) name() string {
	return filepath.Base(v.Dir)
}

func (v *venv) ensureDirs() error {
	if err := os.MkdirAll(v.binDir(), 0755); err != nil {
		return err
	}
	return os.MkdirAll(v.storage().specDir, 0755)
}

func (v *venv) install(spec *Specification) error {
	store := v.storage()
	if err := store.saveSpecification(spec); err != nil {
		return err
	}

	binary := Binary{name: spec.Name, spec: spec, BinDir: v.binDir()}
	return binary.install()
}

func (v *venv) remove(name string) error {
	binary := Binary{name: name, BinDir: v.binDir()}
	if err := binary.uninstall(); err != nil {
		return err
	}

	spec := Specification{Name: name}
	store := v.storage()
	return store.removeSpecification(&spec)
}

// activeVenv returns the currently activated environment, if any:
// install/uninstall/list target it instead of the global location whenever
// $STOWAGE_ENV is set (by a sourced activate script) - the same way `pip
// install` targets an activated virtualenv rather than system site-packages.
func activeVenv() (venv, bool) {
	dir := os.Getenv("STOWAGE_ENV")
	if dir == "" {
		return venv{}, false
	}

	env, err := newVenv(dir)
	if err != nil {
		return venv{}, false
	}

	return env, true
}
