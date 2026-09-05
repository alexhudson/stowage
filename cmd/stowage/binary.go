package main

import (
	"path/filepath"
)

// Binary is a type
type Binary struct {
	name   string
	spec   *Specification
	BinDir string // directory the wrapper is installed into
}

func (b *Binary) getWrapperPath() string {
	return filepath.Join(b.BinDir, b.name)
}
